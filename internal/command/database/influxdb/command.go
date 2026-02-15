package influxdb

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/pkg/utils/template"
	"fmt"
)

type Generator struct{}

type GenCommand struct {
	BaseCommand string
	ArchivePath string
}

func (g *Generator) Generate(
	data *cmdCfg.Config,
) (*commandDomain.DBCommand, error) {

	genCommand := getCommandByVersion(data)

	return &commandDomain.DBCommand{
		Command:  genCommand.BaseCommand,
		DumpPath: genCommand.ArchivePath,
	}, nil
}

func getCommandByVersion(data *cmdCfg.Config) *GenCommand {
	switch data.Database.Options.Version {
	case "2.x":
		return generateComandVersion2x(data)
	default:
		return generateComandVersion3x(data)
	}
}

func generateComandVersion2x(data *cmdCfg.Config) *GenCommand {
	backupPath := fmt.Sprintf("%s", data.DumpName)

	baseCmd := fmt.Sprintf("%s backup --host %s:%s --token %s",
		data.Database.Options.Source,
		data.Database.Options.Host,
		data.Database.Port,
		data.Database.Token,
	)

	if data.Database.Options.Bucket != "" {
		baseCmd += fmt.Sprintf(" --bucket %s", data.Database.Options.Bucket)
	}
	if data.Database.Options.BucketId != "" {
		baseCmd += fmt.Sprintf(" --bucket-id %s", data.Database.Options.BucketId)
	}

	if data.Database.Options.Organization != "" {
		baseCmd += fmt.Sprintf(" --org %s", data.Database.Options.Organization)
	}
	if data.Database.Options.OrganizationId != "" {
		baseCmd += fmt.Sprintf(" --org-id %s", data.Database.Options.OrganizationId)
	}

	if *data.Database.Options.SkipVerify {
		baseCmd += " --skip-verify"
	}

	if data.Database.Options.Start != "" {
		baseCmd += fmt.Sprintf(" --start %s", data.Database.Options.Start)
	}
	if data.Database.Options.End != "" {
		baseCmd += fmt.Sprintf(" --end %s", data.Database.Options.End)
	}
	if data.Database.Options.Filter != "" {
		baseCmd += fmt.Sprintf(" --filter %s", data.Database.Options.Filter)
	}

	baseCmd += fmt.Sprintf(" %s", backupPath)

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

	return &GenCommand{
		BaseCommand: baseCmd,
		ArchivePath: archivePath,
	}
}

func generateComandVersion3x(data *cmdCfg.Config) *GenCommand {
	backupPath := fmt.Sprintf("%s", data.DumpName)
	archivePath := backupPath + ".tar.gz"

	dataDirPath := template.GetFullPath(data.Database.Options.DataDir, data.Database.Options.NodeId)
	baseCmd := fmt.Sprintf("cp -r %s/snapshots %s", dataDirPath, backupPath)
	baseCmd = fmt.Sprintf("%s && cp -r %s/dbs %s", baseCmd, dataDirPath, backupPath)
	baseCmd = fmt.Sprintf("%s && cp -r %s/wal %s", baseCmd, dataDirPath, backupPath)
	baseCmd = fmt.Sprintf("%s && cp -r %s/catalog %s", baseCmd, dataDirPath, backupPath)
	baseCmd = fmt.Sprintf("%s && cp %s/_catalog_checkpoint %s", baseCmd, dataDirPath, backupPath)

	baseCmd = fmt.Sprintf("%s && tar -czf %s -C %s %s",
		baseCmd,
		archivePath,
		data.DumpDirRemote,
		data.DumpNameTemplate,
	)

	if data.RemoveBackup {
		baseCmd = fmt.Sprintf("%s && rm -rf %s", baseCmd, backupPath)
	}

	return &GenCommand{
		BaseCommand: baseCmd,
		ArchivePath: archivePath,
	}
}
