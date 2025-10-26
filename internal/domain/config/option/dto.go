package option

type Options struct {
	AuthSource string `yaml:"auth_source"`
	SSL        *bool  `yaml:"ssl" default:"false"`
	Mode       string `yaml:"mode"`
	Role       string `yaml:"role"`
	Path       string `yaml:"path"`
}
