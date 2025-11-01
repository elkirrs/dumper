package storage

import (
	"dumper/internal/connect"
	"dumper/internal/domain/config/storage"
)

type Config struct {
	Type     string
	DumpName string
	FileSize int64
	Conn     *connect.Connect
	Config   storage.Storage
}
