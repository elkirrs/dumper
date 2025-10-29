package mssql

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type MSQLGenerator struct{}

func (g MSQLGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {

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
			"%s && powershell Compress-Archive -Path '%s' -DestinationPath '%s.gz'",
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
	remotePath := fmt.Sprintf("./%s", fileName)

	baseCmd := fmt.Sprintf(
		"sqlcmd -S %s -U %s -P %s -Q \"BACKUP DATABASE [%s] TO DISK='%s' WITH FORMAT, INIT\"",
		data.Server.Host,
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
	remotePath := fmt.Sprintf("./%s", fileName)

	baseCmd := fmt.Sprintf(
		"sqlpackage /Action:Export /SourceServerName:%s /SourceDatabaseName:%s /SourceUser:%s /SourcePassword:%s /TargetFile:%s",
		data.Server.Host,
		data.Database.Name,
		data.Database.User,
		data.Database.Password,
		remotePath,
	)

	if *data.Database.Options.SSL {
		baseCmd += " /SourceTrustServerCertificate:True"
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: remotePath,
	}
}
