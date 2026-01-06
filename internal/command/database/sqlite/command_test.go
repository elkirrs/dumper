package sqlite_test

import (
	"dumper/internal/command/database/sqlite"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils/mapping"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteGenerator_Generate_AllScenarios(t *testing.T) {
	source := mapping.GetDBSource("sqlite", "")
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Standard SQLite dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Options: option.Options{Path: "test.db", Source: source},
				},
				DumpName:     "dump1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlite3 test.db .dump",
				"> dump1.sql",
			},
			expectedExt: "sql",
		},
		{
			name: "SQLite dump with archive, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Options: option.Options{Path: "archive.db", Source: source},
				},
				DumpName:     "archive1",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlite3 archive.db .dump",
				"| gzip",
				"archive1.sql.gz",
			},
			expectedExt: "sql.gz",
		},
		{
			name: "SQLite dump to server, no archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Options: option.Options{Path: "server.db", Source: source},
				},
				DumpName:     "serverDump",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"sqlite3 server.db .dump",
				"> serverDump.sql",
			},
			expectedExt: "sql",
		},
		{
			name: "SQLite dump to server with archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Options: option.Options{Path: "serverArchive.db", Source: source},
				},
				DumpName:     "serverArchive",
				Archive:      true,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"sqlite3 serverArchive.db .dump",
				"| gzip",
				"serverArchive.sql.gz",
			},
			expectedExt: "sql.gz",
		},
	}

	gen := sqlite.SQLiteGenerator{}

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
