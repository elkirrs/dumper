package storage

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/internal/storage/local"
	"errors"
	"fmt"
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
		//handler = sftp.NewApp(a.ctx, a.config.SFTP)
	case "ftp":
		//handler = ftp.NewApp(a.ctx, a.config.FTP)
	default:
		return errors.New("unsupported storage type: " + s.config.Type)
	}

	if handler == nil {
		return fmt.Errorf("storage handler not initialized for type %s", s.config.Type)
	}

	return handler.Save()
}
