package firebird

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/pkg/utils/template"
	"fmt"
)

type Generator struct{}

func (g *Generator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "fbk"

	backupPath := fmt.Sprintf("%s.%s", data.DumpName, ext)
	backupFilePath := backupPath

	baseCmd := fmt.Sprintf("%s -b", data.Database.Options.Source)

	if data.Database.Options.SkipGarbage {
		baseCmd += " -g"
	}

	if data.Database.Options.SkipIssue {
		baseCmd += " -ignore"
	}

	if data.Database.Options.FastAndStable {
		baseCmd += fmt.Sprintf(" -se service_mgr %s", data.Database.Options.Path)
	} else {
		fileSource := fmt.Sprintf("%s/%s:/%s", "localhost", data.Database.Port, data.Database.Options.Path)
		fileSource = template.GetFullPath(fileSource)
		baseCmd += fmt.Sprintf(" %s", fileSource)
	}

	baseCmd += fmt.Sprintf(" %s -user %s -password %s", backupFilePath, data.Database.User, data.Database.Password)

	if data.Archive {
		originFileName := fmt.Sprintf("%s.%s", data.DumpNameTemplate, ext)
		backupFilePath = fmt.Sprintf("%s.tar.gz", data.DumpName)
		baseCmd += fmt.Sprintf(
			" && tar -czf %s -C %s %s",
			backupFilePath,
			data.DumpDirRemote,
			originFileName,
		)

		if data.RemoveBackup {
			baseCmd += fmt.Sprintf(" && rm -rf %s", backupPath)
		}
	}

	if data.DumpLocation == "server" {
		return &commandDomain.DBCommand{
			Command:  baseCmd,
			DumpPath: backupFilePath,
		}, nil
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: backupFilePath,
	}, nil
}
