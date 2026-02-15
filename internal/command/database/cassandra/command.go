package cassandra

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
	"strings"
)

type Generator struct{}

func (g *Generator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := ".tar.gz"

	baseCmd := fmt.Sprintf("%s snapshot -t %s",
		data.Database.Options.Source,
		data.DumpNameTemplate,
	)

	archiveName := fmt.Sprintf("%s%s", data.DumpName, ext)

	preparedTables, preparedArchive := prepare(data, archiveName)
	if preparedTables != "" {
		baseCmd = fmt.Sprintf("%s %s", baseCmd, preparedTables)
	} else {
		baseCmd = fmt.Sprintf("%s %s", baseCmd, data.Database.Name)
	}

	baseCmd = fmt.Sprintf("%s && cd %s && %s",
		baseCmd,
		data.DumpDirRemote,
		preparedArchive,
	)

	if data.RemoveBackup {
		removeSnapshot := fmt.Sprintf("%s clearsnapshot -t %s",
			data.Database.Options.Source,
			data.DumpNameTemplate,
		)
		baseCmd = fmt.Sprintf("%s && %s", baseCmd, removeSnapshot)
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: archiveName,
	}, nil
}

func prepare(
	data *cmdCfg.Config,
	archiveName string,
) (string, string) {
	var tables string
	var tablesArr []string
	var archive string
	var archiveArr []string

	if len(data.Database.Options.IncTables) > 0 {
		tables += "--kt-list"
		archive += fmt.Sprintf("tar -czf %s $(find ./%s", archiveName, data.Database.Name)
		for _, table := range data.Database.Options.IncTables {
			tablesArr = append(tablesArr, fmt.Sprintf("%s.%s", data.Database.Name, table))
			archiveArr = append(archiveArr, fmt.Sprintf("-path \"./%s/%s-*\"", data.Database.Name, table))
		}

		tables += " " + strings.Join(tablesArr, ",")
		archive += fmt.Sprintf(" \\( %s \\) -path \"*/snapshots/*\")", strings.Join(archiveArr, " -o "))
	} else {
		tables = ""
		archive = fmt.Sprintf("tar -czf %s $(find ./%s -path \"*/snapshots/*\")", archiveName, data.Database.Name)
	}

	return tables, archive
}
