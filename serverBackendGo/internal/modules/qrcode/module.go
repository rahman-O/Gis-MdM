package qrcode

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	qrhttp "github.com/gis-mdm/server-backend-go/internal/modules/qrcode/adapter/http"
	qrpostgres "github.com/gis-mdm/server-backend-go/internal/modules/qrcode/adapter/persistence/postgres"
	qrapp "github.com/gis-mdm/server-backend-go/internal/modules/qrcode/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "qrcode" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleQRCodeEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("qrcode module requires DATABASE_URL")
	}
	repo := qrpostgres.NewConfigRepository(deps.DB)
	svc := qrapp.NewService(repo, deps.Config.BaseURL, deps.Config.FilesDirectory, "", deps.Log)
	qrhttp.NewHandler(svc).Register(groups.Public.Group("/qr"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
