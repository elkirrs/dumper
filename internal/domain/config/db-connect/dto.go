package db_connect

import (
	"dumper/internal/domain/config/database"
	"dumper/internal/domain/config/server"
)

type DBConnect struct {
	Database database.Database
	Server   server.Server
}
