package server

import (
	"context"
	encryptCommand "dumper/internal/command/encrypt"
	"dumper/internal/connect"
	backupDomain "dumper/internal/domain/backup"
	commandConfig "dumper/internal/domain/command-config"
	encryptDomain "dumper/internal/domain/encrypt"
	"dumper/pkg/logging"
	"dumper/pkg/utils/format"
	"dumper/pkg/utils/spiner"
	"fmt"
	"strings"
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

	mkdirCmd := fmt.Sprintf("mkdir -p %s", b.config.DumpDirRemote)
	if msg, err := b.conn.RunCommand(mkdirCmd); err == nil {
		logging.L(b.ctx).Info(
			"Created backup directory",
			logging.StringAttr("dir", b.config.DumpDirRemote),
			logging.StringAttr("msg", msg),
		)
	} else {
		logging.L(b.ctx).Error(
			"error while creating backup directory",
			logging.StringAttr("dir", b.config.DumpDirRemote),
			logging.StringAttr("msg", msg),
			logging.ErrAttr(err),
		)
		return err
	}

	isRemoveDump := b.config.RemoveBackup
	checkCmd := fmt.Sprintf("test -f %s", b.config.DumpName)

	logging.L(b.ctx).Info(
		"Run command found backup in server with name",
		logging.StringAttr("name", b.config.DumpName),
	)

	var totalSize int64

	if msg, err := b.conn.RunCommand(checkCmd); err == nil {
		logging.L(b.ctx).Info(
			"Dump already exists on server",
			logging.StringAttr("name", b.config.DumpName),
			logging.StringAttr("msg", msg),
		)

		fmt.Println("Dump already exists on server:", b.config.DumpName)
		isRemoveDump = false

		totalSize, err = b.FileSize()
		if err != nil {
			return err
		}

		fmt.Printf("\rFile dump size: %s [%d bytes]\n", format.FormatBytes(totalSize), totalSize)

	} else {
		stop := make(chan struct{})
		dumpCreateTimeNow := time.Now()

		logging.L(b.ctx).Info("File dump name", logging.StringAttr("name", b.config.DumpName))
		fmt.Println("File dump name:", b.config.DumpName)

		go spiner.Spinner(stop)

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

		totalSize, err = b.FileSize()
		if err != nil {
			return err
		}

		fmt.Printf("\rDump created successfully in %.2f sec\n", elapsed.Seconds())
		fmt.Printf("\rFile dump size: %s [%d bytes]\n", format.FormatBytes(totalSize), totalSize)

		dumpCreateTimeSec := fmt.Sprintf("%.2f sec", elapsed.Seconds())
		logging.L(b.ctx).Info(
			"The dump was successfully created",
			logging.StringAttr("time", dumpCreateTimeSec),
			logging.Int64Attr("size", totalSize),
		)
	}

	var fileList []backupDomain.FileRemoveList
	fileList = append(fileList, backupDomain.FileRemoveList{
		Name:     b.config.DumpName,
		IsRemove: isRemoveDump,
	})

	if b.config.Encrypt.Type != "" && b.config.Encrypt.Password != "" && *b.config.Encrypt.Enabled {
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
		fileList = append(fileList, backupDomain.FileRemoveList{
			Name:     encryptCmd.Name,
			IsRemove: true,
		})
	}

	b.config.FileRemoveList = fileList
	b.config.FileSize = totalSize

	return nil
}

func (b *BackupServer) FileSize() (int64, error) {
	sizeOutput, err := b.conn.RunCommand(fmt.Sprintf("stat -c %%s %s", b.config.DumpName))

	var totalSize int64

	if err != nil {
		return totalSize, fmt.Errorf("failed to get file size. path: %s err: %v", b.config.DumpName, err)
	}
	sizeOutput = strings.TrimSpace(sizeOutput)

	_, err = fmt.Sscanf(sizeOutput, "%d", &totalSize)
	if err != nil {
		return totalSize, err
	}

	return totalSize, nil
}
