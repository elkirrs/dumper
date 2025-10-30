package app

type Flags struct {
	ConfigFile     string
	DbNameList     string
	All            bool
	FileLog        string
	Input          string
	Crypt          string
	Password       string
	Mode           string
	Recovery       string
	AppSecret      string
	OpenOnlyEncEnv bool
}
