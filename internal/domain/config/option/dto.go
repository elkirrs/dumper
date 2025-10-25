package option

type Options struct {
	AuthSource string `yaml:"auth_source"`
	SSL        *bool  `yaml:"ssl" default:"false"`
	Mode       string `yaml:"mode"`
}
