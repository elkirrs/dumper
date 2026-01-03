package storage

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/azure"
	"dumper/internal/storage/ftp"
	"dumper/internal/storage/local"
	"dumper/internal/storage/minio"
	"dumper/internal/storage/s3"
	"dumper/internal/storage/sftp"
	"errors"
)

type StorageHandler interface {
	Save() error
}

type Storage struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	cfg *storage.Config,
) *Storage {
	return &Storage{
		ctx:    ctx,
		config: cfg,
	}
}

func (s *Storage) Save() error {
	var handler StorageHandler

	switch s.config.Type {
	case "local":
		handler = local.NewApp(s.ctx, s.config)
	case "sftp":
		handler = sftp.NewApp(s.ctx, s.config)
	case "ftp":
		handler = ftp.NewApp(s.ctx, s.config)
	case "azure":
		handler = azure.NewApp(s.ctx, s.config)
	case "s3":
		handler = s3.NewApp(s.ctx, s.config)
	case "minio":
		handler = minio.NewApp(s.ctx, s.config)
	default:
		return errors.New("unsupported storage type: " + s.config.Type)
	}

	return handler.Save()
}
