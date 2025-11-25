package command_config

import (
	"dumper/internal/domain/config/docker"
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/shell"
	"dumper/internal/domain/config/storage"
)

type Database struct {
	User     string
	Password string
	Name     string
	Port     string
	Format   string
	Driver   string
	Options  option.Options
	Archive  bool
	Docker   docker.Docker
}

type Server struct {
	Host string
	Port string
	Key  string
}

type Config struct {
	Database            Database
	Server              Server
	Storages            map[string]storage.Storage
	DumpLocation        string
	Archive             bool
	RemoveBackup        bool
	Command             string
	DumpName            string
	DumpDirRemote       string
	DumpDirLocal        string
	Encrypt             encrypt.Encrypt
	MaxParallelDownload int
	Shell               shell.Shell
}
