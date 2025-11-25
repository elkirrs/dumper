package server

import "dumper/internal/domain/config/shell"

type Server struct {
	Title      string       `yaml:"title,omitempty"`
	Host       string       `yaml:"host" validate:"required"`
	User       string       `yaml:"user" validate:"required"`
	Name       string       `yaml:"name,omitempty"`
	Port       string       `yaml:"port,omitempty"`
	PrivateKey string       `yaml:"private_key,omitempty" validate:"xor=Password"`
	Passphrase string       `yaml:"passphrase,omitempty"`
	Password   string       `yaml:"password,omitempty" validate:"xor=PrivateKey"`
	ConfigPath string       `yaml:"conf_path,omitempty"`
	Shell      *shell.Shell `yaml:"shell,omitempty"`
}

func (s *Server) GetName() string {
	if s.Name != "" {
		return s.Name
	}
	return s.Host
}

func (s *Server) GetPort(port *string) string {
	if s.Port != "" {
		return s.Port
	}
	return *port
}

func (s *Server) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	if s.Name != "" {
		return s.Name
	}

	return s.Host
}

func (s *Server) GetPrivateKey(pathKey *string) string {
	if s.PrivateKey != "" {
		return s.PrivateKey
	}
	return *pathKey
}

func (s *Server) GetPassphrase(passphrase *string) string {
	if s.Passphrase != "" {
		return s.Passphrase
	}
	return *passphrase
}

func (s *Server) GetIsPassphrase(isPassphrase bool) bool {
	if s.Passphrase != "" {
		return isPassphrase
	}

	return isPassphrase
}

func (s *Server) GetPassword(password *string) string {
	if s.Password != "" {
		return s.Password
	}
	return *password
}

func (s *Server) GetShell(globalShell *shell.Shell) shell.Shell {
	if s.Shell == nil && globalShell == nil {
		val := false
		return shell.Shell{Enabled: &val}
	}

	if s.Shell == nil {
		return *globalShell
	}

	if *s.Shell.Enabled && s.Shell.After == "" && s.Shell.Before == "" {
		return *globalShell
	}
	return *s.Shell
}
