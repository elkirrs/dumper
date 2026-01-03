package minio

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type Minio struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *Minio {
	return &Minio{
		ctx:    ctx,
		config: config,
	}
}

func (m *Minio) Save() error {
	awsClient := aws.NewClient(
		m.ctx,
		m.config.Conn,
		m.config.Config,
		m.config.DumpName,
		m.config.FileSize,
	)

	return awsClient.Handler()
}
