package database

import (
	"dumper/internal/domain/config/docker"
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/shell"
)

type Database struct {
	Title      string          `yaml:"title,omitempty"`
	User       string          `yaml:"user"`
	Password   string          `yaml:"password"`
	Name       string          `yaml:"name,omitempty"`
	Server     string          `yaml:"server" validate:"required"`
	Key        string          `yaml:"key"`
	Port       string          `yaml:"port,omitempty"`
	Driver     string          `yaml:"driver"`
	Format     string          `yaml:"format" validate:"required"`
	Options    option.Options  `yaml:"options"`
	RemoveDump *bool           `yaml:"remove_dump"`
	Encrypt    encrypt.Encrypt `yaml:"encrypt"`
	Storages   []string        `yaml:"storages"`
	Archive    *bool           `yaml:"archive"`
	Docker     *docker.Docker  `yaml:"docker"`
	Shell      *shell.Shell    `yaml:"shell"`
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

func (d Database) GetFormat(format string) string {
	if d.Format != "" {
		return d.Format
	}
	return format
}

func (d Database) GetEncryptType(crypt string) string {
	if d.Encrypt.Type != "" {
		return d.Encrypt.Type
	}
	return crypt
}

func (d Database) GetEncryptPass(pass string) string {
	if d.Encrypt.Password != "" {
		return d.Encrypt.Password
	}
	return pass
}

func (d Database) GetTitle() string {
	if d.Title != "" {
		return d.Title
	}

	if d.Name != "" {
		return d.Name
	}

	return d.User
}

func (d Database) IsArchive(isGlobalArchive bool) bool {
	if d.Archive != nil {
		return *d.Archive
	}
	return isGlobalArchive
}

func (d Database) GetDocker(globalDocker *docker.Docker) *docker.Docker {
	if d.Docker != nil {
		return d.Docker
	}
	return globalDocker
}
