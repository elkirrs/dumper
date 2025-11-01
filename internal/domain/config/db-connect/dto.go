package db_connect

import (
	"dumper/internal/domain/config/database"
	"dumper/internal/domain/config/server"
	"dumper/internal/domain/config/storage"
)

type DBConnect struct {
	Database database.Database
	Server   server.Server
	Storages map[string]storage.ListStorages
}
