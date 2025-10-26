package setting

import (
	"dumper/internal/domain/config/encrypt"
	sshConfig "dumper/internal/domain/config/ssh-config"
)

type Settings struct {
	SSH          sshConfig.SSHConfig `yaml:"ssh"`
	Template     string              `yaml:"template" default:"{%srv%}_{%db%}_{%time%}"`
	Archive      *bool               `yaml:"archive" default:"true"`
	Driver       string              `yaml:"driver" validate:"required"`
	DBPort       string              `yaml:"db_port,omitempty"`
	SrvKey       string              `yaml:"server_key,omitempty"`
	SrvPost      string              `yaml:"server_port,omitempty"`
	DumpLocation string              `yaml:"location" default:"server"` // server, local-ssh, local-direct
	DumpFormat   string              `yaml:"format" default:"plain"`    // plain, dump, tar
	DirDump      string              `yaml:"dir_dump" default:"./"`
	DirArchived  string              `yaml:"dir_archived" default:"./archived"`
	Logging      *bool               `yaml:"logging" default:"false"`
	RetryConnect int                 `yaml:"retry_connect" default:"3"`
	RemoveDump   *bool               `yaml:"remove_dump" default:"true"`
	Encrypt      encrypt.Encrypt     `yaml:"encrypt"`
}
