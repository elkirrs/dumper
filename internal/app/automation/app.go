package automation

import (
	"context"
	"dumper/internal/backup"
	connecterror "dumper/internal/connect/connect-error"
	"dumper/internal/domain/app"
	cfg "dumper/internal/domain/config"
	dbConnect "dumper/internal/domain/config/db-connect"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Automation struct {
	ctx context.Context
	cfg *cfg.Config
	env *app.Env
}

func NewApp(
	ctx context.Context,
	cfg *cfg.Config,
	env *app.Env,
) *Automation {
	return &Automation{
		ctx: ctx,
		cfg: cfg,
		env: env,
	}
}

func (a *Automation) Run() error {
	logging.L(a.ctx).Info("Prepare data for create dumps")

	dbList := strings.Split(a.env.DbName, ",")
	countDBs := len(dbList)

	serversDatabases := make(map[string][]dbConnect.DBConnect)

	for _, dbName := range dbList {
		database, ok := a.cfg.Databases[dbName]
		if !ok {
			fmt.Printf("Database %s not found\n", dbName)
			logging.L(a.ctx).Warn("Database not found", logging.StringAttr("name", dbName))
			countDBs--
			continue
		}

		srv, ok := a.cfg.Servers[database.Server]
		if !ok {
			fmt.Printf("Server %s not found\n", database.Server)
			logging.L(a.ctx).Warn("Server not found", logging.StringAttr("name", database.Server))
			countDBs--
			continue
		}

		serversDatabases[database.Server] = append(serversDatabases[database.Server], dbConnect.DBConnect{
			Server:   srv,
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
		go func(connectDBs []dbConnect.DBConnect) {
			defer wg.Done()
			for _, dbConn := range connectDBs {
				select {
				case <-a.ctx.Done():
					logging.L(a.ctx).Info("Backup cancelled by context")
					errCh <- fmt.Errorf("backup cancelled for database %s", dbConn.Database.Name)
					return
				default:
					backupApp := backup.NewApp(a.ctx, a.cfg, dbConn)

					err := utils.WithRetry(
						a.ctx, a.cfg.Settings.RetryConnect,
						func() error {
							return backupApp.Run()
						},
						func(err error) bool {
							var connErr *connecterror.ConnectError
							return errors.As(err, &connErr)
						},
						func(attempt int, err error) {
							var connErr *connecterror.ConnectError
							_ = errors.As(err, &connErr)
							logging.L(a.ctx).Warn(
								"Connection error, retrying",
								logging.StringAttr("db", dbConn.Database.Name),
								logging.StringAttr("addr", connErr.Addr),
								logging.IntAttr("attempt", attempt),
								logging.ErrAttr(err),
							)
						},
					)

					if err != nil {
						errCh <- fmt.Errorf("backup failed for %s: %w", dbConn.Database.Name, err)
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
