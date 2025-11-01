package manual

import (
	"context"
	"dumper/internal/backup"
	remote "dumper/internal/config/remote"
	"dumper/internal/connect"
	connecterror "dumper/internal/connect/connect-error"
	"dumper/internal/domain/app"
	cfg "dumper/internal/domain/config"
	dbConnect "dumper/internal/domain/config/db-connect"
	"dumper/internal/domain/config/server"
	"dumper/internal/domain/config/storage"
	_select "dumper/internal/select"
	t "dumper/internal/temr"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"errors"
	"fmt"
)

type Manual struct {
	ctx context.Context
	cfg *cfg.Config
	env *app.Flags
}

type Remote interface {
	Load() error
	Config() map[string]dbConnect.DBConnect
}

func NewApp(
	ctx context.Context,
	cfg *cfg.Config,
	env *app.Flags,
) *Manual {
	return &Manual{
		ctx: ctx,
		cfg: cfg,
		env: env,
	}
}

func (m *Manual) Run() error {
	logging.L(m.ctx).Info("Prepare server list")
	var serverName string
	var err error

	term := t.New()
	serverList, serverKeys := _select.SelectOptionList(m.cfg.Servers, "")

	if len(m.cfg.Servers) > 1 {
		term.SetList(serverKeys)
		term.SetTitle("Select server ")

		if err = utils.RunWithCtx(m.ctx, func() error { term.Run(); return nil }); err != nil {
			return err
		}

		serverName = term.GetSelect()
	} else {
		serverName = serverKeys[0]
		fmt.Println("\033[32m" + "\U00002714 " + serverName + "\033[0m")
	}

	serverKey := serverList[serverName]
	srv := m.cfg.Servers[serverKey]

	logging.L(m.ctx).Info("Selected server", logging.StringAttr("server", serverKey))

	term.ClearList()
	logging.L(m.ctx).Info("Prepare database list")

	var dataDBConnect map[string]dbConnect.DBConnect

	if srv.ConfigPath != "" {
		dataDBConnect, err = m.prepareRemoteDatabaseList(srv)
		serverKey = ""
		if err != nil {
			return err
		}
	} else {
		dataDBConnect = m.prepareDBConnect()
	}

	dbList, dbKeys := _select.OptionDataBaseList(dataDBConnect, serverKey)

	term.SetList(dbKeys)
	term.SetTitle("Select database ")

	if err = utils.RunWithCtx(m.ctx, func() error { term.Run(); return nil }); err != nil {
		return err
	}

	dbName := term.GetSelect()
	dbKey := dbList[dbName]
	dbConn := dataDBConnect[dbKey]

	logging.L(m.ctx).Info("Selected database", logging.StringAttr("database", dbKey))

	backupApp := backup.NewApp(m.ctx, m.cfg, dbConn)

	err = utils.WithRetry(
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
		return err
	}

	return nil
}

func (m *Manual) prepareRemoteDatabaseList(
	server server.Server,
) (map[string]dbConnect.DBConnect, error) {
	logging.L(m.ctx).Info("Prepare connection")

	conn := connect.New(
		server.Host,
		server.User,
		server.GetPort(m.cfg.Settings.SrvPost),
		m.cfg.Settings.SSH.PrivateKey,
		server.SSHKey,
		m.cfg.Settings.SSH.Passphrase,
		server.Password,
		*m.cfg.Settings.SSH.IsPassphrase,
	)

	fmt.Printf("Trying to connect to server %s...\n", server.Host)
	if err := utils.RunWithCtx(m.ctx, conn.Connect); err != nil {
		logging.L(m.ctx).Error(
			"Error connecting to server",
			logging.StringAttr("server", server.Host),
			logging.ErrAttr(err),
		)
		return nil, &connecterror.ConnectError{Addr: server.Host, Err: err}
	}

	defer func(conn *connect.Connect) {
		_ = conn.Close()
	}(conn)

	logging.L(m.ctx).Info("Trying to test connection to server")
	if err := utils.RunWithCtx(m.ctx, conn.TestConnection); err != nil {
		logging.L(m.ctx).Error("Error testing connection to server")
		return nil, &connecterror.ConnectError{Addr: server.Host, Err: err}
	}

	logging.L(m.ctx).Info("The connection is established")

	var rmt Remote
	rmt = remote.New(m.ctx, conn, server.ConfigPath)

	logging.L(m.ctx).Info("Trying to load remote config")
	if err := utils.RunWithCtx(m.ctx, rmt.Load); err != nil {
		logging.L(m.ctx).Error("Error load remote config")
		return nil, &connecterror.ConnectError{Addr: server.Host, Err: err}
	}

	logging.L(m.ctx).Info("Remote config loaded successfully")

	return rmt.Config(), nil
}

func (m *Manual) prepareDBConnect() map[string]dbConnect.DBConnect {
	connectDBs := make(map[string]dbConnect.DBConnect, len(m.cfg.Databases))
	for idx, database := range m.cfg.Databases {
		connectDBs[idx] = dbConnect.DBConnect{
			Server:   m.cfg.Servers[database.Server],
			Database: database,
			Storages: m.prepareStorages(database.Storages),
		}
	}

	return connectDBs
}

func (m *Manual) prepareStorages(list []string) map[string]storage.ListStorages {
	storages := make(map[string]storage.ListStorages, len(list))
	for _, storageType := range list {
		storages[storageType] = storage.ListStorages{
			Type:    storageType,
			Configs: m.cfg.Storages[storageType],
		}
	}

	if len(storages) == 0 {
		storages["local"] = storage.ListStorages{
			Type: "local",
			Configs: storage.Storage{
				Type: "local",
				Dir:  m.cfg.Settings.DirDump,
			},
		}
	}

	return storages
}
