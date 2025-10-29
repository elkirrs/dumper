package server

import (
	"context"
	encryptCommand "dumper/internal/command/encrypt"
	"dumper/internal/connect"
	backupDomain "dumper/internal/domain/backup"
	commandConfig "dumper/internal/domain/command-config"
	encryptDomain "dumper/internal/domain/encrypt"
	"dumper/internal/download"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"fmt"
	"time"
)

type BackupServer struct {
	ctx    context.Context
	conn   *connect.Connect
	config *commandConfig.Config
}

func NewApp(
	ctx context.Context,
	conn *connect.Connect,
	config *commandConfig.Config,
) *BackupServer {
	return &BackupServer{
		ctx:    ctx,
		conn:   conn,
		config: config,
	}
}

func (b *BackupServer) Run() error {

	isRemoveDump := b.config.RemoveBackup
	checkCmd := fmt.Sprintf("test -f %s", b.config.DumpName)

	logging.L(b.ctx).Info(
		"Run command found backup in server with name",
		logging.StringAttr("name", b.config.DumpName),
	)

	downloadApp := download.NewApp(b.ctx, b.conn, b.config)

	if msg, err := b.conn.RunCommand(checkCmd); err == nil {
		logging.L(b.ctx).Info(
			"Dump already exists on server",
			logging.StringAttr("name", b.config.DumpName),
			logging.StringAttr("msg", msg),
		)

		fmt.Println("Dump already exists on server:", b.config.DumpName)
		isRemoveDump = false
	} else {
		stop := make(chan struct{})
		dumpCreateTimeNow := time.Now()

		logging.L(b.ctx).Info("File dump name", logging.StringAttr("name", b.config.DumpName))
		fmt.Println("File dump name:", b.config.DumpName)

		go utils.Spinner(stop)

		if msg, err := b.conn.RunCommand(b.config.Command); err != nil {
			logging.L(b.ctx).Error(
				"Failed to create dump",
				logging.StringAttr("msg", msg),
				logging.ErrAttr(err),
			)
			return fmt.Errorf("failed to create dump: %v", err)
		}

		close(stop)

		elapsed := time.Since(dumpCreateTimeNow)

		totalSize, err := downloadApp.FileSize()
		if err != nil {
			return err
		}

		fmt.Printf("\rDump created successfully in %.2f sec\n", elapsed.Seconds())
		fmt.Printf("\rFile dump size: %s [%d bytes]\n", utils.FormatBytes(totalSize), totalSize)

		dumpCreateTimeSec := fmt.Sprintf("%.2f sec", elapsed.Seconds())
		logging.L(b.ctx).Info(
			"The dump was successfully created",
			logging.StringAttr("time", dumpCreateTimeSec),
			logging.Int64Attr("size", totalSize),
		)
	}

	var fileList []*backupDomain.FileRemoveList
	fileList = append(fileList, &backupDomain.FileRemoveList{
		Name:     b.config.DumpName,
		IsRemove: isRemoveDump,
	})

	if b.config.Encrypt.Type != "" {
		encOpts := encryptDomain.Options{
			FilePath: b.config.DumpName,
			Password: b.config.Encrypt.Password,
			Type:     b.config.Encrypt.Type,
			Crypt:    "encrypt",
		}
		encryptApp := encryptCommand.NewApp(&encOpts)
		encryptCmd, _ := encryptApp.Generate()

		if msg, err := b.conn.RunCommand(encryptCmd.CMD); err != nil {
			logging.L(b.ctx).Error(
				"Failed to encrypt dump",
				logging.StringAttr("msg", msg),
				logging.ErrAttr(err),
			)
			return fmt.Errorf("failed to encrypt dump on server: %v, msg: %s", err, msg)
		}

		logging.L(b.ctx).Info("File dump encrypted successfully")
		fmt.Println("File dump encrypted successfully")

		b.config.DumpName = encryptCmd.Name
		fileList = append(fileList, &backupDomain.FileRemoveList{
			Name:     encryptCmd.Name,
			IsRemove: true,
		})
	}

	logging.L(b.ctx).Info("Downloading dump", logging.StringAttr("name", b.config.DumpName))
	dumpDownloadTimeNow := time.Now()
	if err := downloadApp.DownloadFile(); err != nil {
		logging.L(b.ctx).Error(
			"Failed to download dump",
			logging.ErrAttr(err),
		)
		return fmt.Errorf("failed to download dump: %v", err)
	}

	dumpDownloadTimeSec := fmt.Sprintf("%.2f sec", time.Since(dumpDownloadTimeNow).Seconds())

	logging.L(b.ctx).Info(
		"The dump was successfully downloaded",
		logging.StringAttr("time", dumpDownloadTimeSec),
	)

	for _, file := range fileList {

		if !file.IsRemove {
			continue
		}

		logging.L(b.ctx).Info("Removing dump on server")
		fmt.Println("Removing dump from server:", file.Name)
		if msg, err := b.conn.RunCommand(fmt.Sprintf("rm -f %s", file.Name)); err != nil {
			logging.L(b.ctx).Error(
				"Failed to remove dump on server",
				logging.StringAttr("msg", msg),
			)
			return fmt.Errorf("failed to delete dump on server: %v", err)
		}
	}

	logging.L(b.ctx).Info("The dump was successfully deleted on server")

	return nil
}
