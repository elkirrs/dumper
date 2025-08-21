package mysql

import (
	"dumper/internal/command"
	"dumper/internal/config"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type MySQLGenerator struct{}

func (g MySQLGenerator) Generate(data *cmdCfg.ConfigData, settings *config.Settings) (string, string) {
	if data.Port == "" {
		data.Port = "3306"
	}

	ext := "sql"

	baseCmd := fmt.Sprintf("/usr/bin/mysqldump -h 127.0.0.1 -P %s -u %s -p%s %s",
		data.Port, data.User, data.Password, data.Name)

	if *settings.Archive {
		baseCmd += " | gzip"
		ext += ".gz"
	}

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	if settings.DumpLocation == "server" {
		return fmt.Sprintf("%s > %s", baseCmd, remotePath), remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("mysql", MySQLGenerator{})
}
