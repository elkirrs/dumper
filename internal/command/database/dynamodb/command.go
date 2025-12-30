package dynamodb

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type DynamoDBGenerator struct{}

func (g DynamoDBGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "json"

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("%s", fileName)

	// DynamoDB backup using AWS CLI
	// We'll use scan to export data to JSON format
	baseCmd := fmt.Sprintf(
		"%s dynamodb scan --table-name %s",
		data.Database.Options.Source,
		data.Database.Name,
	)

	// Add region if specified
	if data.Database.Options.Region != "" {
		baseCmd = fmt.Sprintf("%s --region %s", baseCmd, data.Database.Options.Region)
	}

	// Add AWS profile if specified
	if data.Database.Options.Profile != "" {
		baseCmd = fmt.Sprintf("%s --profile %s", baseCmd, data.Database.Options.Profile)
	}

	// Add endpoint URL if specified (for local DynamoDB)
	if data.Database.Options.Endpoint != "" {
		baseCmd = fmt.Sprintf("%s --endpoint-url %s", baseCmd, data.Database.Options.Endpoint)
	}

	if data.Archive {
		remotePath += ".gz"
		baseCmd = fmt.Sprintf("%s | gzip > %s", baseCmd, remotePath)
	} else {
		baseCmd = fmt.Sprintf("%s > %s", baseCmd, remotePath)
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
