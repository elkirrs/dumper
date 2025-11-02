package ftp

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"golang.org/x/crypto/ssh"
)

type FTP struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(ctx context.Context, config *storage.Config) *FTP {
	return &FTP{
		ctx:    ctx,
		config: config,
	}
}

func (f *FTP) Save() error {
	addr := fmt.Sprintf("%s:%s", f.config.Config.Host, f.config.Config.Port)
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to FTP server: %v", err)
	}
	defer func(c *ftp.ServerConn) { _ = c.Quit() }(c)

	if err := c.Login(f.config.Config.Username, f.config.Config.Password); err != nil {
		return fmt.Errorf("FTP login failed: %v", err)
	}

	targetPath := filepath.Join(f.config.Config.Dir, filepath.Base(f.config.DumpName))
	dir := filepath.Dir(targetPath)
	if err := c.MakeDir(dir); err != nil && !isDirExistsError(err) {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	session, err := f.config.Conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	if err := session.Start(fmt.Sprintf("cat %s", f.config.DumpName)); err != nil {
		return err
	}

	pr, pw := io.Pipe()
	go func() {
		defer func(pw *io.PipeWriter) {
			_ = pw.Close()
		}(pw)
		buf := make([]byte, 32*1024)
		var uploaded int64

		for {
			select {
			case <-f.ctx.Done():
				_ = pw.CloseWithError(fmt.Errorf("FTP upload cancelled by context"))
				return
			default:
			}

			n, readErr := stdout.Read(buf)
			if n > 0 {
				uploaded += int64(n)
				if gp, ok := f.ctx.Value("globalProgress").(*utils.GlobProgress); ok {
					gp.Add(int64(n))
				} else {
					utils.Progress(uploaded, f.config.FileSize)
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

	if err := c.Stor(targetPath, pr); err != nil {
		return fmt.Errorf("failed to upload file via FTP: %v", err)
	}

	utils.SafePrintln("[FTP] Upload complete: %s", targetPath)
	return session.Wait()
}

func isDirExistsError(err error) bool {
	return err != nil && (err.Error() == "550 Create directory operation failed." || err.Error() == "550")
}
