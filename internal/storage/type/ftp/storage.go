package ftp

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils/console"
	"dumper/pkg/utils/stream"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTP struct {
	ctx     context.Context
	config  *storage.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *FTP {
	return &FTP{
		ctx:     ctx,
		config:  config,
		backend: "FTP",
	}
}

func (f *FTP) Save() error {
	addr := fmt.Sprintf("%s:%s", f.config.Config.Host, f.config.Config.Port)
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(10*time.Second))
	targetPath := stream.TargetPath(f.config.Config.Dir, f.config.DumpName)

	if err != nil {
		return &storage.UploadError{
			Backend: f.backend,
			Err:     fmt.Errorf("failed to connect to FTP server: %w", err),
		}
	}

	defer c.Quit()

	if err := c.Login(f.config.Config.Username, f.config.Config.Password); err != nil {
		return &storage.UploadError{
			Backend: f.backend,
			Err:     fmt.Errorf("login failed: %w", err),
		}
	}

	dir := filepath.Dir(targetPath)
	if err := f.checkDirAccessible(c, dir); err != nil {
		return &storage.UploadError{
			Backend: f.backend,
			Err:     err,
		}
	}

	pr, closeSSH, err := stream.SSHStreamer(f.ctx, f.config.Conn, f.config.DumpName, f.config.FileSize)

	if err != nil {
		return &storage.UploadError{
			Backend: f.backend,
			Err:     fmt.Errorf("failed to create SSH session: %v", err),
		}
	}

	defer closeSSH()

	if err := c.Stor(targetPath, pr); err != nil {
		return &storage.UploadError{
			Backend: f.backend,
			Err:     fmt.Errorf("failed to upload file via FTP: %v", err),
		}
	}

	console.SafePrintln("[FTP] Upload complete: %s", targetPath)
	return nil
}

func isDirExistsError(err error) bool {
	return err != nil && (err.Error() == "550 Create directory operation failed." || err.Error() == "550")
}

func makeDirRecursive(c *ftp.ServerConn, path string) error {
	dirs := strings.Split(path, "/")
	curr := ""
	for _, d := range dirs {
		if d == "" {
			continue
		}
		if curr == "" {
			curr = d
		} else {
			curr = curr + "/" + d
		}
		err := c.MakeDir(curr)
		if err != nil && !isDirExistsError(err) {
			return err
		}
	}
	return nil
}

func (f *FTP) checkDirAccessible(c *ftp.ServerConn, dir string) error {
	if err := makeDirRecursive(c, dir); err != nil {
		return fmt.Errorf("FTP directory %s is not accessible: %w", dir, err)
	}

	if err := c.ChangeDir(dir); err != nil {
		return fmt.Errorf("no permission to access FTP directory %s: %w", dir, err)
	}

	_ = c.ChangeDir("/")

	return nil
}
