package digitalOcean

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type DigitalOcean struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *DigitalOcean {
	return &DigitalOcean{
		ctx:    ctx,
		config: config,
	}
}

func (d *DigitalOcean) Save() error {
	awsClient := aws.NewClient(
		d.ctx,
		d.config.Conn,
		d.config.Config,
		d.config.DumpName,
		d.config.FileSize,
	)

	return awsClient.Handler()
}
