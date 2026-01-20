package database

import (
	"dumper/internal/domain/config/docker"
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/shell"
)

type Database struct {
	Title      string           `yaml:"title,omitempty" json:"title,omitempty"`
	User       string           `yaml:"user" json:"user,omitempty"`
	Password   string           `yaml:"password" json:"password,omitempty"`
	Name       string           `yaml:"name" json:"name,omitempty"`
	Server     string           `yaml:"server" validate:"required" json:"server,omitempty"`
	Key        string           `yaml:"key" json:"key,omitempty"`
	Port       string           `yaml:"port" json:"port,omitempty"`
	Driver     string           `yaml:"driver" validate:"required" json:"driver,omitempty"`
	Format     string           `yaml:"format" validate:"required" json:"format,omitempty"`
	Options    *option.Options  `yaml:"options" json:"options,omitempty"`
	RemoveDump *bool            `yaml:"remove_dump" json:"removeDump,omitempty"`
	Encrypt    *encrypt.Encrypt `yaml:"encrypt" json:"encrypt,omitempty"`
	Storages   []string         `yaml:"storages" validate:"required" json:"storages,omitempty"`
	Archive    *bool            `yaml:"archive" json:"archive,omitempty"`
	Docker     *docker.Docker   `yaml:"docker" json:"docker,omitempty"`
	Shell      *shell.Shell     `yaml:"shell" json:"shell,omitempty"`
	DirRemote  string           `yaml:"dir_remote" json:"dirRemote,omitempty"`
	Token      string           `yaml:"token" json:"token,omitempty"`
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

func (d *Database) GetRemoveDump(removeDumpGlobal *bool) bool {
	if d.RemoveDump != nil {
		return *d.RemoveDump
	}
	return *removeDumpGlobal
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
	if len(d.Storages) == 0 {
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
