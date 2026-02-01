package mariadb_test

import (
	"dumper/internal/command/database/mariadb"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils/mapping"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMariaDbGenerator_Generate(t *testing.T) {
	source := mapping.GetDBSource("mariadb", "")
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "No archive, local dump",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "root",
					Password: "123",
					Port:     "3306",
					Name:     "testdb",
					Options:  option.Options{Source: source},
				},
				DumpName:     "dump",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"mariadb-dump",
				"-uroot",
				"-p123",
				"-P3306",
				"testdb",
			},
			expectedExt: "sql",
		},
		{
			name: "Archive enabled, local dump",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "root",
					Password: "123",
					Port:     "3306",
					Name:     "testdb",
					Options:  option.Options{Source: source},
				},
				DumpName:     "dump",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"mariadb-dump",
				"gzip",
			},
			expectedExt: "sql.gz",
		},
		{
			name: "No archive, dump to server",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "root",
					Password: "123",
					Port:     "3306",
					Name:     "testdb",
					Options:  option.Options{Source: source},
				},
				DumpName:     "server_dump",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"mariadb-dump",
				"> server_dump.sql",
			},
			expectedExt: "sql",
		},
		{
			name: "Archive enabled, dump to server",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "root",
					Password: "123",
					Port:     "3306",
					Name:     "testdb",
					Options:  option.Options{Source: source},
				},
				DumpName:     "server_dump",
				Archive:      true,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"mariadb-dump",
				"gzip",
				"> server_dump.sql.gz",
			},
			expectedExt: "sql.gz",
		},
		{
			name: "Dump with only table from list (include tables), dump to server",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "root",
					Password: "123",
					Port:     "3306",
					Name:     "testdb",
					Options: option.Options{
						Source:    source,
						IncTables: []string{"users", "roles"},
					},
				},
				DumpName:     "server_dump",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"mariadb-dump",
				"-uroot -p123 -h127.0.0.1 -P3306 testdb ",
				"users roles",
				"> server_dump.sql",
			},
			expectedExt: "sql",
		},
		{
			name: "Dump with without table from list (exclude tables), dump to server",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "root",
					Password: "123",
					Port:     "3306",
					Name:     "testdb",
					Options: option.Options{
						Source:    source,
						ExcTables: []string{"users", "roles"},
					},
				},
				DumpName:     "server_dump",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"mariadb-dump",
				"-uroot -p123 -h127.0.0.1 -P3306 testdb ",
				"--ignore-table=testdb.users --ignore-table=testdb.roles",
				"> server_dump.sql",
			},
			expectedExt: "sql",
		},
	}

	gen := mariadb.Generator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)
			require.NoError(t, err, "Generate should not return an error")
			require.NotNil(t, cmd, "Returned command must not be nil")

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command, "Command string should not be empty")
			assert.NotEmpty(t, cmd.DumpPath, "DumpPath should not be empty")

			for _, expect := range tt.expectedContains {
				assert.Contains(t, cmd.Command, expect, "Command should contain expected fragment: %s", expect)
			}

			assert.Contains(t, cmd.DumpPath, tt.expectedExt, "DumpPath should have expected extension")
		})
	}
}

func TestMariaDbGenerator_CommandIntegrity(t *testing.T) {
	source := mapping.GetDBSource("mariadb", "")

	cfg := &cmdCfg.Config{
		Database: cmdCfg.Database{
			User:     "admin",
			Password: "secret",
			Port:     "3307",
			Name:     "mydb",
			Options:  option.Options{Source: source},
		},
		DumpName:     "backup",
		Archive:      false,
		DumpLocation: "local",
	}

	gen := mariadb.Generator{}
	cmd, err := gen.Generate(cfg)
	require.NoError(t, err)
	require.NotNil(t, cmd)

	expectedPrefix := "mariadb-dump -uadmin -psecret -h127.0.0.1 -P3307 mydb"
	assert.Contains(t, cmd.Command, expectedPrefix, "Command must be constructed correctly")
	assert.Equal(t, "backup.sql", cmd.DumpPath, "DumpPath should match expected filename")
}
