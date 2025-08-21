package command_config

import "dumper/internal/config"

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
	Options    config.Options
}
