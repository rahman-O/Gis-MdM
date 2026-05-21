package application

import (
	"context"
	"errors"
	_ "image/gif"
	"image"
	"image/png"
	_ "image/jpeg"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	filesdomain "github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrIconDimensionInvalid = errors.New("error.icon.dimension.invalid")

type IconFileStore struct {
	FilesDir string
}

type UploadedFileInserter interface {
	InsertUploadedFile(ctx context.Context, f *filesdomain.UploadedFile) (*filesdomain.UploadedFile, error)
}

type CustomerFilesDir interface {
	CustomerFilesDir(ctx context.Context, customerID int) (string, error)
}

// UploadIconFile saves a square PNG icon (Java IconFileResource parity).
func UploadIconFile(
	ctx context.Context,
	p *platformauth.Principal,
	r io.Reader,
	store IconFileStore,
	customers CustomerFilesDir,
	files UploadedFileInserter,
) (*filesdomain.UploadedFile, error) {
	if p == nil {
		return nil, errors.New("error.permission.denied")
	}
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	if b.Dx() != b.Dy() {
		return nil, ErrIconDimensionInvalid
	}
	scaled := resizeImage(img, 144)
	filesDir, err := customers.CustomerFilesDir(ctx, p.CustomerID)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(store.FilesDir, filesDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	name := uuid.NewString() + ".png"
	full := filepath.Join(dir, name)
	out, err := os.Create(full)
	if err != nil {
		return nil, err
	}
	if err := png.Encode(out, scaled); err != nil {
		out.Close()
		return nil, err
	}
	out.Close()
	now := time.Now().UnixMilli()
	rec := &filesdomain.UploadedFile{
		CustomerID: p.CustomerID,
		FilePath:   name,
		UploadTime: now,
	}
	return files.InsertUploadedFile(ctx, rec)
}

func resizeImage(src image.Image, size int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	sb := src.Bounds()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx := sb.Min.X + (x*sb.Dx())/size
			sy := sb.Min.Y + (y*sb.Dy())/size
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}
