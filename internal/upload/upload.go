package upload

import (
	"context"
	"dumper/internal/connect"
	commandConfig "dumper/internal/domain/command-config"
	storageDomain "dumper/internal/domain/storage"
	"dumper/internal/storage"
	"dumper/pkg/logging"
	"dumper/pkg/utils/progress"
	"fmt"
	"sync"
	"time"
)

type Upload struct {
	ctx    context.Context
	conn   *connect.Connect
	config *commandConfig.Config
}

func New(
	ctx context.Context,
	conn *connect.Connect,
	config *commandConfig.Config,
) *Upload {
	return &Upload{
		ctx:    ctx,
		conn:   conn,
		config: config,
	}
}

func (u *Upload) Uploading() error {
	logging.L(u.ctx).Info("Downloading dump", logging.StringAttr("name", u.config.DumpName))

	var countStorage int8
	var totalAll int64
	var totalDone int64

	totalSize := u.config.FileSize
	dumpDownloadTimeNow := time.Now()
	sem := make(chan struct{}, u.config.MaxParallelDownload)
	totalAll = totalSize * int64(len(u.config.Storages))

	wg := &sync.WaitGroup{}
	errCh := make(chan error, len(u.config.Storages))

	globalProgress := progress.GlobalProgress(totalAll)

	for _, storageItem := range u.config.Storages {
		countStorage++
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			select {
			case <-u.ctx.Done():
				logging.L(u.ctx).Info("download cancelled by context")
				errCh <- fmt.Errorf("download cancelled for storage %s", storageItem.Type)
				return
			default:
				storageConfig := storageDomain.Config{
					Type:     storageItem.Type,
					DumpName: u.config.DumpName,
					FileSize: totalSize,
					Conn:     u.conn,
					Config:   storageItem,
				}

				ctx := context.WithValue(u.ctx, "globalProgress", globalProgress)
				storageApp := storage.NewApp(ctx, &storageConfig)

				if err := storageApp.Save(); err != nil {
					logging.L(u.ctx).Error(
						"Failed to download dump",
						logging.ErrAttr(err),
					)
					errCh <- fmt.Errorf("failed to download dump %s: %w", u.config.DumpName, err)
					return
				}

				dumpDownloadTimeSec := fmt.Sprintf("%.2f sec", time.Since(dumpDownloadTimeNow).Seconds())

				logging.L(u.ctx).Info(
					"The dump was successfully downloaded",
					logging.StringAttr("time", dumpDownloadTimeSec),
					logging.StringAttr("storage", storageItem.Type),
				)
				totalDone += totalSize
			}
		}()
	}

	wg.Wait()
	close(errCh)

	progress.Progress(totalDone, totalAll)

	for err := range errCh {
		countStorage--
		if err != nil {
			fmt.Println(err)
		}
		if countStorage == 0 {
			return fmt.Errorf("failed to download dump")
		}
	}

	for _, file := range u.config.FileRemoveList {
		if !file.IsRemove {
			continue
		}

		logging.L(u.ctx).Info(
			"Removing dump on server",
			logging.StringAttr("file", file.Name),
		)
		fmt.Println("Removing dump from server:", file.Name)
		if msg, err := u.conn.RunCommand(fmt.Sprintf("rm -f %s", file.Name)); err != nil {
			logging.L(u.ctx).Error(
				"Failed to remove dump on server",
				logging.StringAttr("file", file.Name),
				logging.StringAttr("msg", msg),
			)
			return fmt.Errorf("failed to delete dump on server: %v", err)
		}
	}

	logging.L(u.ctx).Info("The dump was successfully deleted on server")

	return nil
}
