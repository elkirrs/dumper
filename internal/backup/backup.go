package backup

import (
	"context"
	"dumper/internal/connect"
	"dumper/pkg/logging"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Backup struct {
	ctx          context.Context
	conn         *connect.Connect
	backupCmd    string
	remotePath   string
	localDir     string
	dumpLocation string
}

func NewApp(
	ctx context.Context,
	conn *connect.Connect,
	backupCmd,
	remotePath,
	localDir,
	dumpLocation string,
) *Backup {
	return &Backup{
		ctx:          ctx,
		conn:         conn,
		backupCmd:    backupCmd,
		remotePath:   remotePath,
		localDir:     localDir,
		dumpLocation: dumpLocation,
	}
}

func (b *Backup) Backup() error {
	switch b.dumpLocation {
	case "server":
		return b.backupByServer()
	case "local-ssh":
		return b.backupByLocalSSH()
	case "local-direct":
		return b.backupLocalDirect()
	default:
		logging.L(b.ctx).Error(
			"Unsupported backup dump location",
			logging.StringAttr("location", b.dumpLocation),
		)
		return fmt.Errorf("unsupported backup dump location: %s", b.dumpLocation)
	}
}

func (b *Backup) backupByServer() error {

	isRemoveDump := true
	checkCmd := fmt.Sprintf("test -f %s", b.remotePath)

	logging.L(b.ctx).Info(
		"Run command found backup in server with name",
		logging.StringAttr("name", b.remotePath),
	)

	if _, err := b.conn.RunCommand(checkCmd); err == nil {
		logging.L(b.ctx).Info("Dump already exists on server", logging.StringAttr("name", b.remotePath))

		fmt.Println("Dump already exists on server:", b.remotePath)
		isRemoveDump = false
	} else {
		dumpCreateTimeNow := time.Now()

		logging.L(b.ctx).Info("Creating dump", logging.StringAttr("name", b.remotePath))
		fmt.Println("Creating dump:", b.remotePath)
		if _, err := b.conn.RunCommand(b.backupCmd); err != nil {
			logging.L(b.ctx).Error("Failed to create dump")
			return fmt.Errorf("failed to create dump: %v", err)
		}

		dumpCreateTimeSec := fmt.Sprintf("%.2f sec", time.Since(dumpCreateTimeNow).Seconds())
		logging.L(b.ctx).Info(
			"The dump was successfully created",
			logging.StringAttr("time", dumpCreateTimeSec),
		)
	}

	logging.L(b.ctx).Info("Downloading dump", logging.StringAttr("name", b.remotePath))
	dumpDownloadTimeNow := time.Now()
	if err := b.downloadFile(); err != nil {
		logging.L(b.ctx).Error("Failed to download dump")
		return fmt.Errorf("failed to download dump: %v", err)
	}

	dumpDownloadTimeSec := fmt.Sprintf("%.2f sec", time.Since(dumpDownloadTimeNow).Seconds())

	logging.L(b.ctx).Info("The dump was successfully downloaded", logging.StringAttr("time", dumpDownloadTimeSec))

	if isRemoveDump {
		logging.L(b.ctx).Info("Removing dump on server")
		fmt.Println("Removing dump from server:", b.remotePath)
		if _, err := b.conn.RunCommand(fmt.Sprintf("rm -f %s", b.remotePath)); err != nil {
			logging.L(b.ctx).Error("Failed to remove dump on server")
			return fmt.Errorf("failed to delete dump on server: %v", err)
		}

		logging.L(b.ctx).Info("The dump was successfully deleted on server")
	}

	return nil
}

func (b *Backup) backupByLocalSSH() error {
	panic("not implement")
}

func (b *Backup) backupLocalDirect() error {
	panic("not implement")
}

func (b *Backup) downloadFile() error {
	localPath := filepath.Join(b.localDir, filepath.Base(b.remotePath))

	sizeOutput, err := b.conn.RunCommand(fmt.Sprintf("stat -c %%s %s", b.remotePath))
	if err != nil {
		return fmt.Errorf("failed to get file size: %v", err)
	}
	sizeOutput = strings.TrimSpace(sizeOutput)

	var totalSize int64

	_, err = fmt.Sscanf(sizeOutput, "%d", &totalSize)
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

	session, err := b.conn.NewSession()
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

	if err := session.Start(fmt.Sprintf("cat %s", b.remotePath)); err != nil {
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
			printProgress(downloaded, totalSize)
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

func printProgress(done, total int64) {
	if total == 0 {
		fmt.Printf("\rDownloaded: %d bytes", done)
		return
	}
	percent := float64(done) / float64(total) * 100
	fmt.Printf("\rDownloading... %.1f%% (%d/%d bytes)", percent, done, total)
}
