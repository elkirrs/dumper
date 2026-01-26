package firebird_test

import (
	"dumper/internal/command/database/firebird"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirebirdDbGenerator_Generate_Variations(t *testing.T) {

	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedArchive  string
	}{
		{
			name: "Minimal backup",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source: "gbak",
						Path:   "/var/lib/firebird/data.fdb",
					},
				},
				DumpName:         "backup1",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup1",
			},
			expectedContains: []string{
				"gbak -b",
				"localhost/3050:/var/lib/firebird/data.fdb",
				"backup1.fbk -user sysdba -password masterkey",
			},
			expectedArchive: "backup1.fbk",
		},
		{
			name: "Skip garbage and skip issue",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source:      "gbak",
						Path:        "/data/db.fdb",
						SkipGarbage: true,
						SkipIssue:   true,
					},
				},
				DumpName:         "backup2",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup2",
			},
			expectedContains: []string{
				"gbak -b -g -ignore",
			},
			expectedArchive: "backup2.fbk",
		},
		{
			name: "FastAndStable mode",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source:        "gbak",
						Path:          "/fast/db.fdb",
						FastAndStable: true,
					},
				},
				DumpName:         "backup3",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup3",
			},
			expectedContains: []string{
				"gbak -b -se service_mgr /fast/db.fdb",
			},
			expectedArchive: "backup3.fbk",
		},
		{
			name: "Archive enabled with remove backup",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source: "gbak",
						Path:   "/data/db.fdb",
					},
				},
				DumpName:         "backup4",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup4",
				Archive:          true,
				RemoveBackup:     true,
			},
			expectedContains: []string{
				"tar -czf backup4.tar.gz -C /backups backup4.fbk",
				"rm -rf backup4.fbk",
			},
			expectedArchive: "backup4.tar.gz",
		},
		{
			name: "Dump to server location",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source: "gbak",
						Path:   "/data/db.fdb",
					},
				},
				DumpName:         "server_backup",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "server_backup",
				DumpLocation:     "server",
			},
			expectedContains: []string{
				"server_backup.fbk -user sysdba -password masterkey",
			},
			expectedArchive: "server_backup.fbk",
		},
	}

	gen := firebird.FirebirdDbGenerator{}

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

func TestFirebirdDbGenerator_EdgeCases(t *testing.T) {
	gen := firebird.FirebirdDbGenerator{}

	tests := []struct {
		name           string
		config         *cmdCfg.Config
		expectedErr    bool
		expectedOutput []string
	}{
		{
			name: "SkipGarbage and FastAndStable together",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source:        "gbak",
						Path:          "/data/db.fdb",
						SkipGarbage:   true,
						FastAndStable: true,
					},
				},
				DumpName:         "backup_edge1",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup_edge1",
			},
			expectedErr: false,
			expectedOutput: []string{
				"-b -g -se service_mgr /data/db.fdb",
				"backup_edge1.fbk -user sysdba -password masterkey",
			},
		},
		{
			name: "Archive without RemoveBackup",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source: "gbak",
						Path:   "/data/db.fdb",
					},
				},
				DumpName:         "backup_edge2",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup_edge2",
				Archive:          true,
				RemoveBackup:     false,
			},
			expectedErr: false,
			expectedOutput: []string{
				"tar -czf backup_edge2.tar.gz -C /backups backup_edge2.fbk",
			},
		},
		{
			name: "Empty Source should still generate command",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source: "",
						Path:   "/data/db.fdb",
					},
				},
				DumpName:         "backup_edge3",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup_edge3",
			},
			expectedErr: false,
			expectedOutput: []string{
				" -b localhost/3050:/data/db.fdb backup_edge3.fbk -user sysdba -password masterkey",
			},
		},
		{
			name: "Empty Path should still generate command",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "sysdba",
					Password: "masterkey",
					Port:     "3050",
					Options: option.Options{
						Source: "gbak",
						Path:   "",
					},
				},
				DumpName:         "backup_edge4",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup_edge4",
			},
			expectedErr: false,
			expectedOutput: []string{
				"gbak -b localhost/3050: backup_edge4.fbk -user sysdba -password masterkey",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)
			if tt.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cmd)

			for _, fragment := range tt.expectedOutput {
				assert.Contains(t, cmd.Command, fragment)
			}
		})
	}
}
