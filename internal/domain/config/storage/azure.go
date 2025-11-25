package storage

type AzureAD struct {
	Type         string `yaml:"type" validate:"required"`
	TenantID     string `yaml:"tenant_id" validate:"required"`
	ClientID     string `yaml:"client_id" validate:"required"`
	ClientSecret string `yaml:"client_secret" validate:"required"`
	Endpoint     string `yaml:"endpoint" validate:"required,url"`
	Container    string `yaml:"container" validate:"required"`
}

type AzureSharedKey struct {
	Type      string `yaml:"type" validate:"required"`
	Name      string `yaml:"name" validate:"required"`
	SharedKey string `yaml:"shared_key" validate:"required"`
	Endpoint  string `yaml:"endpoint" validate:"required,url"`
	Container string `yaml:"container" validate:"required"`
}
