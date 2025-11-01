package local

import (
	"context"
	storageDomain "dumper/internal/domain/storage"
	"dumper/pkg/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type Local struct {
	ctx    context.Context
	config *storageDomain.Config
}

func NewApp(
	ctx context.Context,
	config *storageDomain.Config,
) *Local {
	return &Local{
		ctx:    ctx,
		config: config,
	}
}

func (l *Local) Save() error {
	localPath := filepath.Join(l.config.Config.Dir, filepath.Base(l.config.DumpName))

	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}

	defer func(outFile *os.File) {
		_ = outFile.Close()
		return
	}(outFile)

	session, err := l.config.Conn.NewSession()
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

	if err := session.Start(fmt.Sprintf("cat %s", l.config.DumpName)); err != nil {
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
			utils.Progress(downloaded, l.config.FileSize)
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
