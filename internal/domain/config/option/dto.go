package option

type Options struct {
	AuthSource string   `yaml:"auth_source"`
	SSL        *bool    `yaml:"ssl" default:"false"`
	Mode       string   `yaml:"mode"`
	Role       string   `yaml:"role"`
	Path       string   `yaml:"path"`
	Source     string   `yaml:"source"`
	IncTables  []string `yaml:"inc_tables"`
	ExcTables  []string `yaml:"exc_tables"`
}
