package mysql

import (
	"dumper/internal/command"
	"dumper/internal/config"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type MSQLGenerator struct{}

func (g MSQLGenerator) Generate(data *cmdCfg.ConfigData, settings *config.Settings) (string, string) {
	if data.Port == "" {
		data.Port = "3306"
	}
	return fmt.Sprintf("mysqldump -u%s -p%s -h127.0.0.1 -P%s %s",
		data.User, data.Password, data.Port, data.Name), ""
}

func init() {
	command.Register("mysql", MSQLGenerator{})
}
