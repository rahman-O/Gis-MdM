package app

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	"github.com/gis-mdm/server-backend-go/internal/modules/applications"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth"
	platformjwt "github.com/gis-mdm/server-backend-go/internal/platform/jwt"
	"github.com/gis-mdm/server-backend-go/internal/modules/configfiles"
	"github.com/gis-mdm/server-backend-go/internal/modules/configurations"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices"
	"github.com/gis-mdm/server-backend-go/internal/modules/files"
	"github.com/gis-mdm/server-backend-go/internal/modules/groups"
	"github.com/gis-mdm/server-backend-go/internal/modules/hints"
	"github.com/gis-mdm/server-backend-go/internal/modules/icons"
	"github.com/gis-mdm/server-backend-go/internal/modules/notifications"
	"github.com/gis-mdm/server-backend-go/internal/modules/passwordreset"
	pluginaudit "github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit"
	plugindevicelog "github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog"
	plugindeviceinfo "github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo"
	pluginmessaging "github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging"
	pluginplatform "github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform"
	pluginpush "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push"
	"github.com/gis-mdm/server-backend-go/internal/modules/publicapi"
	"github.com/gis-mdm/server-backend-go/internal/modules/push"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode"
	"github.com/gis-mdm/server-backend-go/internal/modules/roles"
	"github.com/gis-mdm/server-backend-go/internal/modules/settings"
	"github.com/gis-mdm/server-backend-go/internal/modules/signup"
	"github.com/gis-mdm/server-backend-go/internal/modules/stats"
	"github.com/gis-mdm/server-backend-go/internal/modules/summary"
	"github.com/gis-mdm/server-backend-go/internal/modules/sync"
	"github.com/gis-mdm/server-backend-go/internal/modules/twofactor"
	"github.com/gis-mdm/server-backend-go/internal/modules/updates"
	"github.com/gis-mdm/server-backend-go/internal/modules/users"
)

func allModules(jwtProvider *platformjwt.Provider) []module.Module {
	return []module.Module{
		auth.NewModule(jwtProvider),
		signup.New(),
		passwordreset.New(),
		users.New(),
		twofactor.New(),
		roles.New(),
		customers.New(),
		settings.New(),
		hints.New(),
		summary.New(),
		devices.New(),
		groups.New(),
		applications.New(),
		configurations.New(),
		configfiles.New(),
		files.New(),
		icons.New(),
		publicapi.New(),
		stats.New(),
		sync.New(),
		push.New(),
		notifications.New(),
		updates.New(),
		qrcode.New(),
		pluginplatform.New(),
		pluginaudit.New(),
		pluginpush.New(),
		pluginmessaging.New(),
		plugindeviceinfo.New(),
		plugindevicelog.New(),
	}
}

func registerModules(groups module.RouteGroups, deps module.Dependencies, jwtProvider *platformjwt.Provider) error {
	for _, mod := range allModules(jwtProvider) {
		if err := mod.Register(groups, deps); err != nil {
			return fmt.Errorf("%s: %w", mod.Name(), err)
		}
	}
	return nil
}
