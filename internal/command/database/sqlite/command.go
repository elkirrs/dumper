package sqlite

import (
	command "dumper/internal/command/database"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type SQLiteGenerator struct{}

func (g SQLiteGenerator) Generate(data *cmdCfg.ConfigData, settings *setting.Settings) (string, string) {
	ext := "sql"

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	baseCmd := fmt.Sprintf("sqlite3 %s .dump", data.Options.Path)

	if *settings.Archive {
		remotePath += ".gz"
		baseCmd = fmt.Sprintf("%s | gzip > %s", baseCmd, remotePath)
	} else {
		baseCmd = fmt.Sprintf("%s > %s", baseCmd, remotePath)
	}

	if settings.DumpLocation == "server" {
		return baseCmd, remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("sqlite", SQLiteGenerator{})
}
