package neo4j

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
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
		"%s database dump %s --to-path=.",
		data.Database.Options.Source,
		data.Database.Name,
	)

	// Move the dump file to the desired name
	baseCmd = fmt.Sprintf("%s && mv %s.dump %s", baseCmd, data.Database.Name, fileName)

	if data.Archive {
		remotePath += ".gz"
		baseCmd = fmt.Sprintf("%s && gzip -c %s > %s && rm %s", baseCmd, fileName, remotePath, fileName)
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
