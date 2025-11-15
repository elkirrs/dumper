package shell

type Shell struct {
	After   string `yaml:"after"`
	Before  string `yaml:"before"`
	Enabled bool   `yaml:"enabled" default:"true"`
}
