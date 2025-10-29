package app

import (
	"context"
	"dumper/internal/app/automation"
	"dumper/internal/app/manual"
	_ "dumper/internal/command/database/mariadb"
	_ "dumper/internal/command/database/mongodb"
	_ "dumper/internal/command/database/mssql"
	_ "dumper/internal/command/database/mysql"
	_ "dumper/internal/command/database/postgres"
	_ "dumper/internal/command/database/redis"
	_ "dumper/internal/command/database/sqlite"
	"dumper/internal/domain/app"
	cfg "dumper/internal/domain/config"
	"dumper/pkg/logging"
	"fmt"
	"strings"
)

type App struct {
	ctx context.Context
	cfg *cfg.Config
	env *app.Env
}

func NewApp(ctx context.Context, cfg *cfg.Config, env *app.Env) *App {
	return &App{
		ctx: ctx,
		cfg: cfg,
		env: env,
	}
}

func (a *App) MustRun() error {
	if err := a.Run(); err != nil {
		logging.L(a.ctx).Error("App failed to run")
		return fmt.Errorf("%v", err)
	}
	return nil
}

func (a *App) Run() error {
	if a.env.All == false && a.env.DbName != "" {
		logging.L(a.ctx).Info("Running the app with the parameters specified (db list)")
		automationDumpApp := automation.NewApp(a.ctx, a.cfg, a.env)
		return automationDumpApp.Run()
	}

	if a.env.All == true && a.env.DbName == "" {
		logging.L(a.ctx).Info("Running the app with the parameters specified (db all)")
		var keys []string
		for key := range a.cfg.Databases {
			keys = append(keys, key)
		}
		a.env.DbName = strings.Join(keys, ",")

		automationDumpApp := automation.NewApp(a.ctx, a.cfg, a.env)
		return automationDumpApp.Run()
	}

	logging.L(a.ctx).Info("Running the app in manual mode with db selection")
	manualDumpApp := manual.NewApp(a.ctx, a.cfg, a.env)
	return manualDumpApp.Run()
}
