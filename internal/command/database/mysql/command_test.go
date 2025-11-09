package mysql_test

import (
	"dumper/internal/command/database/mysql"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQLGenerator_Generate_AllScenarios(t *testing.T) {
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Standard MySQL dump (no archive, local)",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "3306",
					User:     "root",
					Password: "password",
					Name:     "testdb",
				},
				DumpName:     "dump1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"mysqldump",
				"-h 127.0.0.1",
				"-u root",
				"-ppassword",
				"testdb",
			},
			expectedExt: "sql",
		},
		{
			name: "Archived MySQL dump (gzip, local)",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "3307",
					User:     "admin",
					Password: "1234",
					Name:     "prod_db",
				},
				DumpName:     "backup1",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"mysqldump",
				"| gzip",
				"prod_db",
			},
			expectedExt: "sql.gz",
		},
		{
			name: "Server dump (no archive)",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "3308",
					User:     "sa",
					Password: "pw123",
					Name:     "serverdb",
				},
				DumpName:     "dump_server",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"> ./dump_server.sql",
				"serverdb",
			},
			expectedExt: "sql",
		},
		{
			name: "Server dump (with archive)",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "3309",
					User:     "admin",
					Password: "pw",
					Name:     "archivedb",
				},
				DumpName:     "dump_archive",
				Archive:      true,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"| gzip",
				"> ./dump_archive.sql.gz",
				"archivedb",
			},
			expectedExt: "sql.gz",
		},
	}

	gen := mysql.MySQLGenerator{}

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
