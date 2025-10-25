package mssql

import (
	"dumper/internal/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type MSSQLGenerator struct{}

func (g MSSQLGenerator) Generate(data *cmdCfg.ConfigData, settings *setting.Settings) (string, string) {

	var baseCmd, remotePath string

	switch data.DumpFormat {
	case "bac":
		baseCmd, remotePath = genFormatBak(data)
	case "bacpac":
		baseCmd, remotePath = genFormatBacpac(data)
	}

	if *settings.Archive {
		baseCmd = fmt.Sprintf(
			"%s && powershell Compress-Archive -Path '%s' -DestinationPath '%s.gz'",
			baseCmd, remotePath, remotePath,
		)
		remotePath += ".gz"
	}

	if settings.DumpLocation == "server" {
		return baseCmd, remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("mssql", MSSQLGenerator{})
}

func genFormatBak(
	data *cmdCfg.ConfigData,
) (string, string) {
	ext := "bak"
	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	baseCmd := fmt.Sprintf(
		"sqlcmd -S %s -U %s -P %s -Q \"BACKUP DATABASE [%s] TO DISK='%s' WITH FORMAT, INIT\"",
		data.Name, data.User, data.Password, data.Name, remotePath,
	)

	return baseCmd, remotePath
}

func genFormatBacpac(
	data *cmdCfg.ConfigData,
) (string, string) {
	ext := "bacpac"
	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	baseCmd := fmt.Sprintf(
		"sqlpackage /Action:Export /SourceServerName:%s /SourceDatabaseName:%s /SourceUser:%s /SourcePassword:%s /TargetFile:%s",
		data.Host, data.Name, data.User, data.Password, remotePath,
	)

	if *data.Options.SSL {
		baseCmd += " /SourceTrustServerCertificate:True"
	}

	return baseCmd, remotePath
}
