package backblaze

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type Backblaze struct {
	ctx     context.Context
	config  *storage.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *Backblaze {
	return &Backblaze{
		ctx:     ctx,
		config:  config,
		backend: "Backblaze B2",
	}
}

func (b *Backblaze) Save() error {
	awsClient := aws.NewClient(
		b.ctx,
		b.config.Conn,
		b.config.Config,
		b.config.DumpName,
		b.config.FileSize,
		b.backend,
	)

	return awsClient.Handler()
}
