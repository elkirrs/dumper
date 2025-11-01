package storage

type Storage struct {
	Type     string
	Dir      string `json:"dir"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ListStorages struct {
	Type    string
	Configs Storage
}
