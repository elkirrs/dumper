package storage

import (
	"dumper/internal/connect"
	"dumper/internal/domain/config/storage"
	"fmt"
)

type Config struct {
	Type     string
	DumpName string
	FileSize int64
	Conn     *connect.Connect
	Config   storage.Storage
}

type Uploader interface {
	Save() error
}

type UploadError struct {
	Backend string
	Err     error
}

func (e *UploadError) Error() string {
	return fmt.Sprintf("[%s] upload failed: %v", e.Backend, e.Err)
}
