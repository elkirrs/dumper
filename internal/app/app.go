package app

import (
	"context"
	"dumper/internal/backup"
	"dumper/internal/command"
	_ "dumper/internal/command/mariadb"
	_ "dumper/internal/command/mongodb"
	_ "dumper/internal/command/mssql"
	_ "dumper/internal/command/mysql"
	_ "dumper/internal/command/postgres"
	"dumper/internal/config"
	"dumper/internal/config/remote"
	"dumper/internal/connect"
	cmdCfg "dumper/internal/domain/command-config"
	_select "dumper/internal/select"
	t "dumper/internal/temr"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
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

type Remote interface {
	Load() error
	Config() map[string]config.DBConnect
}

func NewApp(ctx context.Context, cfg *config.Config, env *Env) *App {
	return &App{
		ctx: ctx,
		cfg: cfg,
		env: env,
	}
}

type ConnectionError struct {
	Addr string
	Err  error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection error to %s: %v", e.Addr, e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
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
		var keys []string
		for key := range a.cfg.Databases {
			keys = append(keys, key)
		}
		a.env.DbName = strings.Join(keys, ",")
		return a.RunDumpDB()
	}

	logging.L(a.ctx).Info("Running the app in manual mode with db selection")
	return a.RunDumpManual()
}

