package mariadb

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"fmt"
)

type Generator struct{}

func (g *Generator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "sql"

	baseCmd := fmt.Sprintf(
		"%s -u%s -p%s -h127.0.0.1 -P%s %s",
		data.Database.Options.Source,
		data.Database.User,
		data.Database.Password,
		data.Database.Port,
		data.Database.Name,
	)

	tables := prepareTables(
		&data.Database.Options,
		data.Database.Name+".",
	)

	if tables != "" {
		baseCmd = fmt.Sprintf("%s %s", baseCmd, tables)
	}

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

func prepareTables(
	options *option.Options,
	prefix string,
) string {
	out := ""

	if options.IncTables != nil {
		for _, table := range options.IncTables {
			out += fmt.Sprintf(" %s", table)
		}
	}

	if options.IncTables == nil && options.ExcTables != nil {
		for _, table := range options.ExcTables {
			out += fmt.Sprintf(" %s%s%s", "--ignore-table=", prefix, table)
		}
	}

	return out
}
