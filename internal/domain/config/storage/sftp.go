package storage

type SFTP struct {
	Dir        string `yaml:"dir" validate:"required"`
	Host       string `yaml:"host" validate:"required"`
	Port       string `yaml:"port" validate:"required"`
	Username   string `yaml:"username" validate:"required"`
	PrivateKey string `yaml:"private_key" validate:"required"`
	Passphrase string `yaml:"passphrase"`
}
