package setting

import (
	"dumper/internal/domain/config/docker"
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/shell"
	sshConfig "dumper/internal/domain/config/ssh-config"
)

type Settings struct {
	SSH                 *sshConfig.SSHConfig `yaml:"ssh"`
	DirRemote           string               `yaml:"dir_remote" default:"./"`
	Template            string               `yaml:"template" default:"{%srv%}_{%db%}_{%time%}"`
	Archive             *bool                `yaml:"archive" default:"true"`
	Driver              string               `yaml:"driver" validate:"required"`
	DBPort              string               `yaml:"db_port,omitempty"`
	SrvPost             string               `yaml:"server_port,omitempty"`
	DumpLocation        string               `yaml:"location" default:"server"`
	DumpFormat          string               `yaml:"format" default:"plain"`
	DirDump             string               `yaml:"dir_dump" default:"./"`
	DirArchived         string               `yaml:"dir_archived" default:"./archived"`
	Logging             *bool                `yaml:"logging" default:"true"`
	RetryConnect        int                  `yaml:"retry_connect" default:"3"`
	RemoveDump          *bool                `yaml:"remove_dump" default:"true"`
	Encrypt             *encrypt.Encrypt     `yaml:"encrypt"`
	Storages            []string             `yaml:"storages"`
	MaxParallelDownload int                  `yaml:"parallel_download" default:"2"`
	Docker              *docker.Docker       `yaml:"docker"`
	Shell               *shell.Shell         `yaml:"shell"`
}
