package storage

type Storage struct {
	// Common
	Type string `yaml:"type" validate:"required"`

	// Local
	Dir string `yaml:"dir"`

	// FTP / SFTP / Common remote
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	// SFTP only
	PrivateKey string `yaml:"private_key"`
	Passphrase string `yaml:"passphrase"`

	// Azure Common
	Endpoint  string `yaml:"endpoint"`
	Container string `yaml:"container"`
	AuthType  string `yaml:"auth_type" default:"SharedKey"`

	// Azure AD
	TenantID     string `yaml:"tenant_id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	// Azure SharedKey
	Name      string `yaml:"name"`
	SharedKey string `yaml:"shared_key"`

	// S3
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`

	// Cloudflare
	AccountID string `yaml:"account_id"`

	// Google Cloud
	Credential     string `yaml:"credential"`
	CredentialFile string `yaml:"credential_file"`
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
