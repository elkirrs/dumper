package config

import (
	"dumper/internal/domain/config/database"
	"dumper/internal/domain/config/server"
	"dumper/internal/domain/config/setting"
	"dumper/internal/domain/config/storage"
)

type Config struct {
	Settings  setting.Settings             `yaml:"settings" validate:"required" json:"settings"`
	Databases map[string]database.Database `yaml:"databases" validate:"required" json:"databases,omitempty"`
	Servers   map[string]server.Server     `yaml:"servers" validate:"required" json:"servers,omitempty"`
	Storages  map[string]storage.Storage   `yaml:"storages" validate:"required" json:"storages,omitempty"`
}
