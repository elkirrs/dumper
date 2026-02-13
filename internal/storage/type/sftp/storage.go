package sftp

import (
	"context"
	"dumper/internal/connect"
	connectDomain "dumper/internal/domain/connect"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils/console"
	"dumper/pkg/utils/stream"
	"fmt"
	"io"
	"path/filepath"

	"github.com/pkg/sftp"
)

type SFTP struct {
	ctx     context.Context
	config  *storage.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *SFTP {
	return &SFTP{
		ctx:     ctx,
		config:  config,
		backend: "SFTP",
	}
}

func (s *SFTP) Save() error {
	connectDto := &connectDomain.Connect{
		Server:       s.config.Config.Host,
		Port:         s.config.Config.Port,
		Username:     s.config.Config.Username,
		Password:     s.config.Config.Password,
		PrivateKey:   s.config.Config.PrivateKey,
		Passphrase:   s.config.Config.Passphrase,
		IsPassphrase: true,
	}

	tClient := connect.NewApp(s.ctx, connectDto)

	if err := tClient.Connect(); err != nil {
		return &storage.UploadError{
			Backend: s.backend,
			Err:     fmt.Errorf("failed to connect target SFTP: %v", err),
		}
	}

	targetClient, err := sftp.NewClient(tClient.Client())
	if err != nil {
		return &storage.UploadError{
			Backend: s.backend,
			Err:     fmt.Errorf("failed to create target SFTP client: %v", err),
		}
	}
	defer targetClient.Close()

	targetPath := stream.TargetPath(s.config.Config.Dir, s.config.DumpName)
	dir := filepath.Dir(targetPath)

	if err := s.checkDirAccessible(targetClient, dir); err != nil {
		return &storage.UploadError{
			Backend: s.backend,
			Err:     err,
		}
	}

	pr, closeSSH, err := stream.SSHStreamer(s.ctx, s.config.Conn, s.config.DumpName, s.config.FileSize)

	if err != nil {
		return &storage.UploadError{
			Backend: s.backend,
			Err:     fmt.Errorf("failed to create SSH session: %v", err),
		}
	}

	defer closeSSH()

	dstFile, err := targetClient.Create(targetPath)
	if err != nil {
		return &storage.UploadError{
			Backend: s.backend,
			Err:     fmt.Errorf("failed to create remote file: %w", err),
		}
	}

	defer dstFile.Close()

	if _, err := io.Copy(dstFile, pr); err != nil {
		return &storage.UploadError{
			Backend: s.backend,
			Err:     fmt.Errorf("failed to upload to SFTP: %w", err),
		}
	}

	console.SafePrintln("[SFTP] Upload complete: %s", targetPath)
	return nil
}

func (s *SFTP) checkDirAccessible(client *sftp.Client, dir string) error {
	if err := client.MkdirAll(dir); err != nil {
		return fmt.Errorf("SFTP directory %s is not accessible: %w", dir, err)
	}

	if _, err := client.Stat(dir); err != nil {
		return fmt.Errorf("no permission to access SFTP directory %s: %w", dir, err)
	}

	return nil
}
