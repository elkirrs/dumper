package command_config

import (
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/option"
)

type ConfigData struct {
	User       string
	Password   string
	Name       string
	Port       string
	Key        string
	Host       string
	DumpName   string
	DumpFormat string
	Driver     string
	Options    option.Options
	Encrypt    encrypt.Encrypt
}
