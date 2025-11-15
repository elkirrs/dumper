package mariadb

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type MariaDbGenerator struct{}

func (g MariaDbGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "sql"

	baseCmd := fmt.Sprintf(
		"mariadb-dump -u%s -p%s -h127.0.0.1 -P%s %s",
		data.Database.User,
		data.Database.Password,
		data.Database.Port,
		data.Database.Name,
	)

	if data.Archive {
		baseCmd += " | gzip"
		ext += ".gz"
	}

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("%s", fileName)

	if data.DumpLocation == "server" {
		return &commandDomain.DBCommand{
			Command:  fmt.Sprintf("%s > %s", baseCmd, remotePath),
			DumpPath: remotePath,
		}, nil
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: remotePath,
	}, nil
}
