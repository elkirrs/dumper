package download

import (
	"context"
	"dumper/internal/connect"
	commandConfig "dumper/internal/domain/command-config"
	"dumper/pkg/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Download struct {
	ctx    context.Context
	conn   *connect.Connect
	config *commandConfig.Config
}

func NewApp(
	ctx context.Context,
	conn *connect.Connect,
	config *commandConfig.Config,
) *Download {
	return &Download{
		ctx:    ctx,
		conn:   conn,
		config: config,
	}
}

func (d *Download) DownloadFile() error {
	localPath := filepath.Join(d.config.DumpDirLocal, filepath.Base(d.config.DumpName))

	var totalSize int64

	totalSize, err := d.FileSize()
	if err != nil {
		return err
	}

	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}

	defer func(outFile *os.File) {
		_ = outFile.Close()
		return
	}(outFile)

	session, err := d.conn.NewSession()
	if err != nil {
		return err
	}

	defer func(session *ssh.Session) {
		_ = session.Close()
		return
	}(session)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	if err := session.Start(fmt.Sprintf("cat %s", d.config.DumpName)); err != nil {
		return err
	}

	var downloaded int64
	buf := make([]byte, 32*1024)
	for {
		n, readErr := stdout.Read(buf)
		if n > 0 {
			if _, err := outFile.Write(buf[:n]); err != nil {
				return err
			}
			downloaded += int64(n)
			utils.Progress(downloaded, totalSize)
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	fmt.Println("\nDownload complete:", localPath)

	return session.Wait()
}

func (d *Download) FileSize() (int64, error) {
	sizeOutput, err := d.conn.RunCommand(fmt.Sprintf("stat -c %%s %s", d.config.DumpName))

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
