package azure

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"golang.org/x/crypto/ssh"
)

type Azure struct {
	ctx    context.Context
	config *storage.Config
	client *azblob.Client
}

func NewApp(ctx context.Context, config *storage.Config) *Azure {
	return &Azure{
		ctx:    ctx,
		config: config,
	}
}

func (a *Azure) authType() error {

	switch a.config.Config.AuthType {
	case "SharedKey":
		return a.clientSharedKey()
	case "AAD":
		return a.clientADD()
	default:
		return fmt.Errorf("unsupported auth type: %s", a.config.Config.AuthType)
	}
}

func (a *Azure) Save() error {
	blobName := filepath.Base(a.config.DumpName)
	containerName := a.config.Config.Container
	targetPath := blobName

	err := a.authType()
	if err != nil {
		return err
	}

	containerClient := a.client.ServiceClient().NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(targetPath)

	session, err := a.config.Conn.NewSession()
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

	if err := session.Start(fmt.Sprintf("cat %s", a.config.DumpName)); err != nil {
		return fmt.Errorf("failed to start remote command: %v", err)
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
			case <-a.ctx.Done():
				_ = pw.CloseWithError(fmt.Errorf("azure upload cancelled by context"))
				return
			default:
			}

			n, readErr := stdout.Read(buf)
			if n > 0 {
				uploaded += int64(n)
				if gp, ok := a.ctx.Value("globalProgress").(*utils.GlobProgress); ok {
					gp.Add(int64(n))
				} else {
					utils.Progress(uploaded, a.config.FileSize)
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

	_, err = blobClient.UploadStream(a.ctx, pr, &azblob.UploadStreamOptions{
		BlockSize: 32 * 1024,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to azure: %v", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %v", err)
	}

	fmt.Println("\n[Azure] Upload complete:", targetPath)
	return nil
}

func (a *Azure) clientSharedKey() error {
	cred, err := azblob.NewSharedKeyCredential(
		a.config.Config.Name,
		a.config.Config.SharedKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create shared key credential: %v", err)
	}

	a.client, err = azblob.NewClientWithSharedKeyCredential(
		a.config.Config.Endpoint,
		cred,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create azure blob client: %v", err)
	}

	return nil
}

func (a *Azure) clientADD() error {
	cred, err := azidentity.NewClientSecretCredential(
		a.config.Config.TenantID,
		a.config.Config.ClientID,
		a.config.Config.ClientSecret,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create azure credential: %v", err)
	}

	a.client, err = azblob.NewClient(a.config.Config.Endpoint, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create azure client: %v", err)
	}

	return nil
}
