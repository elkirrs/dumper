package database

import (
	"dumper/internal/domain/config/docker"
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/shell"
)

type Database struct {
	Title      string           `yaml:"title,omitempty"`
	User       string           `yaml:"user"`
	Password   string           `yaml:"password"`
	Name       string           `yaml:"name"`
	Server     string           `yaml:"server" validate:"required"`
	Key        string           `yaml:"key"`
	Port       string           `yaml:"port"`
	Driver     string           `yaml:"driver" validate:"required"`
	Format     string           `yaml:"format" validate:"required"`
	Options    *option.Options  `yaml:"options"`
	RemoveDump *bool            `yaml:"remove_dump"`
	Encrypt    *encrypt.Encrypt `yaml:"encrypt"`
	Storages   []string         `yaml:"storages" validate:"required"`
	Archive    *bool            `yaml:"archive"`
	Docker     *docker.Docker   `yaml:"docker"`
	Shell      *shell.Shell     `yaml:"shell"`
	DirRemote  string           `yaml:"dir_remote"`
}

func (d *Database) GetName() string {
	if d.Name != "" {
		return d.Name
	}
	return d.User
}

func (d *Database) GetPort(port *string) string {
	if d.Port != "" {
		return d.Port
	}
	return *port
}

func (d *Database) GetDriver(driver *string) string {
	if d.Driver != "" {
		return d.Driver
	}
	return *driver
}

func (d *Database) GetRemoveDump(removeDump bool) bool {
	if d.RemoveDump != nil {
		return *d.RemoveDump
	}
	return removeDump
}

func (d *Database) GetFormat(format *string) string {
	if d.Format != "" {
		return d.Format
	}
	return *format
}

func (d *Database) GetEncrypt(encryptGlobal *encrypt.Encrypt) encrypt.Encrypt {

	if d.Encrypt == nil && encryptGlobal == nil {
		val := false
		return encrypt.Encrypt{
			Enabled: &val,
		}
	}

	if d.Encrypt == nil && encryptGlobal != nil {
		return *encryptGlobal
	}

	if *d.Encrypt.Enabled && d.Encrypt.Password == "" && d.Encrypt.Type == "" {
		return *encryptGlobal
	}

	return *d.Encrypt
}

func (d *Database) GetTitle() string {
	if d.Title != "" {
		return d.Title
	}

	if d.Name != "" {
		return d.Name
	}

	return d.User
}

func (d *Database) IsArchive(isGlobalArchive bool) bool {
	if d.Archive != nil {
		return *d.Archive
	}
	return isGlobalArchive
}

func (d *Database) GetDocker(globalDocker *docker.Docker) docker.Docker {
	if d.Docker == nil && globalDocker == nil {
		val := false
		return docker.Docker{
			Enabled: &val,
		}
	}

	if d.Docker == nil && globalDocker != nil {
		return *globalDocker
	}

	if *d.Docker.Enabled && d.Docker.Command == "" {
		return *globalDocker
	}

	return *globalDocker
}

func (d *Database) GetOptions() option.Options {
	if d.Options != nil {
		return *d.Options
	}
	return option.Options{}
}

func (d *Database) GetShell(serverShell *shell.Shell) shell.Shell {
	if d.Shell == nil && serverShell == nil {
		val := false
		return shell.Shell{Enabled: &val}
	}

	if d.Shell == nil && serverShell != nil {
		return *serverShell
	}

	if *d.Shell.Enabled && d.Shell.Before == "" && d.Shell.After == "" {
		return *serverShell
	}

	return *d.Shell
}

func (d *Database) GetStorages(globalStorages *[]string) []string {
	if d.Storages == nil || len(d.Storages) == 0 {
		return *globalStorages
	}
	return d.Storages
}

func (d *Database) GetDirRemote(globalDirRemote *string) string {
	if d.DirRemote != "" {
		return d.DirRemote
	}
	return *globalDirRemote
}
