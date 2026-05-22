package application

import (
	"context"
	"errors"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

// commitApplicationAPK moves a temp upload into tenant storage and sets public download URL(s).
func (s *Service) commitApplicationAPK(ctx context.Context, customerID int, app *domain.Application) error {
	if s == nil || s.store == nil || app == nil {
		return nil
	}
	f := apkUploadFields{
		filePath: app.FilePath, url: app.URL, urlArmeabi: app.URLArmeabi, urlArm64: app.URLArm64,
		arch: app.Arch, split: app.Split, version: app.Version,
	}
	if err := s.commitAPKFields(ctx, customerID, &f); err != nil {
		return err
	}
	app.FilePath, app.URL, app.URLArmeabi, app.URLArm64, app.Split = f.filePath, f.url, f.urlArmeabi, f.urlArm64, f.split
	return nil
}

func (s *Service) commitVersionAPK(ctx context.Context, customerID int, ver *domain.ApplicationVersion) error {
	if s == nil || s.store == nil || ver == nil {
		return nil
	}
	f := apkUploadFields{
		filePath: ver.FilePath, url: ver.URL, urlArmeabi: ver.URLArmeabi, urlArm64: ver.URLArm64,
		arch: ver.Arch, split: ver.Split, version: ver.Version,
	}
	if err := s.commitAPKFields(ctx, customerID, &f); err != nil {
		return err
	}
	ver.FilePath, ver.URL, ver.URLArmeabi, ver.URLArm64, ver.Split = f.filePath, f.url, f.urlArmeabi, f.urlArm64, f.split
	return nil
}

type apkUploadFields struct {
	filePath   *string
	url        *string
	urlArmeabi *string
	urlArm64   *string
	arch       *string
	split      *bool
	version    *string
}

func (s *Service) commitAPKFields(ctx context.Context, customerID int, f *apkUploadFields) error {
	if f == nil {
		return nil
	}
	fp := strings.TrimSpace(derefStr(f.filePath))
	if fp == "" {
		return nil
	}
	filesDir, err := s.repo.CustomerFilesDir(ctx, customerID)
	if err != nil {
		return err
	}
	rel := fp
	if storage.IsTempUploadPath(fp) {
		rel, err = s.moveTempAPK(filesDir, fp)
		if err != nil {
			return err
		}
		f.filePath = &rel
	}
	publicURL := storage.BuildPublicURL(s.baseURL, filesDir, rel)
	if strings.TrimSpace(publicURL) == "" {
		return errors.New("BASE_URL is not configured; cannot publish APK download URL")
	}
	arch := strings.TrimSpace(derefStr(f.arch))
	split := f.split != nil && *f.split
	if strings.TrimSpace(derefStr(f.url)) == "" &&
		strings.TrimSpace(derefStr(f.urlArmeabi)) == "" &&
		strings.TrimSpace(derefStr(f.urlArm64)) == "" {
		switch {
		case split && arch == "armeabi":
			f.urlArmeabi = &publicURL
		case split && arch == "arm64":
			f.urlArm64 = &publicURL
		default:
			f.url = &publicURL
			if f.split != nil {
				b := false
				f.split = &b
			}
		}
	}
	return nil
}

func (s *Service) moveTempAPK(filesDir, tmpPath string) (string, error) {
	rel, err := s.store.MoveToCustomer(filesDir, "", tmpPath, "")
	if errors.Is(err, storage.ErrExists) {
		if name, nameErr := storage.NameFromTmpPath(tmpPath); nameErr == nil {
			_ = s.store.DeleteRelative(filesDir, name)
		}
		rel, err = s.store.MoveToCustomer(filesDir, "", tmpPath, "")
	}
	return rel, err
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
