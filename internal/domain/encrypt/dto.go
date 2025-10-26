package encrypt

type Options struct {
	FilePath string
	Password string
	Key      string
	Type     string
	Crypt    string
}

type DataCrypt struct {
	CMD  string
	Name string
}
