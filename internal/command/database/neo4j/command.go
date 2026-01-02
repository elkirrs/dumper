package neo4j

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/pkg/utils"
	"fmt"
)

type Neo4jGenerator struct{}

func (g Neo4jGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "dump"

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("%s", fileName)

	// Neo4j dump command using neo4j-admin
	// Format: neo4j-admin database dump <database-name> --to-path=<path>
	baseCmd := fmt.Sprintf(
		"%s database dump %s --to-path=%s",
		data.Database.Options.Source,
		data.Database.Name,
		data.DumpDirRemote,
	)

	// TODO Need add Neo4j Enterprise hot backup

	// Move the dump file to the desired name
	nameTmpDumpFile := fmt.Sprintf("%s.%s", data.Database.Name, ext)
	pathDumpFile := utils.GetFullPath(data.DumpDirRemote, nameTmpDumpFile)
	baseCmd = fmt.Sprintf("%s && mv %s %s", baseCmd, pathDumpFile, fileName)

	if data.Archive {
		remotePath += ".gz"
		baseCmd = fmt.Sprintf("%s && gzip -c %s > %s", baseCmd, fileName, remotePath)

		if data.RemoveBackup {
			baseCmd = fmt.Sprintf("%s && rm %s", baseCmd, fileName)
		}
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
