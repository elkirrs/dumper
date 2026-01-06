package cloudflare

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type Cloudflare struct {
	ctx     context.Context
	config  *storage.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *Cloudflare {
	return &Cloudflare{
		ctx:     ctx,
		config:  config,
		backend: "Cloudflare R2",
	}
}

func (c *Cloudflare) Save() error {
	c.config.Config.Region = "auto"
	awsClient := aws.NewClient(
		c.ctx,
		c.config.Conn,
		c.config.Config,
		c.config.DumpName,
		c.config.FileSize,
		c.backend,
	)

	return awsClient.Handler()
}
