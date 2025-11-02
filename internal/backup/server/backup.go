package server

import (
	"context"
	encryptCommand "dumper/internal/command/encrypt"
	"dumper/internal/connect"
	backupDomain "dumper/internal/domain/backup"
	commandConfig "dumper/internal/domain/command-config"
	configStorageDomain "dumper/internal/domain/config/storage"
	encryptDomain "dumper/internal/domain/encrypt"
	storageDomain "dumper/internal/domain/storage"
	"dumper/internal/storage"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"fmt"
	"strings"
	"sync"
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

		fmt.Printf("\rFile dump size: %s [%d bytes]\n", utils.FormatBytes(totalSize), totalSize)

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

		totalSize, err = b.FileSize()
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

	wg := &sync.WaitGroup{}
	errCh := make(chan error, len(b.config.Storages))
	var countStorage int8
	sem := make(chan struct{}, b.config.MaxParallelDownload)
	var totalAll int64

	totalAll = totalSize * int64(len(b.config.Storages))
	globalProgress := utils.GlobalProgress(totalAll)

	for _, storageItem := range b.config.Storages {
		countStorage++
		wg.Add(1)
		go func(item configStorageDomain.ListStorages) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			select {
			case <-b.ctx.Done():
				logging.L(b.ctx).Info("download cancelled by context")
				errCh <- fmt.Errorf("download cancelled for storage %s", storageItem.Type)
				return
			default:
				storageConfig := storageDomain.Config{
					Type:     storageItem.Type,
					DumpName: b.config.DumpName,
					FileSize: totalSize,
					Conn:     b.conn,
					Config:   storageItem.Configs,
				}

				ctx := context.WithValue(b.ctx, "globalProgress", globalProgress)
				storageApp := storage.NewApp(ctx, &storageConfig)

				if err := storageApp.Save(); err != nil {
					logging.L(b.ctx).Error(
						"Failed to download dump",
						logging.ErrAttr(err),
					)
					errCh <- fmt.Errorf("failed to download dump %s: %w", b.config.DumpName, err)
					return
				}

				dumpDownloadTimeSec := fmt.Sprintf("%.2f sec", time.Since(dumpDownloadTimeNow).Seconds())

				logging.L(b.ctx).Info(
					"The dump was successfully downloaded",
					logging.StringAttr("time", dumpDownloadTimeSec),
					logging.StringAttr("storage", storageItem.Type),
				)
			}
		}(storageItem)
	}

	wg.Wait()
	close(errCh)

	utils.Progress(totalAll, totalAll)

	for err := range errCh {
		countStorage--
		if err != nil {
			fmt.Println(err)
		}
		if countStorage == 0 {
			return fmt.Errorf("failed to download dump")
		}
	}

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

func (b *BackupServer) FileSize() (int64, error) {
	sizeOutput, err := b.conn.RunCommand(fmt.Sprintf("stat -c %%s %s", b.config.DumpName))

	var totalSize int64

	if err != nil {
		return totalSize, fmt.Errorf("failed to get file size: %v", err)
	}
	sizeOutput = strings.TrimSpace(sizeOutput)

	_, err = fmt.Sscanf(sizeOutput, "%d", &totalSize)
	if err != nil {
		return totalSize, err
	}

	return totalSize, nil
}
