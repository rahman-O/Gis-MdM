package application

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/files/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrFileUsed         = errors.New("error.used.file")
	ErrFileExists       = errors.New("error.duplicate.file")
	ErrUnsafePath       = errors.New("unsafe path")
	ErrSizeLimit        = errors.New("error.size.limit.exceeded")
	ErrInternal         = errors.New("error.internal.server")
	ErrSaveFile         = errors.New("error.file.save")
)

// Service implements file library use cases.
type Service struct {
	files    port.FileRepository
	customer port.CustomerRepository
	apps     port.ApplicationLookup
	store    *storage.LocalStore
	baseURL  string
	push     port.PushNotifier
}

func NewService(
	files port.FileRepository,
	customer port.CustomerRepository,
	apps port.ApplicationLookup,
	store *storage.LocalStore,
	baseURL string,
	push port.PushNotifier,
) *Service {
	if push == nil {
		push = port.NoopPush()
	}
	return &Service{files: files, customer: customer, apps: apps, store: store, baseURL: baseURL, push: push}
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal, filter string) ([]domain.FileView, error) {
	if p == nil || !p.CanBrowseFiles() {
		return nil, ErrPermissionDenied
	}
	rows, err := s.files.List(ctx, p.CustomerID, filter)
	if err != nil {
		return nil, err
	}
	return s.toViews(ctx, p.CustomerID, rows)
}

func (s *Service) Remove(ctx context.Context, p *platformauth.Principal, id int, filePath string, external bool) error {
	if p == nil || !p.CanEditFiles() {
		return ErrPermissionDenied
	}
	if !storage.IsSafePath(filePath) {
		return ErrPermissionDenied
	}
	usedCfg, _ := s.files.IsUsedByConfiguration(ctx, id)
	usedIcon, _ := s.files.IsUsedByIcon(ctx, id)
	if usedCfg || usedIcon {
		return ErrFileUsed
	}
	meta, err := s.customer.GetMeta(ctx, p.CustomerID)
	if err != nil {
		return err
	}
	if err := s.files.Delete(ctx, id); err != nil {
		return err
	}
	if !external && filePath != "" {
		_ = s.store.DeleteRelative(meta.FilesDir, filePath)
	}
	return nil
}

func (s *Service) GetLimit(ctx context.Context, p *platformauth.Principal) (domain.LimitResponse, error) {
	out := domain.LimitResponse{}
	if p == nil {
		return out, nil
	}
	count, _ := s.customer.CountCustomers(ctx)
	if count <= 1 {
		return out, nil
	}
	meta, err := s.customer.GetMeta(ctx, p.CustomerID)
	if err != nil || meta.Master || meta.SizeLimit <= 0 {
		return out, err
	}
	bytes, err := s.store.DirSizeBytes(meta.FilesDir)
	if err != nil {
		return out, nil
	}
	out.SizeUsed = int(bytes / 1048576)
	out.SizeLimit = meta.SizeLimit
	return out, nil
}

func (s *Service) Upload(ctx context.Context, p *platformauth.Principal, name string, r io.Reader, parseAPK bool) (*domain.FileUploadResult, error) {
	if p == nil || !p.CanEditFiles() {
		return nil, ErrPermissionDenied
	}
	if err := s.checkQuota(ctx, p.CustomerID, 0); err != nil {
		return nil, err
	}
	path, err := s.store.CreateTemp(name, r)
	if err != nil {
		return nil, ErrInternal
	}
	result := &domain.FileUploadResult{Name: name, ServerPath: path}
	if parseAPK {
		result.FileDetails = ParseAPK(path)
		s.enrichUploadHints(ctx, p.CustomerID, result)
	}
	return result, nil
}

