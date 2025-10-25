package database

import "dumper/internal/domain/config/option"

type Database struct {
	User       string         `yaml:"user"`
	Password   string         `yaml:"password"`
	Name       string         `yaml:"name,omitempty"`
	Server     string         `yaml:"server" validate:"required"`
	Key        string         `yaml:"key"`
	Port       string         `yaml:"port,omitempty"`
	Driver     string         `yaml:"driver"`
	Options    option.Options `yaml:"options"`
	RemoveDump *bool          `yaml:"remove_dump"`
}

func (d Database) GetName() string {
	if d.Name != "" {
		return d.Name
	}
	return d.User
}

func (d Database) GetPort(port string) string {
	if d.Port != "" {
		return d.Port
	}
	return port
}

func (d Database) GetDriver(driver string) string {
	if d.Driver != "" {
		return d.Driver
	}
	return driver
}
func (d Database) GetAuthSource() string {
	if d.Options.AuthSource != "" {
		return d.Options.AuthSource
	}
	return d.GetName()
}

func (d Database) GetRemoveDump(removeDump bool) bool {
	if d.RemoveDump != nil {
		return *d.RemoveDump
	}
	return removeDump
}

func (d Database) GetDisplayName(key string) string {
	if d.Name != "" {
		return d.Name
	}
	return key
}
