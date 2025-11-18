package docker

type Docker struct {
	Command string `yaml:"command"`
	Enabled *bool  `yaml:"enabled" default:"false"`
}
