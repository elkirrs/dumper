package storage

type Storage struct {
	Type       string
	Dir        string `json:"dir"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Passphrase string `json:"passphrase"`

	//TenantID     string `json:"tenant_id"`
	//ClientID     string `json:"client_id"`
	//ClientSecret string `json:"client_secret"`
	//Endpoint     string `json:"endpoint"`
	//Container    string `json:"container"`

	//Region    string `json:"region"`
	//Bucket    string `json:"bucket"`
	//AccessKey string `json:"access_key"`
	//SecretKey string `json:"secret_key"`
}

type ListStorages struct {
	Type    string
	Configs Storage
}

func (s Storage) GetPrivateKey(pathKey string) string {
	if s.PrivateKey != "" {
		return s.PrivateKey
	}
	return pathKey
}
