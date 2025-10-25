package server

type Server struct {
	Host       string `yaml:"host" validate:"required"`
	User       string `yaml:"user" validate:"required"`
	Name       string `yaml:"name,omitempty"`
	Port       string `yaml:"port,omitempty"`
	SSHKey     string `yaml:"key,omitempty"`
	Password   string `yaml:"password,omitempty"`
	ConfigPath string `yaml:"conf_path,omitempty"`
}

func (s Server) GetName() string {
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
