package storage

type Local struct {
	Type string `yaml:"type" validate:"required"`
	Dir  string `yaml:"dir" validate:"required"`
}
