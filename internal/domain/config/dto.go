package config

import (
	"dumper/internal/domain/config/database"
	"dumper/internal/domain/config/server"
	"dumper/internal/domain/config/setting"
)

type Config struct {
	Settings  setting.Settings             `yaml:"settings" validate:"required" json:"settings"`
	Databases map[string]database.Database `yaml:"databases" validate:"required" json:"databases,omitempty"`
	Servers   map[string]server.Server     `yaml:"servers" validate:"required" json:"servers,omitempty"`
	Licence   string                       `json:"licence,omitempty"`
}