func (a *App) RunDumpManual() error {
	logging.L(a.ctx).Info("Prepare server list")

	var serverName string
	var err error

	term := t.New()
	serverList, serverKeys := _select.SelectOptionList(a.cfg.Servers, "")

	if len(a.cfg.Servers) > 1 {
		term.SetList(serverKeys)
		term.SetTitle("Select server ")

		if err = runWithCtx(a.ctx, func() error { term.Run(); return nil }); err != nil {
			return err
		}

		serverName = term.GetSelect()
	} else {
		serverName = serverKeys[0]
		fmt.Println("\033[32m" + "\U00002714 " + serverName + "\033[0m")
	}

	serverKey := serverList[serverName]
	server := a.cfg.Servers[serverKey]

	logging.L(a.ctx).Info("Selected server", logging.StringAttr("server", serverKey))

	term.ClearList()
	logging.L(a.ctx).Info("Prepare database list")

	var dataDBConnect map[string]config.DBConnect

	if server.ConfigPath != "" {
		dataDBConnect, err = a.prepareRemoteDatabaseList(server)
		serverKey = ""
		if err != nil {
			return err
		}
	} else {
		dataDBConnect = a.prepareDBConnect()
	}

	dbList, dbKeys := _select.OptionDataBaseList(dataDBConnect, serverKey)

	term.SetList(dbKeys)
	term.SetTitle("Select database ")

	if err = runWithCtx(a.ctx, func() error { term.Run(); return nil }); err != nil {
		return err
	}

	dbName := term.GetSelect()
	dbKey := dbList[dbName]
	dbConnect := dataDBConnect[dbKey]

	logging.L(a.ctx).Info("Selected database", logging.StringAttr("database", dbKey))

	err = withRetry(
		a.ctx, a.cfg.Settings.RetryConnect,
		func() error {
			return a.runBackup(dbConnect)
		},
		func(err error) bool {
			var connErr *ConnectionError
			return errors.As(err, &connErr)
		},
		func(attempt int, err error) {
			var connErr *ConnectionError
			_ = errors.As(err, &connErr)
			logging.L(a.ctx).Warn(
				"Connection error, retrying",
				logging.StringAttr("db", dbConnect.Database.Name),
				logging.StringAttr("addr", connErr.Addr),
				logging.IntAttr("attempt", attempt),
				logging.ErrAttr(err),
			)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) RunDumpDB() error {
	logging.L(a.ctx).Info("Prepare data for create dumps")

	dbList := strings.Split(a.env.DbName, ",")
	countDBs := len(dbList)

	serversDatabases := make(map[string][]config.DBConnect)

	for _, dbName := range dbList {
		database, ok := a.cfg.Databases[dbName]
		if !ok {
			fmt.Printf("Database %s not found\n", dbName)
			logging.L(a.ctx).Warn("Database not found", logging.StringAttr("name", dbName))
			countDBs--
			continue
		}

		server, ok := a.cfg.Servers[database.Server]
		if !ok {
			fmt.Printf("Server %s not found\n", database.Server)
			logging.L(a.ctx).Warn("Server not found", logging.StringAttr("name", database.Server))
			countDBs--
			continue
		}

		serversDatabases[database.Server] = append(serversDatabases[database.Server], config.DBConnect{
			Server:   server,
			Database: database,
		})

		if countDBs == 0 {
			logging.L(a.ctx).Error("Database and server no key matches check the configuration file")
			return errors.New("database and server no key matches check the configuration file")
		}
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error, len(dbList))

	for _, dbInfoList := range serversDatabases {
		dbListCopy := dbInfoList
		wg.Add(1)
		go func(connectDBs []config.DBConnect) {
			defer wg.Done()
			for _, dbConnect := range connectDBs {
				select {
				case <-a.ctx.Done():
					logging.L(a.ctx).Info("Backup cancelled by context")
					errCh <- fmt.Errorf("backup cancelled for database %s", dbConnect.Database.Name)
					return
				default:
					err := withRetry(
						a.ctx, a.cfg.Settings.RetryConnect,
						func() error {
							return a.runBackup(dbConnect)
						},
						func(err error) bool {
							var connErr *ConnectionError
							return errors.As(err, &connErr)
						},
						func(attempt int, err error) {
							var connErr *ConnectionError
							_ = errors.As(err, &connErr)
							logging.L(a.ctx).Warn(
								"Connection error, retrying",
								logging.StringAttr("db", dbConnect.Database.Name),
								logging.StringAttr("addr", connErr.Addr),
								logging.IntAttr("attempt", attempt),
								logging.ErrAttr(err),
							)
						},
					)

					if err != nil {
						errCh <- fmt.Errorf("backup failed for %s: %w", dbConnect.Database.Name, err)
						return
					}
				}
			}
		}(dbListCopy)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	logging.L(a.ctx).Info("All requested database backups are done")

	return nil
}

func (a *App) runBackup(dbConnect config.DBConnect) error {
	server := dbConnect.Server
	db := dbConnect.Database

	dataFormat := utils.TemplateData{
		Server:   server.GetName(),
		Database: db.GetName(),
		Template: a.cfg.Settings.Template,
	}
	nameFile := utils.GetTemplateFileName(dataFormat)
	logging.L(a.ctx).Info("Generated template", logging.StringAttr("name", nameFile))

	cmdData := &cmdCfg.ConfigData{
		User:       db.User,
		Password:   db.Password,
		Name:       db.GetName(),
		Port:       db.GetPort(a.cfg.Settings.DBPort),
		Key:        server.SSHKey,
		Host:       server.Host,
		DumpName:   nameFile,
		DumpFormat: a.cfg.Settings.DumpFormat,
		Driver:     db.GetDriver(a.cfg.Settings.Driver),
		Options:    db.Options,
	}

	logging.L(a.ctx).Info("Prepare command for dump")

	cmdApp := command.NewApp(&a.cfg.Settings, cmdData)
	cmdStr, remotePath, err := cmdApp.GetCommand()

	if err != nil {
		logging.L(a.ctx).Error("failed generate command")
		return fmt.Errorf("failed generate command: %w", err)
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

	fmt.Printf("Trying to connect to server %s...\n", server.Host)
	if err := runWithCtx(a.ctx, conn.Connect); err != nil {
		logging.L(a.ctx).Error(
			"Error connecting to server",
			logging.StringAttr("server", server.Host),
			logging.ErrAttr(err),
		)
		return &ConnectionError{Addr: server.Host, Err: err}
	}

	defer func(conn *connect.Connect) {
		_ = conn.Close()
	}(conn)

	logging.L(a.ctx).Info("Trying to test connection to server")
	if err := runWithCtx(a.ctx, conn.TestConnection); err != nil {
		logging.L(a.ctx).Error("Error testing connection to server")
		return &ConnectionError{Addr: server.Host, Err: err}
	}
	logging.L(a.ctx).Info("The connection is established")

	logging.L(a.ctx).Info("Preparing for backup creation")

	backupApp := backup.NewApp(
		a.ctx,
		conn,
		cmdStr,
		remotePath,
		a.cfg.Settings.DirDump,
		a.cfg.Settings.DumpLocation,
		db.GetRemoveDump(*a.cfg.Settings.RemoveDump),
	)

	if err := runWithCtx(a.ctx, backupApp.Backup); err != nil {
		logging.L(a.ctx).Error("Error creating backup database")
		return err
	}
	logging.L(a.ctx).Info("Backup was successfully created and downloaded")

	if a.cfg.Settings.DirArchived != "" {
		logging.L(a.ctx).Info("Search for old backups")
		dbNamePrefix := fmt.Sprintf("%s_%s", server.GetName(), db.GetName())

		if err := runWithCtx(a.ctx, func() error {
			return utils.ArchivedLocalFile(dbNamePrefix, remotePath, a.cfg.Settings.DirDump, a.cfg.Settings.DirArchived)
		}); err != nil {
			logging.L(a.ctx).Error("Error archiving old backups")
			return err
		}

		logging.L(a.ctx).Info("Archived old backups", logging.StringAttr("path", a.cfg.Settings.DirArchived))
	}

	return nil
}

func runWithCtx(ctx context.Context, f func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- f()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled")
	case err := <-done:
		return err
	}
}

func withRetry(
	ctx context.Context,
	maxRetries int,
	fn func() error,
	shouldRetry func(err error) bool,
	onRetry func(attempt int, err error),
) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if shouldRetry != nil && shouldRetry(err) {
			if onRetry != nil {
				onRetry(attempt, err)
			}

			if attempt == maxRetries {
				logging.L(ctx).Error("Failed retrying connection",
					logging.IntAttr("attempts", maxRetries),
					logging.ErrAttr(err),
				)
				return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
			}

			delay := utils.ExponentialBackoff(attempt)

			logging.L(ctx).Error("Connection error, retrying after",
				logging.StringAttr("time", delay.String()),
				logging.ErrAttr(err),
			)
			fmt.Printf("Connection error, retrying after %.2fs\n", delay.Seconds())

			select {
			case <-ctx.Done():
				return fmt.Errorf("operation cancelled: %w", ctx.Err())
			case <-time.After(delay):
				continue
			}
		}

		return err
	}
	return err
}

func (a *App) prepareRemoteDatabaseList(
	server config.Server,
) (map[string]config.DBConnect, error) {
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

	fmt.Printf("Trying to connect to server %s...\n", server.Host)
	if err := runWithCtx(a.ctx, conn.Connect); err != nil {
		logging.L(a.ctx).Error(
			"Error connecting to server",
			logging.StringAttr("server", server.Host),
			logging.ErrAttr(err),
		)
		return nil, &ConnectionError{Addr: server.Host, Err: err}
	}

	defer func(conn *connect.Connect) {
		_ = conn.Close()
	}(conn)

	logging.L(a.ctx).Info("Trying to test connection to server")
	if err := runWithCtx(a.ctx, conn.TestConnection); err != nil {
		logging.L(a.ctx).Error("Error testing connection to server")
		return nil, &ConnectionError{Addr: server.Host, Err: err}
	}

	logging.L(a.ctx).Info("The connection is established")

	var rmt Remote
	rmt = remote.New(a.ctx, conn, server.ConfigPath)

	logging.L(a.ctx).Info("Trying to load remote config")
	if err := runWithCtx(a.ctx, rmt.Load); err != nil {
		logging.L(a.ctx).Error("Error load remote config")
		return nil, &ConnectionError{Addr: server.Host, Err: err}
	}

	logging.L(a.ctx).Info("Remote config loaded successfully")

	return rmt.Config(), nil
}

func (a *App) prepareDBConnect() map[string]config.DBConnect {
	connectDBs := make(map[string]config.DBConnect, len(a.cfg.Databases))
	for idx, database := range a.cfg.Databases {
		connectDBs[idx] = config.DBConnect{
			Server:   a.cfg.Servers[database.Server],
			Database: database,
		}
	}

	return connectDBs
}
