package sftp

import (
	"context"
	"dumper/internal/connect"
	connectDomain "dumper/internal/domain/connect"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils"
	"fmt"
	"io"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTP struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *SFTP {
	return &SFTP{
		ctx:    ctx,
		config: config,
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
		return fmt.Errorf("failed to connect target SFTP: %w", err)
	}

	targetClient, err := sftp.NewClient(tClient.Client())
	if err != nil {
		return fmt.Errorf("failed to create target SFTP client: %v", err)
	}
	defer func(targetClient *sftp.Client) {
		_ = targetClient.Close()
	}(targetClient)

	session, err := s.config.Conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := session.Start(fmt.Sprintf("cat %s", s.config.DumpName)); err != nil {
		return fmt.Errorf("failed to start remote command: %v", err)
	}

	targetPath := filepath.Join(s.config.Config.Dir, filepath.Base(s.config.DumpName))
	if err := targetClient.MkdirAll(filepath.Dir(targetPath)); err != nil {
		return fmt.Errorf("failed to create remote directory: %v", err)
	}

	dstFile, err := targetClient.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer func(dstFile *sftp.File) {
		_ = dstFile.Close()
	}(dstFile)

	pr := utils.StreamToPipe(s.ctx, stdout, s.config.FileSize)

	if _, err := io.Copy(dstFile, pr); err != nil {
		return fmt.Errorf("failed to upload to SFTP: %v", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %v", err)
	}

	utils.SafePrintln("[SFTP] Upload complete: %s", targetPath)
	return nil
}
