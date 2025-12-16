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

	dir := filepath.Dir(localPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create local directory: %v", err)
	}

	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	session, err := l.config.Conn.NewSession()
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

	if err := session.Start(fmt.Sprintf("cat %s", l.config.DumpName)); err != nil {
		return fmt.Errorf("failed to start remote command: %v", err)
	}

	pr, pw := io.Pipe()
	defer func(pr *io.PipeReader) {
		_ = pr.Close()
	}(pr)

	go func() {
		defer func(pw *io.PipeWriter) {
			_ = pw.Close()
		}(pw)
		buf := make([]byte, 32*1024)
		var downloaded int64

		for {
			select {
			case <-l.ctx.Done():
				_ = pw.CloseWithError(fmt.Errorf("download cancelled by context"))
				return
			default:
			}

			n, readErr := stdout.Read(buf)
			if n > 0 {
				downloaded += int64(n)
				if gp, ok := l.ctx.Value("globalProgress").(*utils.GlobProgress); ok {
					gp.Add(int64(n))
				} else {
					utils.Progress(downloaded, l.config.FileSize)
				}
				if _, err := pw.Write(buf[:n]); err != nil {
					return
				}
			}

			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				_ = pw.CloseWithError(readErr)
				return
			}
		}
	}()

	if _, err := io.Copy(outFile, pr); err != nil {
		return fmt.Errorf("failed to write to local file: %v", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %v", err)
	}

	utils.SafePrintln("[Local] Upload complete: %s", localPath)
	return nil
}
