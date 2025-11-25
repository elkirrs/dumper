package ssh_config

type SSHConfig struct {
	PrivateKey   string `yaml:"private_key"`
	Passphrase   string `yaml:"passphrase"`
	IsPassphrase *bool  `yaml:"is_passphrase" default:"false"`
	Password     string `yaml:"password"`
}
