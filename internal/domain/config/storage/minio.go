package storage

type MinIO struct {
	Type      string `yaml:"type" validate:"required"`
	Region    string `yaml:"region" validate:"required"`
	AccessKey string `yaml:"access_key" validate:"required"`
	SecretKey string `yaml:"secret_key" validate:"required"`
	Bucket    string `yaml:"bucket" validate:"required"`
	Endpoint  string `yaml:"endpoint" validate:"required"`
}
