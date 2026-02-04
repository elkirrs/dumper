package db2

import (
	commandDomain "dumper/internal/domain/command"
	commandConfig "dumper/internal/domain/command-config"
	"fmt"
)

type Generator struct{}

func (g *Generator) Generate(data *commandConfig.Config) (*commandDomain.DBCommand, error) {

	backupPath := fmt.Sprintf("%s", data.DumpName)
	dirToCreateDump := data.DumpName

	var backupMode string

	if data.Database.Options.BackupMode == "online" {
		backupMode = " online"
	}

	baseCmd := fmt.Sprintf(
		"%s backup database %s%s to %s",
		data.Database.Options.Source,
		data.Database.Name,
		backupMode,
		backupPath,
	)

	baseCmd = fmt.Sprintf("mkdir -p %s && %s", dirToCreateDump, baseCmd)
	archivePath := backupPath + ".tar.gz"

	baseCmd = fmt.Sprintf("%s && tar -czf %s -C %s %s",
		baseCmd,
		archivePath,
		data.DumpDirRemote,
		data.DumpNameTemplate,
	)

	if data.RemoveBackup {
		baseCmd = fmt.Sprintf("%s && rm -rf %s", baseCmd, backupPath)
	}

	if data.DumpLocation == "server" {
		return &commandDomain.DBCommand{
			Command:  baseCmd,
			DumpPath: archivePath,
		}, nil
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: archivePath,
	}, nil
}
