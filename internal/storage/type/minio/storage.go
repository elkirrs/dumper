package minio

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type Minio struct {
	ctx     context.Context
	config  *storage.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *Minio {
	return &Minio{
		ctx:     ctx,
		config:  config,
		backend: "MinIO",
	}
}

func (m *Minio) Save() error {
	awsClient := aws.NewClient(
		m.ctx,
		m.config.Conn,
		m.config.Config,
		m.config.DumpName,
		m.config.FileSize,
		m.backend,
	)

	return awsClient.Handler()
}
