package application

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/publicapi/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/publicapi/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
	sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

var (
	ErrInvalidHash      = errors.New("Invalid hash")
	ErrDeviceNotFound   = errors.New("error.notfound.device")
	ErrDuplicateApp     = errors.New("error.duplicate.application")
	ErrMissingParams    = errors.New("error.params.missing")
	ErrPermissionDenied = errors.New("error.permission.denied")
)

type RebrandingConfig struct {
	AppName     string
	VendorName  string
	VendorLink  string
	SignupLink  string
	TermsLink   string
	LogoPath    string
	HashSecret  string
}

type Service struct {
	repo   port.DeviceRepository
	store  *storage.LocalStore
	base   string
	rebrand RebrandingConfig
}

func NewService(repo port.DeviceRepository, store *storage.LocalStore, baseURL string, cfg RebrandingConfig) *Service {
	return &Service{repo: repo, store: store, base: strings.TrimRight(baseURL, "/"), rebrand: cfg}
}

func (s *Service) GetName() domain.NameResponse {
	return domain.NameResponse{
		AppName:    s.rebrand.AppName,
		VendorName: s.rebrand.VendorName,
		VendorLink: s.rebrand.VendorLink,
		SignupLink: s.rebrand.SignupLink,
		TermsLink:  s.rebrand.TermsLink,
	}
}

func (s *Service) LogoPath() string { return s.rebrand.LogoPath }

func (s *Service) UploadApplication(ctx context.Context, appJSON string, fileName string, r io.Reader) error {
	var req domain.UploadAppRequest
	if err := json.Unmarshal([]byte(appJSON), &req); err != nil {
		return ErrMissingParams
	}
	req.DeviceID = strings.Trim(req.DeviceID, "\"")
	req.Hash = strings.Trim(req.Hash, "\"")
	if req.Name == "" || req.Pkg == "" || req.Version == "" || req.DeviceID == "" || req.Hash == "" {
		return ErrMissingParams
	}
	if fileName != "" && r != nil {
		if req.LocalPath == "" || req.FileName == "" {
			return ErrMissingParams
		}
	}
	if !storage.IsSafePath(req.LocalPath) || !storage.IsSafePath(req.FileName) {
		return ErrPermissionDenied
	}
	want := sharedcrypto.DeviceUploadHash(req.DeviceID, s.rebrand.HashSecret)
	if !strings.EqualFold(want, req.Hash) {
		return ErrInvalidHash
	}
	dev, err := s.repo.FindDeviceByNumber(ctx, req.DeviceID)
	if err != nil {
		return err
	}
	if dev == nil {
		return ErrDeviceNotFound
	}
	dup, err := s.repo.HasDuplicateApp(ctx, dev.CustomerID, req.Pkg, req.Version)
	if err != nil {
		return err
	}
	if dup {
		return ErrDuplicateApp
	}
	filesDir, err := s.repo.CustomerFilesDir(ctx, dev.CustomerID)
	if err != nil {
		return err
	}
	var url string
	if fileName != "" && r != nil {
		destDir := filepath.Join(s.store.CustomerRoot(filesDir), filepath.FromSlash(req.LocalPath))
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return err
		}
		dest := filepath.Join(destDir, req.FileName)
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, r); err != nil {
			out.Close()
			return err
		}
		out.Close()
		url = storage.BuildPublicURL(s.base, filesDir, filepath.ToSlash(filepath.Join(req.LocalPath, req.FileName)))
	}
	return s.repo.InsertApplication(ctx, dev.CustomerID, req.Name, req.Pkg, req.Version, url, req)
}
