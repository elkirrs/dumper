package automation

import (
	"context"
	"dumper/internal/backup"
	"dumper/internal/connect"
	connecterror "dumper/internal/connect/connect-error"
	"dumper/internal/domain/app"
	cfg "dumper/internal/domain/config"
	dbConnect "dumper/internal/domain/config/db-connect"
	"dumper/internal/domain/config/storage"
	connectDomain "dumper/internal/domain/connect"
	"dumper/pkg/logging"
	"dumper/pkg/utils/retry"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Automation struct {
	ctx context.Context
	cfg *cfg.Config
	env *app.Flags
}

func NewApp(
	ctx context.Context,
	cfg *cfg.Config,
	env *app.Flags,
) *Automation {
	return &Automation{
		ctx: ctx,
		cfg: cfg,
		env: env,
	}
}

func (m *Automation) Run() error {
	logging.L(m.ctx).Info("Prepare data for create dumps")

	dbList := strings.Split(m.env.DbNameList, ",")
	countDBs := len(dbList)

	serversDatabases := make(map[string][]dbConnect.DBConnect)
	var dataDBConnect map[string]dbConnect.DBConnect
	dataDBConnect = m.prepareDBConnect()

	for _, dbName := range dbList {
		dbC, ok := dataDBConnect[dbName]
		if !ok {
			fmt.Printf("Database %s not found\n", dbName)
			logging.L(m.ctx).Warn("Database not found", logging.StringAttr("name", dbName))
			countDBs--
			continue
		}

		serversDatabases[dbC.Database.Server] = append(serversDatabases[dbC.Database.Server], dbC)

		if countDBs == 0 {
			logging.L(m.ctx).Error("Database and server no key matches check the configuration file")
			return errors.New("database and server no key matches check the configuration file")
		}
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error, len(dbList))

	for _, dbInfoList := range serversDatabases {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, dbConn := range dbInfoList {
				select {
				case <-m.ctx.Done():
					logging.L(m.ctx).Info("Backup cancelled by context")
					errCh <- fmt.Errorf("backup cancelled for database %s", dbConn.Database.Name)
					return
				default:
					connectDto := &connectDomain.Connect{
						Server:       dbConn.Server.Host,
						Port:         dbConn.Server.GetPort(&m.cfg.Settings.SrvPost),
						Username:     dbConn.Server.User,
						Password:     dbConn.Server.GetPassword(&m.cfg.Settings.SSH.Password),
						PrivateKey:   dbConn.Server.GetPrivateKey(&m.cfg.Settings.SSH.PrivateKey),
						Passphrase:   dbConn.Server.GetPassphrase(&m.cfg.Settings.SSH.Passphrase),
						IsPassphrase: dbConn.Server.GetIsPassphrase(*m.cfg.Settings.SSH.IsPassphrase),
					}
					connectApp := connect.NewApp(m.ctx, connectDto)
					backupApp := backup.NewApp(m.ctx, m.cfg, dbConn, connectApp)

					err := retry.WithRetry(
						m.ctx, m.cfg.Settings.RetryConnect,
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
							logging.L(m.ctx).Warn(
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
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	logging.L(m.ctx).Info("All requested database backups are done")

	return nil
}

func (m *Automation) prepareDBConnect() map[string]dbConnect.DBConnect {
	connectDBs := make(map[string]dbConnect.DBConnect, len(m.cfg.Databases))
	for idx, database := range m.cfg.Databases {
		storageList := database.GetStorages(&m.cfg.Settings.Storages)
		connectDBs[idx] = dbConnect.DBConnect{
			Server:   m.cfg.Servers[database.Server],
			Database: database,
			Storages: m.prepareStorages(storageList),
		}
	}

	return connectDBs
}

func (m *Automation) prepareStorages(list []string) map[string]storage.Storage {
	storages := make(map[string]storage.Storage, len(list))
	for _, storageType := range list {
		st := m.cfg.Storages[storageType]
		st.PrivateKey = st.GetPrivateKey(m.cfg.Settings.SSH.PrivateKey)
		storages[storageType] = st
	}
	return storages
}
