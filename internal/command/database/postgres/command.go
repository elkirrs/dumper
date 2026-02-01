package postgres

import (
	commandDomain "dumper/internal/domain/command"
	commandConfig "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"fmt"
)

type Generator struct{}

func (g Generator) Generate(data *commandConfig.Config) (*commandDomain.DBCommand, error) {
	formatFlag := "-Fp" // plain SQL
	ext := "sql"

	switch data.Database.Format {
	case "dump":
		formatFlag = "-Fc"
		ext = "dump"
	case "tar":
		formatFlag = "-Ft"
		ext = "tar"
	}

	tables := prepareTables(&data.Database.Options)

	baseCmd := fmt.Sprintf(
		"%s --dbname=postgresql://%s:%s@127.0.0.1:%s/%s %s --clean --if-exists --no-owner %s",
		data.Database.Options.Source,
		data.Database.User,
		data.Database.Password,
		data.Database.Port,
		data.Database.Name,
		tables,
		formatFlag,
	)

	if data.Archive && formatFlag == "-Fp" { // gzip only for plain
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
) string {
	out := ""

	if options.IncTables != nil {
		for _, table := range options.IncTables {
			out += fmt.Sprintf(" %s%s%s", "--table=", "public.", table)
		}
	}

	if options.IncTables == nil && options.ExcTables != nil {
		for _, table := range options.ExcTables {
			out += fmt.Sprintf(" %s%s%s", "--exclude-table=", "public.", table)
		}
	}

	return out
}
