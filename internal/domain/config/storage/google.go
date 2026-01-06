package storage

type GoogleCloud struct {
	Type           string `yaml:"type" validate:"required"`
	Bucket         string `yaml:"bucket" validate:"required"`
	Credential     string `yaml:"credential"`
	CredentialFile string `yaml:"credential_file"`
}
