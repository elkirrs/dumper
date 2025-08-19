package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Settings  Settings            `yaml:"settings" validate:"required" json:"settings"`
	Databases map[string]Database `yaml:"databases" validate:"required" json:"databases,omitempty"`
	Servers   map[string]Server   `yaml:"servers" validate:"required" json:"servers,omitempty"`
	Licence   string              `json:"licence,omitempty"`
}

type Settings struct {
	SSH          SSHConfig `yaml:"ssh"`
	Template     string    `yaml:"template" default:"{%srv%}_{%db%}_{%time%}"`
	Archive      *bool     `yaml:"archive" default:"true"`
	Driver       string    `yaml:"driver" validate:"required"`
	DBPort       string    `yaml:"db_port,omitempty"`
	SrvKey       string    `yaml:"server_key,omitempty"`
	SrvPost      string    `yaml:"server_port,omitempty"`
	DumpLocation string    `yaml:"location" default:"server"` // server, local-ssh, local-direct
	DumpFormat   string    `yaml:"format" default:"plain"`    // plain, dump, tar
	DirDump      string    `yaml:"dir_dump" default:"./"`
	DirArchived  string    `yaml:"dir_archived" default:"./archived"`
	Logging      *bool     `yaml:"logging" default:"false"`
	RetryConnect int       `yaml:"retry_connect" default:"3"`
}

type Database struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name,omitempty"`
	Server   string `yaml:"server" validate:"required"`
	Key      string `yaml:"key"`
	Port     string `yaml:"port,omitempty"`
}

type Server struct {
	Host     string `yaml:"host" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Name     string `yaml:"name,omitempty"`
	Port     string `yaml:"port,omitempty"`
	SSHKey   string `yaml:"key,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type SSHConfig struct {
	PrivateKey   string `yaml:"private_key"`
	Passphrase   string `yaml:"passphrase"`
	IsPassphrase *bool  `yaml:"is_passphrase" validate:"required"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := defaults.Set(&config); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	for k, server := range config.Servers {
		if server.Name == "" {
			server.Name = server.Host
			config.Servers[k] = server
		}
	}

	return &config, nil
}

func (s Server) GetDisplayName() string {
	if s.Name != "" {
		return s.Name
	}
	return s.Host
}

func (s Server) GetPort(port string) string {
	if s.Port != "" {
		return s.Port
	}
	return port
}

func (d Database) GetDisplayName() string {
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
