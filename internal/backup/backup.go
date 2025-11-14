package backup

import (
	"context"
	backupLocalDirect "dumper/internal/backup/local-direct"
	backupLocalSSH "dumper/internal/backup/local-ssh"
	backupByServer "dumper/internal/backup/server"
	command "dumper/internal/command/database"
	"dumper/internal/connect"
	connecterror "dumper/internal/connect/connect-error"
	commandConfig "dumper/internal/domain/command-config"
	"dumper/internal/domain/config"
	dbConnect "dumper/internal/domain/config/db-connect"
	encryptConfigDomain "dumper/internal/domain/config/encrypt"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"fmt"
)

type Backup struct {
	ctx       context.Context
	cfg       *config.Config
	conn      *connect.Connect
	dbConnect dbConnect.DBConnect
	cmdConfig *commandConfig.Config
}

func NewApp(
	ctx context.Context,
	cfg *config.Config,
	dbConnect dbConnect.DBConnect,
) *Backup {
	return &Backup{
		ctx:       ctx,
		cfg:       cfg,
		dbConnect: dbConnect,
	}
}

func (b *Backup) Run() error {

	b.prepareBackupConfig()

	logging.L(b.ctx).Info("Prepare command for dump")

	cmdApp := command.NewApp(b.cmdConfig)
	cmdDB, err := cmdApp.GetCommand()

	if err != nil {
		logging.L(b.ctx).Error("failed generate command")
		return fmt.Errorf("failed generate command: %w", err)
	}

	b.cmdConfig.Command = cmdDB.Command
	b.cmdConfig.DumpName = cmdDB.DumpPath

	logging.L(b.ctx).Info("Prepare connection")

	conn := connect.New(
		b.dbConnect.Server.Host,
		b.dbConnect.Server.User,
		b.dbConnect.Server.GetPort(b.cfg.Settings.SrvPost),
		b.cfg.Settings.SSH.PrivateKey,
		b.cfg.Settings.SSH.Passphrase,
		b.dbConnect.Server.Password,
		*b.cfg.Settings.SSH.IsPassphrase,
	)
	b.conn = conn

	fmt.Printf("Trying to connect to server %s...\n", b.dbConnect.Server.Host)
	if err := utils.RunWithCtx(b.ctx, conn.Connect); err != nil {
		logging.L(b.ctx).Error(
			"Error connecting to server",
			logging.StringAttr("server", b.dbConnect.Server.Host),
			logging.ErrAttr(err),
		)
		return &connecterror.ConnectError{
			Addr: b.dbConnect.Server.Host,
			Err:  err,
		}
	}

	defer func(conn *connect.Connect) {
		_ = conn.Close()
	}(conn)

	logging.L(b.ctx).Info("Trying to test connection to server")
	if err := utils.RunWithCtx(b.ctx, conn.TestConnection); err != nil {
		logging.L(b.ctx).Error("Error testing connection to server")
		return &connecterror.ConnectError{
			Addr: b.dbConnect.Server.Host,
			Err:  err,
		}
	}
	logging.L(b.ctx).Info("The connection is established")

	logging.L(b.ctx).Info("Preparing for backup creation")

	if err := utils.RunWithCtx(b.ctx, b.backup); err != nil {
		logging.L(b.ctx).Error("Error creating backup database")
		return err
	}
	logging.L(b.ctx).Info("Backup was successfully created and downloaded")

	if b.cfg.Settings.DirArchived != "" {
		logging.L(b.ctx).Info("Search for old backups")
		dbNamePrefix := fmt.Sprintf("%s_%s",
			b.dbConnect.Server.GetName(),
			b.dbConnect.Database.GetName(),
		)

		if err := utils.RunWithCtx(b.ctx, func() error {
			return utils.ArchivedLocalFile(dbNamePrefix, b.cmdConfig.DumpName, b.cfg.Settings.DirDump, b.cfg.Settings.DirArchived)
		}); err != nil {
			logging.L(b.ctx).Error("Error archiving old backups")
			return err
		}

		logging.L(b.ctx).Info("Archived old backups", logging.StringAttr("path", b.cfg.Settings.DirArchived))
	}

	return nil
}

func (b *Backup) backup() error {
	switch b.cfg.Settings.DumpLocation {
	case "server":
		byServer := backupByServer.NewApp(b.ctx, b.conn, b.cmdConfig)
		return byServer.Run()
	case "local-ssh":
		localSSH := backupLocalSSH.NewApp(b.ctx, b.conn, b.cmdConfig)
		return localSSH.Run()
	case "local-direct":
		localDirect := backupLocalDirect.NewApp(b.ctx, b.conn, b.cmdConfig)
		return localDirect.Run()
	default:
		logging.L(b.ctx).Error(
			"Unsupported backup dump location",
			logging.StringAttr("location", b.cfg.Settings.DumpLocation),
		)
		return fmt.Errorf("unsupported backup dump location: %s", b.cfg.Settings.DumpLocation)
	}
}

func (b *Backup) backupByLocalSSH() error {
	panic("not implement")
}

func (b *Backup) backupLocalDirect() error {
	panic("not implement")
}

func (b *Backup) prepareBackupConfig() {
	logging.L(b.ctx).Info("Prepare command config")

	dataFormat := utils.TemplateData{
		Server:   b.dbConnect.Server.GetName(),
		Database: b.dbConnect.Database.GetName(),
		Template: b.cfg.Settings.Template,
	}
	nameFile := utils.GetTemplateFileName(dataFormat)
	fullPath := utils.GetFullPath(b.cfg.Settings.DirRemote, nameFile)

	b.cmdConfig = &commandConfig.Config{
		Database: commandConfig.Database{
			User:     b.dbConnect.Database.User,
			Password: b.dbConnect.Database.Password,
			Name:     b.dbConnect.Database.GetName(),
			Port:     b.dbConnect.Database.GetPort(b.cfg.Settings.DBPort),
			Format:   b.dbConnect.Database.GetFormat(b.cfg.Settings.DumpFormat),
			Driver:   b.dbConnect.Database.GetDriver(b.cfg.Settings.Driver),
			Options:  b.dbConnect.Database.Options,
		},
		Server: commandConfig.Server{
			Host: b.dbConnect.Server.Host,
			Port: b.dbConnect.Server.Port,
			Key:  b.dbConnect.Server.GetPrivateKey(b.cfg.Settings.SSH.PrivateKey),
		},
		Storages:      b.dbConnect.Storages,
		DumpLocation:  b.cfg.Settings.DumpLocation,
		Archive:       b.dbConnect.Database.IsArchive(*b.cfg.Settings.Archive),
		DumpDirLocal:  b.cfg.Settings.DirDump,
		DumpName:      fullPath,
		DumpDirRemote: b.cfg.Settings.DirRemote,
		RemoveBackup:  *b.cfg.Settings.RemoveDump,
		Encrypt: encryptConfigDomain.Encrypt{
			Type:     b.dbConnect.Database.GetEncryptType(b.cfg.Settings.Encrypt.Type),
			Password: b.dbConnect.Database.GetEncryptPass(b.cfg.Settings.Encrypt.Password),
		},
		MaxParallelDownload: b.cfg.Settings.MaxParallelDownload,
	}
}
