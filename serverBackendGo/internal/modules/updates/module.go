package updates

import (
	"context"
	"net/url"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/module"
	updhttp "github.com/gis-mdm/server-backend-go/internal/modules/updates/adapter/http"
	updapp "github.com/gis-mdm/server-backend-go/internal/modules/updates/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "updates" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleUpdatesEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	manifest := deps.Config.UpdateManifestURL
	if u, err := url.Parse(deps.Config.BaseURL); err == nil && u.Host != "" {
		manifest = updapp.SubstituteDomain(manifest, u.Host)
	}
	single := func(ctx context.Context) (bool, error) {
		if deps.DB == nil {
			return true, nil
		}
		var n int
		err := deps.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&n)
		return n <= 1, err
	}
	svc := updapp.NewService(strings.TrimSpace(manifest), single)
	updhttp.NewHandler(svc).Register(groups.Private.Group("/update"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
