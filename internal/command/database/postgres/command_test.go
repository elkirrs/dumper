package postgres_test

import (
	"dumper/internal/command/database/postgres"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPSQLGenerator_Generate_AllScenarios(t *testing.T) {
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Plain SQL dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Format:   "plain",
					User:     "postgres",
					Password: "pass",
					Port:     "5432",
					Name:     "testdb",
				},
				DumpName:     "plain1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"pg_dump",
				"-Fp",
				"--dbname=postgresql://postgres:pass@127.0.0.1:5432/testdb",
			},
			expectedExt: "sql",
		},
		{
			name: "Plain SQL with gzip archive, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Format:   "plain",
					User:     "postgres",
					Password: "pass",
					Port:     "5433",
					Name:     "archive_db",
				},
				DumpName:     "plain_gzip",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"| gzip",
				"-Fp",
			},
			expectedExt: "sql.gz",
		},
		{
			name: "Custom format dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Format:   "dump",
					User:     "admin",
					Password: "pwd123",
					Port:     "5432",
					Name:     "customdb",
				},
				DumpName:     "custom1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"-Fc",
				"customdb",
			},
			expectedExt: "dump",
		},
		{
			name: "Tar format dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Format:   "tar",
					User:     "admin",
					Password: "pwd",
					Port:     "5432",
					Name:     "tardb",
				},
				DumpName:     "tar1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"-Ft",
				"tardb",
			},
			expectedExt: "tar",
		},
		{
			name: "Plain SQL dump to server",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Format:   "plain",
					User:     "postgres",
					Password: "pass",
					Port:     "5432",
					Name:     "serverdb",
				},
				DumpName:     "server_plain",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"> ./server_plain.sql",
				"-Fp",
			},
			expectedExt: "sql",
		},
		{
			name: "Plain SQL gzip to server",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Format:   "plain",
					User:     "postgres",
					Password: "pass",
					Port:     "5432",
					Name:     "serverdb",
				},
				DumpName:     "server_gzip",
				Archive:      true,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"| gzip",
				"> ./server_gzip.sql.gz",
			},
			expectedExt: "sql.gz",
		},
	}

	gen := postgres.PSQLGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)
			require.NoError(t, err, "Generate() should not return an error")
			require.NotNil(t, cmd, "DBCommand should not be nil")

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command, "Command string must not be empty")
			assert.NotEmpty(t, cmd.DumpPath, "DumpPath must not be empty")

			for _, expected := range tt.expectedContains {
				assert.Contains(t, cmd.Command, expected, "Command should contain %s", expected)
			}

			assert.Contains(t, cmd.DumpPath, tt.expectedExt, "DumpPath should contain extension %s", tt.expectedExt)
		})
	}
}
