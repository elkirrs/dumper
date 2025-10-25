package sqlite

import (
	"dumper/internal/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type SQLiteGenerator struct{}

func (g SQLiteGenerator) Generate(data *cmdCfg.ConfigData, settings *setting.Settings) (string, string) {
	ext := "sql"

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	baseCmd := fmt.Sprintf("sqlite3 %s \".output %s\" \".dump\"",
		data.Name, remotePath,
	)

	if *settings.Archive {
		baseCmd = fmt.Sprintf("%s | gzip", baseCmd)
		ext += ".gz"
		remotePath += ".gz"
	}

	if settings.DumpLocation == "server" {
		return baseCmd, remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("sqlite", SQLiteGenerator{})
}
