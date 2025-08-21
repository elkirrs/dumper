package app

import (
	"context"
	"dumper/internal/backup"
	"dumper/internal/command"
	_ "dumper/internal/command/mongo"
	_ "dumper/internal/command/mysql"
	_ "dumper/internal/command/postgres"
	"dumper/internal/config"
	"dumper/internal/connect"
	cmdCfg "dumper/internal/domain/command-config"
	_select "dumper/internal/select"
	t "dumper/internal/temr"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"errors"
	"fmt"
	"math/rand"
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

type DBInfo struct {
	Server   config.Server
	Database config.Database
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

	m := t.New()
	serverList, serverKeys := _select.SelectOptionList(a.cfg.Servers, "")
	m.SetList(serverKeys)
	m.SetTitle("Select server ")

	if err := runWithCtx(a.ctx, func() error { m.Run(); return nil }); err != nil {
		return err
	}

	serverName := m.GetSelect()
	serverKey := serverList[serverName]
	server := a.cfg.Servers[serverKey]

	logging.L(a.ctx).Info("Selected server", logging.StringAttr("server", serverKey))

	m.ClearList()
	logging.L(a.ctx).Info("Prepare database list")

	dbList, dbKeys := _select.SelectOptionList(a.cfg.Databases, serverKey)
	m.SetList(dbKeys)
	m.SetTitle("Select database ")

	if err := runWithCtx(a.ctx, func() error { m.Run(); return nil }); err != nil {
		return err
	}

	dbName := m.GetSelect()
	dbKey := dbList[dbName]
	db := a.cfg.Databases[dbKey]

	logging.L(a.ctx).Info("Selected database", logging.StringAttr("database", dbKey))

	err := withRetry(
		a.ctx, a.cfg.Settings.RetryConnect,
		func() error {
			return a.runBackup(server, db)
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
				logging.StringAttr("db", db.Name),
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

	serversDatabases := make(map[string][]DBInfo)

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

		serversDatabases[database.Server] = append(serversDatabases[database.Server], DBInfo{
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
		go func(dbInfos []DBInfo) {
			defer wg.Done()
			for _, dbInfo := range dbInfos {
				select {
				case <-a.ctx.Done():
					logging.L(a.ctx).Info("Backup cancelled by context")
					errCh <- fmt.Errorf("backup cancelled for database %s", dbInfo.Database.Name)
					return
				default:
					err := withRetry(
						a.ctx, a.cfg.Settings.RetryConnect,
						func() error {
							return a.runBackup(dbInfo.Server, dbInfo.Database)
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
								logging.StringAttr("db", dbInfo.Database.Name),
								logging.StringAttr("addr", connErr.Addr),
								logging.IntAttr("attempt", attempt),
								logging.ErrAttr(err),
							)
						},
					)

					if err != nil {
						errCh <- fmt.Errorf("backup failed for %s: %w", dbInfo.Database.Name, err)
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

func (a *App) runBackup(server config.Server, db config.Database) error {
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
	backupApp := backup.NewApp(a.ctx, conn, cmdStr, remotePath, a.cfg.Settings.DirDump, a.cfg.Settings.DumpLocation)

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

			delay := exponentialBackoff(attempt)

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
func exponentialBackoff(attempt int) time.Duration {
	base := time.Duration(1<<uint(attempt-1)) * time.Second
	jitterRange := int64(float64(base) * 0.6)
	jitter := time.Duration(rand.Int63n(jitterRange) - jitterRange/2)
	delay := base + jitter
	if delay < 0 {
		delay = 0
	}
	return delay
}
