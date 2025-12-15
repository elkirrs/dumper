package redis

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type RedisGenerator struct{}

func (g RedisGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "rdb"

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("%s", fileName)

	host := "127.0.0.1"

	baseCmd := fmt.Sprintf(
		"%s -h %s -p %s -a %s --rdb",
		data.Database.Options.Source, host, data.Database.Port, data.Database.Password,
	)

	if data.Database.Options.Mode == "save" {
		saveCmd := fmt.Sprintf(
			"%s -h %s -p %s -a %s SAVE && ",
			data.Database.Options.Source, host, data.Database.Port, data.Database.Password,
		)
		baseCmd = saveCmd + baseCmd
	}

	if data.Archive {
		remotePath += ".gz"
		baseCmd = fmt.Sprintf("%s - | gzip > %s", baseCmd, remotePath)
	} else {
		baseCmd = fmt.Sprintf("%s %s", baseCmd, remotePath)
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
