package storage

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/type/azure"
	"dumper/internal/storage/type/backblaze"
	"dumper/internal/storage/type/cloudflare"
	"dumper/internal/storage/type/ftp"
	"dumper/internal/storage/type/local"
	"dumper/internal/storage/type/minio"
	"dumper/internal/storage/type/s3"
	"dumper/internal/storage/type/sftp"
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
	case "r2":
		handler = cloudflare.NewApp(s.ctx, s.config)
	case "b2":
		handler = backblaze.NewApp(s.ctx, s.config)
	default:
		return errors.New("unsupported storage type: " + s.config.Type)
	}

	return handler.Save()
}