func (s *Service) Commit(ctx context.Context, p *platformauth.Principal, in domain.UploadedFile) (*domain.UploadedFile, error) {
	if p == nil || !p.CanEditFiles() {
		return nil, ErrPermissionDenied
	}
	meta, err := s.customer.GetMeta(ctx, p.CustomerID)
	if err != nil {
		return nil, err
	}
	if in.ID == nil || *in.ID == 0 {
		if in.External {
			return s.createExternal(ctx, p, in)
		}
		return s.createInternal(ctx, p, meta, in)
	}
	if err := s.updateExisting(ctx, p, meta, in); err != nil {
		return nil, err
	}
	return &in, nil
}

func (s *Service) createExternal(ctx context.Context, p *platformauth.Principal, in domain.UploadedFile) (*domain.UploadedFile, error) {
	if strings.TrimSpace(in.ExternalURL) == "" {
		return nil, ErrInternal
	}
	in.CustomerID = p.CustomerID
	in.FilePath = ""
	in.UploadTime = time.Now().UnixMilli()
	if err := s.files.Insert(ctx, &in); err != nil {
		return nil, err
	}
	in.URL = in.ExternalURL
	return &in, nil
}

func (s *Service) createInternal(ctx context.Context, p *platformauth.Principal, meta *domain.CustomerMeta, in domain.UploadedFile) (*domain.UploadedFile, error) {
	if !storage.IsSafePath(in.FilePath) || !isUnderTemp(in.TmpPath) {
		return nil, ErrPermissionDenied
	}
	name := in.FileName
	if name == "" {
		var err error
		name, err = storage.NameFromTmpPath(in.TmpPath)
		if err != nil {
			return nil, ErrPermissionDenied
		}
		in.FilePath = name
	}
	if dup, _ := s.files.CountByPath(ctx, p.CustomerID, nil, in.FilePath); dup > 0 {
		return nil, ErrFileExists
	}
	rel, err := s.store.MoveToCustomer(meta.FilesDir, in.Subdir, in.TmpPath, name)
	if err == storage.ErrExists {
		return nil, ErrFileExists
	}
	if err != nil {
		return nil, ErrSaveFile
	}
	in.FilePath = rel
	in.CustomerID = p.CustomerID
	in.UploadTime = time.Now().UnixMilli()
	if err := s.files.Insert(ctx, &in); err != nil {
		return nil, err
	}
	in.URL = storage.BuildPublicURL(s.baseURL, meta.FilesDir, rel)
	return &in, nil
}

func (s *Service) updateExisting(ctx context.Context, p *platformauth.Principal, meta *domain.CustomerMeta, in domain.UploadedFile) error {
	db, err := s.files.GetByID(ctx, p.CustomerID, *in.ID)
	if err != nil || db == nil {
		return ErrInternal
	}
	if !in.External && in.TmpPath != "" {
		if !isUnderTemp(in.TmpPath) {
			return ErrPermissionDenied
		}
		_ = s.store.DeleteRelative(meta.FilesDir, db.FilePath)
		rel, err := s.store.MoveToCustomer(meta.FilesDir, "", in.TmpPath, filepath.Base(db.FilePath))
		if err != nil {
			return ErrSaveFile
		}
		in.FilePath = rel
		in.UploadTime = time.Now().UnixMilli()
	}
	if !storage.IsSafePath(in.FilePath) {
		return ErrPermissionDenied
	}
	in.CustomerID = p.CustomerID
	return s.files.Update(ctx, &in)
}

func (s *Service) GetApplicationsByURL(ctx context.Context, p *platformauth.Principal, url string) (any, error) {
	if p == nil || !p.CanBrowseFiles() {
		return nil, ErrPermissionDenied
	}
	return s.apps.SearchByURL(ctx, p.CustomerID, url)
}

func (s *Service) GetFileConfigurations(ctx context.Context, p *platformauth.Principal, fileID int) ([]domain.FileConfigurationLink, error) {
	if p == nil || !p.CanBrowseFiles() {
		return nil, ErrPermissionDenied
	}
	return s.files.GetFileConfigurations(ctx, p.CustomerID, int(p.ID), fileID)
}

