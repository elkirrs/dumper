package dynamodb_test

import (
	"dumper/internal/command/database/dynamodb"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils/mapping"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDynamoDBGenerator_Generate_AllScenarios(t *testing.T) {
	source := mapping.GetDBSource("dynamodb", "")
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Standard DynamoDB dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "MyTable",
					Options: option.Options{
						Source: source,
						Region: "us-east-1",
					},
				},
				DumpName:     "dump1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"aws dynamodb scan",
				"--table-name MyTable",
				"--region us-east-1",
				"> dump1.json",
			},
			expectedExt: "json",
		},
		{
			name: "DynamoDB dump with profile",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "ProductsTable",
					Options: option.Options{
						Source:  source,
						Region:  "eu-west-1",
						Profile: "production",
					},
				},
				DumpName:     "products",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"aws dynamodb scan",
				"--table-name ProductsTable",
				"--region eu-west-1",
				"--profile production",
				"> products.json",
			},
			expectedExt: "json",
		},
		{
			name: "DynamoDB dump with archive, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "UsersTable",
					Options: option.Options{
						Source: source,
						Region: "us-west-2",
					},
				},
				DumpName:     "archive1",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"aws dynamodb scan",
				"--table-name UsersTable",
				"--region us-west-2",
				"| gzip",
				"archive1.json.gz",
			},
			expectedExt: "json.gz",
		},
		{
			name: "DynamoDB local endpoint dump",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "LocalTable",
					Options: option.Options{
						Source:   source,
						Region:   "local",
						Endpoint: "http://localhost:8000",
					},
				},
				DumpName:     "localDump",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"aws dynamodb scan",
				"--table-name LocalTable",
				"--region local",
				"--endpoint-url http://localhost:8000",
				"> localDump.json",
			},
			expectedExt: "json",
		},
		{
			name: "DynamoDB dump to server, no archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "ServerTable",
					Options: option.Options{
						Source: source,
						Region: "ap-southeast-1",
					},
				},
				DumpName:     "serverDump",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"aws dynamodb scan",
				"--table-name ServerTable",
				"--region ap-southeast-1",
				"> serverDump.json",
			},
			expectedExt: "json",
		},
		{
			name: "DynamoDB dump to server with archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "ArchiveTable",
					Options: option.Options{
						Source: source,
						Region: "us-east-1",
					},
				},
				DumpName:     "serverArchive",
				Archive:      true,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"aws dynamodb scan",
				"--table-name ArchiveTable",
				"--region us-east-1",
				"| gzip",
				"serverArchive.json.gz",
			},
			expectedExt: "json.gz",
		},
	}

	gen := dynamodb.Generator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)
			require.NoError(t, err, "Generate() should not return an error")
			require.NotNil(t, cmd, "DBCommand should not be nil")

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command, "Command string should not be empty")
			assert.NotEmpty(t, cmd.DumpPath, "DumpPath should not be empty")

			for _, expected := range tt.expectedContains {
				assert.Contains(t, cmd.Command, expected, "Command should contain %s", expected)
			}

			assert.Contains(t, cmd.DumpPath, tt.expectedExt, "DumpPath should contain extension %s", tt.expectedExt)
		})
	}
}

func TestDynamoDBGenerator_CommandIntegrity(t *testing.T) {
	source := mapping.GetDBSource("dynamodb", "")

	cfg := &cmdCfg.Config{
		Database: cmdCfg.Database{
			Name: "TestTable",
			Options: option.Options{
				Source: source,
				Region: "us-east-1",
			},
		},
		DumpName:     "backup",
		Archive:      false,
		DumpLocation: "local",
	}

	gen := dynamodb.Generator{}
	cmd, err := gen.Generate(cfg)
	require.NoError(t, err)
	require.NotNil(t, cmd)

	expectedPrefix := "aws dynamodb scan --table-name TestTable"
	assert.Contains(t, cmd.Command, expectedPrefix, "Command must be constructed correctly")
	assert.Equal(t, "backup.json", cmd.DumpPath, "DumpPath should match expected filename")
}
