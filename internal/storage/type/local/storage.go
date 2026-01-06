package local

import (
	"context"
	storageDomain "dumper/internal/domain/storage"
	"dumper/pkg/utils/console"
	"dumper/pkg/utils/stream"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	ctx     context.Context
	config  *storageDomain.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storageDomain.Config,
) *Local {
	return &Local{
		ctx:     ctx,
		config:  config,
		backend: "Local",
	}
}

func (l *Local) Save() error {
	localPath := stream.TargetPath(l.config.Config.Dir, l.config.DumpName)

	dir := filepath.Dir(localPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return &storageDomain.UploadError{
			Backend: l.backend,
			Err:     fmt.Errorf("failed to create local directory: %v", err),
		}
	}

	outFile, err := os.Create(localPath)
	if err != nil {
		return &storageDomain.UploadError{
			Backend: l.backend,
			Err:     fmt.Errorf("failed to create local file: %v", err),
		}
	}
	defer outFile.Close()

	pr, closeSSH, err := stream.SSHStreamer(l.ctx, l.config.Conn, l.config.DumpName, l.config.FileSize)

	if err != nil {
		return &storageDomain.UploadError{
			Backend: l.backend,
			Err:     fmt.Errorf("failed to create SSH session: %v", err),
		}
	}

	defer closeSSH()

	if _, err := io.Copy(outFile, pr); err != nil {
		outFile.Close()
		return &storageDomain.UploadError{
			Backend: l.backend,
			Err:     fmt.Errorf("failed to write to local file: %v", err),
		}
	}

	console.SafePrintln("[Local] Upload complete: %s", localPath)
	return nil
}
