package command_config

import (
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/option"
)

type Database struct {
	User     string
	Password string
	Name     string
	Port     string
	Format   string
	Driver   string
	Options  option.Options
}

type Server struct {
	Host string
	Port string
	Key  string
}

type Config struct {
	Database     Database
	Server       Server
	DumpLocation string
	Archive      bool
	RemoveBackup bool
	Command      string
	DumpName     string
	DumpDirLocal string
	Encrypt      encrypt.Encrypt
}
