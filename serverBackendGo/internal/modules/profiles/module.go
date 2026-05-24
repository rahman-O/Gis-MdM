package profiles

import (
	"context"
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
	profilehttp "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/http"
	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	profileapp "github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	pushapp "github.com/gis-mdm/server-backend-go/internal/platform/push/application"
)

// Module registers profile routes (017-device-control-plane US3).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "profiles" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleProfilesEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("profiles module requires DATABASE_URL")
	}
	repo := profilepostgres.NewProfileRepository(deps.DB)
	draft := profileapp.NewDraftService(repo)
	compiler := profileapp.NewArtifactCompiler(deps.DB, deps.Config.BaseURL, deps.Config.FilesDirectory)
	var rollout port.RolloutStore
	if deps.Config.ModuleProfileRolloutEnabled {
		rollout = profilepostgres.NewAssignmentRepository(deps.DB)
	}
	publish := profileapp.NewPublishService(repo, rollout, deps.DB, draft, compiler)
	versionDelete := profileapp.NewVersionDeleteService(repo, deps.DB)
	hub := profileapp.NewHubService(deps.DB, repo, deps.Config.ProfileStalePublishDays)
	h := profilehttp.NewHandler(draft, publish, hub, versionDelete)
	profGroup := groups.Private.Group("/profiles")
	h.Register(profGroup)
	if deps.Config.ModuleProfileRolloutEnabled {
		assign := profileapp.NewAssignmentService(deps.DB)
		rollout := profileapp.NewRolloutStatusService(deps.DB)
		enable := profileapp.NewEnableService(deps.DB)
		profilehttp.NewRolloutHandlers(assign, rollout, enable, draft).RegisterOnProfile(profGroup)
	}
	onboarding := profileapp.NewOnboardingService(deps.DB)
	profilehttp.NewOnboardingHandler(onboarding).Register(groups.Private.Group("/onboarding"))
	if deps.Config.PushNotifierEnabled {
		queue := notifpostgres.NewQueueRepository(deps.DB)
		pushapp.NewDomainEventsWorker(deps.DB, queue, deps.Log).Start(context.Background())
	}
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
