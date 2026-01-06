package google

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils/console"
	"dumper/pkg/utils/stream"
	"fmt"
	"io"

	googleClient "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCS struct {
	ctx     context.Context
	config  *storage.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *GCS {
	return &GCS{
		ctx:     ctx,
		config:  config,
		backend: "GoogleCloud",
	}
}

func (gc *GCS) Save() error {

	var opts []option.ClientOption

	if gc.config.Config.Credential != "" {
		opts = append(
			opts,
			option.WithCredentialsJSON([]byte(gc.config.Config.Credential)),
		)
	} else if gc.config.Config.CredentialFile != "" {
		opts = append(
			opts,
			option.WithCredentialsFile(gc.config.Config.CredentialFile),
		)
	}

	client, err := googleClient.NewClient(gc.ctx, opts...)

	if err != nil {
		return &storage.UploadError{
			Backend: gc.backend,
			Err:     fmt.Errorf("failed to create GoogleCloud client: %w", err),
		}
	}

	defer client.Close()

	pr, closeSSH, err := stream.SSHStreamer(gc.ctx, gc.config.Conn, gc.config.DumpName, gc.config.FileSize)

	if err != nil {
		return &storage.UploadError{
			Backend: gc.backend,
			Err:     fmt.Errorf("failed to create SSH session: %v", err),
		}
	}

	defer closeSSH()
	targetPath := stream.TargetPath(gc.config.Config.Dir, gc.config.DumpName)

	writer := client.Bucket(gc.config.Config.Bucket).Object(targetPath).NewWriter(gc.ctx)
	writer.ContentType = "application/octet-stream"
	writer.ChunkSize = 32 * 1024 * 1024

	if _, err := io.Copy(writer, pr); err != nil {
		_ = writer.Close()
		return &storage.UploadError{Backend: gc.backend, Err: err}
	}

	if err := writer.Close(); err != nil {
		return &storage.UploadError{Backend: gc.backend, Err: err}
	}

	console.SafePrintln("[GCS] Upload complete: %s", targetPath)
	return nil
}
