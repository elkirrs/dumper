package encrypt

type Encrypt struct {
	Type     string `yaml:"type"`
	Password string `yaml:"password"`
	Enabled  *bool  `yaml:"enabled" default:"true"`
}
