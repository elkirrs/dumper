package mssql

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"fmt"
)

type Generator struct{}

func (g *Generator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {

	var baseCmd *commandDomain.DBCommand

	switch data.Database.Format {
	case "bac":
		baseCmd = genFormatBak(data)
	case "bacpac":
		baseCmd = genFormatBacpac(data)
	default:
		return nil, fmt.Errorf("unsupported database format: %s", data.Database.Format)
	}

	if data.Archive {
		baseCmd.Command = fmt.Sprintf(
			"%s && tar -czf %s.gz %s",
			baseCmd.Command, baseCmd.DumpPath, baseCmd.DumpPath,
		)
		baseCmd.DumpPath += ".gz"
	}

	if data.DumpLocation == "server" {
		return baseCmd, nil
	}

	return baseCmd, nil
}

func genFormatBak(
	data *cmdCfg.Config,
) *commandDomain.DBCommand {
	ext := "bak"
	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("%s", fileName)

	baseCmd := fmt.Sprintf(
		"%s -S %s -C -U %s -P %s -Q \"BACKUP DATABASE [%s] TO DISK=\\\"%s\\\" WITH FORMAT, INIT\"",
		data.Database.Options.Source,
		"localhost",
		data.Database.User,
		data.Database.Password,
		data.Database.Name,
		remotePath,
	)

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: remotePath,
	}
}

func genFormatBacpac(
	data *cmdCfg.Config,
) *commandDomain.DBCommand {
	ext := "bacpac"
	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("%s", fileName)

	baseCmd := fmt.Sprintf(
		"%s /Action:Export /SourceServerName:%s /SourceDatabaseName:%s /SourceUser:%s /SourcePassword:%s /TargetFile:%s",
		data.Database.Options.Source,
		"localhost",
		data.Database.Name,
		data.Database.User,
		data.Database.Password,
		remotePath,
	)

	tables := prepareTables(&data.Database.Options)

	if tables != "" {
		baseCmd = fmt.Sprintf("%s %s", baseCmd, tables)
	}

	if !*data.Database.Options.SSL {
		baseCmd += " /SourceTrustServerCertificate:True"
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: remotePath,
	}
}

func prepareTables(
	options *option.Options,
) string {
	out := ""

	if options.IncTables != nil {
		for _, table := range options.IncTables {
			out += fmt.Sprintf(" /p:TableData=\"dbo.%s\"", table)
		}
	}

	return out
}
