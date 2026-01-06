package yandex

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type Yandex struct {
	ctx     context.Context
	config  *storage.Config
	backend string
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *Yandex {
	return &Yandex{
		ctx:     ctx,
		config:  config,
		backend: "Yandex Cloud",
	}
}

func (y *Yandex) Save() error {

	awsClient := aws.NewClient(
		y.ctx,
		y.config.Conn,
		y.config.Config,
		y.config.DumpName,
		y.config.FileSize,
		y.backend,
	)
	return awsClient.Handler()
}
