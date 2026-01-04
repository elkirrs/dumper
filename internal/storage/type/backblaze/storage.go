package backblaze

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type Backblaze struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *Backblaze {
	return &Backblaze{
		ctx:    ctx,
		config: config,
	}
}

func (b *Backblaze) Save() error {
	awsClient := aws.NewClient(
		b.ctx,
		b.config.Conn,
		b.config.Config,
		b.config.DumpName,
		b.config.FileSize,
	)

	return awsClient.Handler()
}
