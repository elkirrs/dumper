package cassandra_test

import (
	"dumper/internal/command/database/cassandra"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCassandraGenerator_Generate_Variations(t *testing.T) {

	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedArchive  string
	}{
		{
			name: "minimal config - full keyspace snapshot",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "testks",
					Options: option.Options{
						Source: "/opt/cassandra/bin/nodetool",
					},
				},
				DumpName:         "dump",
				DumpDirRemote:    "/opt/cassandra/data/data",
				DumpNameTemplate: "dump",
			},
			expectedContains: []string{
				"snapshot -t dump testks",
				"cd /opt/cassandra/data/data",
				`tar -czf dump.tar.gz $(find ./testks -path "*/snapshots/*")`,
			},
			expectedArchive: "dump.tar.gz",
		},

		{
			name: "include single table",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "metrics",
					Options: option.Options{
						Source:    "nodetool",
						IncTables: []string{"cpu"},
					},
				},
				DumpName:         "table_dump",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "table_dump",
			},
			expectedContains: []string{
				"snapshot -t table_dump --kt-list metrics.cpu",
				`-path "./metrics/cpu-*"`,
			},
			expectedArchive: "table_dump.tar.gz",
		},

		{
			name: "include multiple tables",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "logs",
					Options: option.Options{
						Source:    "nodetool",
						IncTables: []string{"events", "users"},
					},
				},
				DumpName:         "multi",
				DumpDirRemote:    "/backup",
				DumpNameTemplate: "multi",
			},
			expectedContains: []string{
				"--kt-list logs.events,logs.users",
				`-path "./logs/events-*"`,
				`-path "./logs/users-*"`,
			},
			expectedArchive: "multi.tar.gz",
		},

		{
			name: "remove snapshot enabled",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "cleanupks",
					Options: option.Options{
						Source: "nodetool",
					},
				},
				DumpName:         "cleanup",
				DumpDirRemote:    "/backup",
				DumpNameTemplate: "cleanup",
				RemoveBackup:     true,
			},
			expectedContains: []string{
				"clearsnapshot -t cleanup",
			},
			expectedArchive: "cleanup.tar.gz",
		},

		{
			name: "custom dump directory respected",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "analytics",
					Options: option.Options{
						Source: "nodetool",
					},
				},
				DumpName:         "analytics_dump",
				DumpDirRemote:    "/mnt/backups",
				DumpNameTemplate: "analytics_dump",
			},
			expectedContains: []string{
				"cd /mnt/backups",
			},
			expectedArchive: "analytics_dump.tar.gz",
		},
	}

	gen := cassandra.Generator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cmd, err := gen.Generate(tt.config)

			require.NoError(t, err)
			require.NotNil(t, cmd)

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command)
			assert.NotEmpty(t, cmd.DumpPath)

			for _, fragment := range tt.expectedContains {
				assert.Contains(t, cmd.Command, fragment)
			}

			assert.Equal(t, tt.expectedArchive, cmd.DumpPath)
		})
	}
}

func TestCassandraGenerator_CommandIntegrity(t *testing.T) {

	cfg := &cmdCfg.Config{
		Database: cmdCfg.Database{
			Name: "prod",
			Options: option.Options{
				Source: "/usr/bin/nodetool",
			},
		},
		DumpName:         "prod_backup",
		DumpDirRemote:    "/data",
		DumpNameTemplate: "prod_backup",
	}

	gen := cassandra.Generator{}

	cmd, err := gen.Generate(cfg)

	require.NoError(t, err)

	assert.Contains(
		t,
		cmd.Command,
		"/usr/bin/nodetool snapshot -t prod_backup prod",
	)

	assert.Equal(t, "prod_backup.tar.gz", cmd.DumpPath)
}
