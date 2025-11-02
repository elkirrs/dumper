package server

type Server struct {
	Title      string `yaml:"title,omitempty"`
	Host       string `yaml:"host" validate:"required"`
	User       string `yaml:"user" validate:"required"`
	Name       string `yaml:"name,omitempty"`
	Port       string `yaml:"port,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty"`
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

func (s Server) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	if s.Name != "" {
		return s.Name
	}

	return s.Host
}

func (s Server) GetPrivateKey(pathKey string) string {
	if s.PrivateKey != "" {
		return s.PrivateKey
	}
	return pathKey
}
