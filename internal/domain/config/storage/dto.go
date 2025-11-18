package storage

type Storage struct {
	Type       string
	Dir        string `yaml:"dir"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	PrivateKey string `yaml:"private_key"`
	Passphrase string `yaml:"passphrase"`

	TenantID     string `yaml:"tenant_id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Endpoint     string `yaml:"endpoint"`
	Container    string `yaml:"container"`

	//Region    string `yaml:"region"`
	//Bucket    string `yaml:"bucket"`
	//AccessKey string `yaml:"access_key"`
	//SecretKey string `yaml:"secret_key"`
}

type ListStorages struct {
	Type    string
	Configs Storage
}

func (s Storage) GetPrivateKey(pathKey string) string {
	if s.PrivateKey != "" {
		return s.PrivateKey
	}
	return pathKey
}
