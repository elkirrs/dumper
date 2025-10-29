package sqlite

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type SQLiteGenerator struct{}

func (g SQLiteGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "sql"

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	baseCmd := fmt.Sprintf("sqlite3 %s .dump", data.Database.Options.Path)

	if data.Archive {
		remotePath += ".gz"
		baseCmd = fmt.Sprintf("%s | gzip > %s", baseCmd, remotePath)
	} else {
		baseCmd = fmt.Sprintf("%s > %s", baseCmd, remotePath)
	}

	if data.DumpLocation == "server" {
		return &commandDomain.DBCommand{
			Command:  baseCmd,
			DumpPath: remotePath,
		}, nil
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: remotePath,
	}, nil
}
