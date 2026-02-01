package db2_test

import (
	"dumper/internal/command/database/db2"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB2Generator_Generate_Variations(t *testing.T) {
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedArchive  string
	}{
		{
			name: "offline backup minimal",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "appdb",
					Options: option.Options{
						Source:     "db2",
						BackupMode: "offline",
					},
				},
				DumpName:         "appdb_dump",
				DumpDirRemote:    "/database/dumps",
				DumpNameTemplate: "appdb_dump",
			},
			expectedContains: []string{
				"mkdir -p appdb_dump",
				"db2 backup database appdb to appdb_dump",
				"tar -czf appdb_dump.tar.gz",
			},
			expectedArchive: "appdb_dump.tar.gz",
		},
		{
			name: "online backup",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "appdb",
					Options: option.Options{
						Source:     "db2",
						BackupMode: "online",
					},
				},
				DumpName:         "online_dump",
				DumpDirRemote:    "/database/dumps",
				DumpNameTemplate: "online_dump",
			},
			expectedContains: []string{
				"db2 backup database appdb online to online_dump",
			},
			expectedArchive: "online_dump.tar.gz",
		},
		{
			name: "remove backup after archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "appdb",
					Options: option.Options{
						Source:     "db2",
						BackupMode: "online",
					},
				},
				DumpName:         "cleanup_dump",
				DumpDirRemote:    "/database/dumps",
				DumpNameTemplate: "cleanup_dump",
				RemoveBackup:     true,
			},
			expectedContains: []string{
				"rm -rf cleanup_dump",
			},
			expectedArchive: "cleanup_dump.tar.gz",
		},
		{
			name: "custom db name and source",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name: "TESTDB",
					Options: option.Options{
						Source:     "/opt/ibm/db2/bin/db2",
						BackupMode: "online",
					},
				},
				DumpName:         "testdb_backup",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "testdb_backup",
			},
			expectedContains: []string{
				"/opt/ibm/db2/bin/db2 backup database TESTDB online",
			},
			expectedArchive: "testdb_backup.tar.gz",
		},
	}

	gen := db2.Generator{}

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

func TestDB2Generator_CommandIntegrity(t *testing.T) {
	cfg := &cmdCfg.Config{
		Database: cmdCfg.Database{
			Name: "appdb",
			Options: option.Options{
				Source:     "db2",
				BackupMode: "online",
			},
		},
		DumpName:         "integrity",
		DumpDirRemote:    "/database/dumps",
		DumpNameTemplate: "integrity",
	}

	gen := db2.Generator{}
	cmd, err := gen.Generate(cfg)

	require.NoError(t, err)

	assert.Contains(
		t,
		cmd.Command,
		"db2 backup database appdb online to integrity",
	)
	assert.Contains(
		t,
		cmd.Command,
		"tar -czf integrity.tar.gz",
	)

	assert.Equal(t, "integrity.tar.gz", cmd.DumpPath)
}
