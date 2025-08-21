package backup

import (
	"context"
	"dumper/internal/connect"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
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
	removeDump   bool
}

func NewApp(
	ctx context.Context,
	conn *connect.Connect,
	backupCmd,
	remotePath,
	localDir,
	dumpLocation string,
	removeDump bool,
) *Backup {
	return &Backup{
		ctx:          ctx,
		conn:         conn,
		backupCmd:    backupCmd,
		remotePath:   remotePath,
		localDir:     localDir,
		dumpLocation: dumpLocation,
		removeDump:   removeDump,
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

	isRemoveDump := b.removeDump
	checkCmd := fmt.Sprintf("test -f %s", b.remotePath)

	logging.L(b.ctx).Info(
		"Run command found backup in server with name",
		logging.StringAttr("name", b.remotePath),
	)

	if msg, err := b.conn.RunCommand(checkCmd); err == nil {
		logging.L(b.ctx).Info(
			"Dump already exists on server",
			logging.StringAttr("name", b.remotePath),
			logging.StringAttr("msg", msg),
		)

		fmt.Println("Dump already exists on server:", b.remotePath)
		isRemoveDump = false
	} else {
		stop := make(chan struct{})
		dumpCreateTimeNow := time.Now()

		logging.L(b.ctx).Info("File dump name", logging.StringAttr("name", b.remotePath))
		fmt.Println("File dump name:", b.remotePath)

		go utils.Spinner(stop)

		if msg, err := b.conn.RunCommand(b.backupCmd); err != nil {
			logging.L(b.ctx).Error(
				"Failed to create dump",
				logging.StringAttr("msg", msg),
				logging.ErrAttr(err),
			)
			return fmt.Errorf("failed to create dump: %v", err)
		}

		close(stop)

		elapsed := time.Since(dumpCreateTimeNow)

		totalSize, err := b.fileSize()
		if err != nil {
			return err
		}

		fmt.Printf("\rDump created successfully in %.2f sec\n", elapsed.Seconds())
		fmt.Printf("\rFile dump size: %s [%d bytes]\n", utils.FormatBytes(totalSize), totalSize)

		dumpCreateTimeSec := fmt.Sprintf("%.2f sec", elapsed.Seconds())
		logging.L(b.ctx).Info(
			"The dump was successfully created",
			logging.StringAttr("time", dumpCreateTimeSec),
			logging.Int64Attr("size", totalSize),
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
		if msg, err := b.conn.RunCommand(fmt.Sprintf("rm -f %s", b.remotePath)); err != nil {
			logging.L(b.ctx).Error(
				"Failed to remove dump on server",
				logging.StringAttr("msg", msg),
			)
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

	var totalSize int64

	totalSize, err := b.fileSize()
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

func (b *Backup) fileSize() (int64, error) {
	sizeOutput, err := b.conn.RunCommand(fmt.Sprintf("stat -c %%s %s", b.remotePath))

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