func (s *Service) UpdateFileConfigurations(ctx context.Context, p *platformauth.Principal, req domain.LinkConfigurationsToFileRequest) error {
	if p == nil || !p.CanEditFiles() {
		return ErrPermissionDenied
	}
	for _, link := range req.Configurations {
		if link.ID != nil && !link.Upload {
			if err := s.files.DeleteConfigurationFile(ctx, *link.ID); err != nil {
				return err
			}
		}
		if link.ID == nil && link.Upload {
			if err := s.files.InsertConfigurationFile(ctx, link.ConfigurationID, link.FileID, link.FileName); err != nil {
				return err
			}
		}
		if link.Notify {
			s.push.NotifyConfigurationUpdate(link.ConfigurationID)
		}
	}
	return nil
}

func (s *Service) toViews(ctx context.Context, customerID int, files []domain.UploadedFile) ([]domain.FileView, error) {
	meta, _ := s.customer.GetMeta(ctx, customerID)
	views := make([]domain.FileView, 0, len(files))
	for _, f := range files {
		v := domain.FileView{
			ID:               f.ID,
			FilePath:         f.FilePath,
			Description:      f.Description,
			UploadTime:       f.UploadTime,
			DevicePath:       f.DevicePath,
			External:         f.External,
			ReplaceVariables: f.ReplaceVariables,
		}
		if f.External {
			v.URL = f.ExternalURL
		} else if meta != nil && s.store != nil {
			v.URL = storage.BuildPublicURL(s.baseURL, meta.FilesDir, f.FilePath)
			full := filepath.Join(s.store.CustomerRoot(meta.FilesDir), filepath.FromSlash(f.FilePath))
			if info, err := os.Stat(full); err == nil {
				v.Size = info.Size()
			} else {
				v.Size = -1
			}
		}
		if f.ID != nil {
			v.UsedByConfigurations, _ = s.files.UsingConfigurationNames(ctx, customerID, *f.ID)
			v.UsedByIcons, _ = s.files.UsingIconNames(ctx, customerID, *f.ID)
		}
		views = append(views, v)
	}
	return views, nil
}

func (s *Service) checkQuota(ctx context.Context, customerID int, extraBytes int64) error {
	count, _ := s.customer.CountCustomers(ctx)
	if count <= 1 {
		return nil
	}
	meta, err := s.customer.GetMeta(ctx, customerID)
	if err != nil || meta.Master || meta.SizeLimit <= 0 {
		return err
	}
	used, err := s.store.DirSizeBytes(meta.FilesDir)
	if err != nil {
		return nil
	}
	totalMB := (used + extraBytes) / 1048576
	if int(totalMB) > meta.SizeLimit {
		return ErrSizeLimit
	}
	return nil
}

func (s *Service) enrichUploadHints(ctx context.Context, customerID int, result *domain.FileUploadResult) {
	if result.FileDetails == nil || result.FileDetails.Pkg == "" {
		return
	}
	pkg := result.FileDetails.Pkg
	if result.FileDetails.VersionCode > 0 {
		if v, _ := s.apps.FindVersionByPkgCode(ctx, customerID, pkg, result.FileDetails.VersionCode); v != nil {
			if result.FileDetails.Version != "" && v.Version != nil && *v.Version != result.FileDetails.Version {
				// version code conflict — caller may map to duplicate message
			}
		}
	}
	if result.FileDetails.Version != "" {
		if v, _ := s.apps.FindVersionByPkgVersion(ctx, customerID, pkg, result.FileDetails.Version); v != nil {
			ex := true
			result.Exists = &ex
		}
	}
	if apps, _ := s.apps.FindAppsByPkg(ctx, customerID, pkg); len(apps) > 0 {
		result.Application = apps[0]
	}
}

func isUnderTemp(path string) bool {
	if path == "" {
		return false
	}
	tmp := os.TempDir()
	abs, err := filepath.Abs(path)
	if err != nil {
		return strings.HasPrefix(path, tmp)
	}
	return strings.HasPrefix(abs, tmp)
}
