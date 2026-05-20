package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gis-mdm/server-backend-go/internal/config"
	"github.com/gis-mdm/server-backend-go/internal/module"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/database"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/middleware"
	"github.com/gis-mdm/server-backend-go/internal/platform/jwt"
	"github.com/gis-mdm/server-backend-go/internal/platform/logger"
)

// Run bootstraps config, database, HTTP server, and modules.
func Run() error {
	cfg := config.Load()
	log := logger.New()
	log.Info("starting server", slog.String("port", cfg.ServerPort))

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	if db != nil {
		defer db.Close()
		if err := database.Migrate(db); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	} else {
		log.Info("DATABASE_URL not set; running without persistence")
	}

	jwtProvider := jwt.NewProvider(jwt.Config{
		Secret:          cfg.JWTSecret,
		ValiditySeconds: cfg.JWTValiditySeconds,
		RememberSeconds: cfg.JWTValidityRememberSecs,
	})

	engine := httpx.NewEngine(cfg)
	middleware.SetupSessions(engine, cfg.SessionSecret)

	var lookup platformauth.UserLookup = noopLookup{}
	var enrich platformauth.PrincipalEnricher
	if db != nil {
		repo := postgres.NewUserRepository(db)
		lookup = repo
		enrich = repo
	}
	groups := httpx.BuildRouteGroups(engine, cfg, httpx.AuthWiring{
		JWT:    jwtProvider,
		Lookup: lookup,
		Enrich: enrich,
	})

	if cfg.SwaggerEnabled {
		registerSwagger(groups.Engine)
	}

	modDeps := module.Dependencies{Config: cfg, DB: db, Log: log}
	if err := registerModules(toModuleGroups(groups), modDeps, jwtProvider); err != nil {
		return fmt.Errorf("register modules: %w", err)
	}

	addr := ":" + cfg.ServerPort
	log.Info("listening", slog.String("addr", addr))
	return engine.Run(addr)
}

type noopLookup struct{}

func (noopLookup) LookupByLogin(context.Context, string) (*platformauth.Principal, error) {
	return nil, nil
}
