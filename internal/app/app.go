package app

import (
	"context"
	"dumper/internal/backup"
	"dumper/internal/command"
	_ "dumper/internal/command/mysql"
	_ "dumper/internal/command/postgres"
	"dumper/internal/config"
	"dumper/internal/connect"
	cmdCfg "dumper/internal/domain/command-config"
	_select "dumper/internal/select"
	t "dumper/internal/temr"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"fmt"
)

type Env struct {
	ConfigFile string
	DbName     string
	All        bool
	FileLog    string
}

type App struct {
	ctx context.Context
	cfg *config.Config
	env *Env
}

func NewApp(ctx context.Context, cfg *config.Config, env *Env) *App {
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
		return a.RunDumpDB()
	}

	if a.env.All == true && a.env.DbName == "" {
		logging.L(a.ctx).Info("Running the app with the parameters specified (db all)")
		return a.RunDumpAll()
	}

	logging.L(a.ctx).Info("Running the app in manual mode with db selection")
	return a.RunDumpManual()
}

func (a *App) RunDumpManual() error {
	logging.L(a.ctx).Info("Prepare server list")

	m := t.New()
	serverList, serverKeys := _select.SelectOptionList(a.cfg.Servers, "")
	m.SetList(serverKeys)
	m.SetTitle("Select server ")
	m.Run()

	serverName := m.GetSelect()
	serverKey := serverList[serverName]
	server := a.cfg.Servers[serverKey]

	logging.L(a.ctx).Info("Selected server", logging.StringAttr("server", serverKey))

	m.ClearList()

	logging.L(a.ctx).Info("Prepare database list")

	dbList, dbKeys := _select.SelectOptionList(a.cfg.Databases, serverKey)
	m.SetList(dbKeys)
	m.SetTitle("Select database ")
	m.Run()

	dbName := m.GetSelect()
	dbKey := dbList[dbName]
	db := a.cfg.Databases[dbKey]

	logging.L(a.ctx).Info("Selected database", logging.StringAttr("database", dbKey))

	dataFormat := utils.TemplateData{
		Server:   serverName,
		Database: dbName,
		Template: a.cfg.Settings.Template,
	}

	nameFile := utils.GetTemplateFileName(dataFormat)

	logging.L(a.ctx).Info("Generate template", logging.StringAttr("name", nameFile))

	cmdData := &cmdCfg.ConfigData{
		User:       db.User,
		Password:   db.Password,
		Name:       db.GetDisplayName(),
		Port:       db.GetPort(a.cfg.Settings.DBPort),
		Key:        server.SSHKey,
		Host:       server.Host,
		DumpName:   nameFile,
		DumpFormat: a.cfg.Settings.DumpFormat,
	}

	logging.L(a.ctx).Info("Prepare command for dump")

	cmdApp := command.NewApp(&a.cfg.Settings, cmdData)
	cmdStr, remotePath, err := cmdApp.GetCommand()

	if err != nil {
		logging.L(a.ctx).Error("error generate command")
		return fmt.Errorf("error generate command: %w", err)
	}

	logging.L(a.ctx).Info("Prepare connection")
	conn := connect.New(
		server.Host,
		server.User,
		server.GetPort(a.cfg.Settings.SrvPost),
		a.cfg.Settings.SSH.PrivateKey,
		server.SSHKey,
		a.cfg.Settings.SSH.Passphrase,
		server.Password,
		*a.cfg.Settings.SSH.IsPassphrase,
	)

	if err := conn.Connect(); err != nil {
		logging.L(a.ctx).Error("Error connecting to server")
		return err
	}

	defer func(conn *connect.Connect) {
		_ = conn.Close()
	}(conn)

	logging.L(a.ctx).Info("Run test connection to server")

	if err := conn.TestConnection(); err != nil {
		logging.L(a.ctx).Error("Error connection to server")
		return err
	}

	logging.L(a.ctx).Info("The connection is established")

	logging.L(a.ctx).Info("Preparing for backup creation")

	backupApp := backup.NewApp(a.ctx, conn, cmdStr, remotePath, a.cfg.Settings.DirDump, a.cfg.Settings.DumpLocation)
	if err := backupApp.Backup(); err != nil {
		logging.L(a.ctx).Error("Error create backup database")
		return err
	}

	logging.L(a.ctx).Info("Backup was successfully created and downloaded")

	if a.cfg.Settings.DirArchived != "" {
		logging.L(a.ctx).Info("Search for old backups")

		dbNamePrefix := fmt.Sprintf("%s_%s", serverName, dbName)
		err = utils.ArchivedLocalFile(dbNamePrefix, remotePath, a.cfg.Settings.DirDump, a.cfg.Settings.DirArchived)
		if err != nil {
			logging.L(a.ctx).Error("Error archiving old backups")
			return err
		}

		logging.L(a.ctx).Info("Archived old backups", logging.StringAttr("path", a.cfg.Settings.DirArchived))
	}

	return nil
}

func (a *App) RunDumpAll() error {
	panic("implement me")
}

func (a *App) RunDumpDB() error {
	panic("implement me")
}
