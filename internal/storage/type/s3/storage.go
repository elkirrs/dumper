package s3

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/client/aws"
)

type S3 struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *S3 {
	return &S3{
		ctx:    ctx,
		config: config,
	}
}

func (s *S3) Save() error {
	awsClient := aws.NewClient(
		s.ctx,
		s.config.Conn,
		s.config.Config,
		s.config.DumpName,
		s.config.FileSize,
	)

	return awsClient.Handler()
}
