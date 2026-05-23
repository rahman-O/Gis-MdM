package auth

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	authhttp "github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/http"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	platformcrypto "github.com/gis-mdm/server-backend-go/internal/platform/crypto"
	"github.com/gis-mdm/server-backend-go/internal/platform/email"
	platformjwt "github.com/gis-mdm/server-backend-go/internal/platform/jwt"
)

// Module registers authentication HTTP routes.
type Module struct {
	jwt *platformjwt.Provider
}

// NewModule creates an auth module with JWT provider.
func NewModule(jwt *platformjwt.Provider) *Module {
	return &Module{jwt: jwt}
}

func (m *Module) Name() string { return "auth" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleAuthEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("auth module requires DATABASE_URL")
	}
	repo := postgres.NewUserRepository(deps.DB)
	emailSvc := email.NewService(deps.Config.EmailConfigured, deps.Log)
	var rsaKeys *platformcrypto.RSAKeys
	if deps.Config.TransmitPassword {
		rsaKeys = platformcrypto.NewRSAKeys(deps.Config.FilesDirectory)
		_ = rsaKeys.EnsureKeys()
	}
	svc := application.NewService(
		repo,
		m.jwt,
		emailSvc,
		rsaKeys,
		deps.Config.TransmitPassword,
		deps.Config.CustomerSignup,
	)
	authhttp.Register(groups.Public, authhttp.NewHandler(svc))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
