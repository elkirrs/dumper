package postgres

import (
	command "dumper/internal/command/database"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type PSQLGenerator struct{}

func (g PSQLGenerator) Generate(data *cmdCfg.ConfigData, settings *setting.Settings) (string, string) {
	formatFlag := "-Fp" // plain SQL
	ext := "sql"

	switch data.DumpFormat {
	case "dump":
		formatFlag = "-Fc"
		ext = "dump"
	case "tar":
		formatFlag = "-Ft"
		ext = "tar"
	}

	baseCmd := fmt.Sprintf("/usr/bin/pg_dump --dbname=postgresql://%s:%s@127.0.0.1:%s/%s --clean --if-exists --no-owner %s",
		data.User, data.Password, data.Port, data.Name, formatFlag)

	if *settings.Archive && formatFlag == "-Fp" { // gzip only for plain
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
	command.Register("psql", PSQLGenerator{})
}
