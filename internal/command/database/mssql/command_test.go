package mssql_test

import (
	"dumper/internal/command/database/mssql"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMSQLGenerator_Generate_AllScenarios(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectErr        bool
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Generate .bak backup (no archive)",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "localhost"},
				Database: cmdCfg.Database{
					Format:   "bac",
					User:     "sa",
					Password: "Passw0rd!",
					Name:     "MyDB",
					Options:  option.Options{SSL: &falseVal},
				},
				DumpName:     "backup1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlcmd",
				"BACKUP DATABASE [MyDB]",
				"TO DISK='backup1.bak'",
			},
			expectedExt: ".bak",
		},
		{
			name: "Generate .bacpac export without SSL",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "localhost"},
				Database: cmdCfg.Database{
					Format:   "bacpac",
					User:     "sa",
					Password: "Passw0rd!",
					Name:     "MyDB",
					Options:  option.Options{SSL: &falseVal},
				},
				DumpName:     "export1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlpackage",
				"/Action:Export",
				"/SourceServerName:localhost",
				"/TargetFile:export1.bacpac",
			},
			expectedExt: ".bacpac",
		},
		{
			name: "Generate .bacpac export with SSL enabled",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "mssql-server"},
				Database: cmdCfg.Database{
					Format:   "bacpac",
					User:     "admin",
					Password: "pwd123",
					Name:     "ProdDB",
					Options:  option.Options{SSL: &trueVal},
				},
				DumpName:     "exportSSL",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlpackage",
				"/SourceTrustServerCertificate:True",
			},
			expectedExt: ".bacpac",
		},
		{
			name: "Generate .bak and compress archive",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "localhost"},
				Database: cmdCfg.Database{
					Format:   "bac",
					User:     "sa",
					Password: "Passw0rd!",
					Name:     "MyDB",
					Options:  option.Options{SSL: &falseVal},
				},
				DumpName:     "backup2",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"powershell Compress-Archive",
				"backup2.bak",
				"backup2.bak.gz",
			},
			expectedExt: ".bak.gz",
		},
		{
			name: "Generate .bacpac and compress archive",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "localhost"},
				Database: cmdCfg.Database{
					Format:   "bacpac",
					User:     "sa",
					Password: "Passw0rd!",
					Name:     "MyDB",
					Options:  option.Options{SSL: &falseVal},
				},
				DumpName:     "export2",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"powershell Compress-Archive",
				"export2.bacpac.gz",
			},
			expectedExt: ".bacpac.gz",
		},
		{
			name: "Server dump (no compression)",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "server123"},
				Database: cmdCfg.Database{
					Format:   "bac",
					User:     "sa",
					Password: "1234",
					Name:     "DB1",
					Options:  option.Options{SSL: &falseVal},
				},
				DumpName:     "dumpServer",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"sqlcmd",
				"TO DISK='dumpServer.bak'",
			},
			expectedExt: ".bak",
		},
		{
			name: "Unsupported database format",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "localhost"},
				Database: cmdCfg.Database{
					Format:   "unknown",
					User:     "sa",
					Password: "pass",
					Name:     "DBX",
					Options:  option.Options{SSL: &falseVal},
				},
				DumpName:     "badDump",
				Archive:      false,
				DumpLocation: "local",
			},
			expectErr: true,
		},
	}

	gen := mssql.MSQLGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)

			if tt.expectErr {
				require.Error(t, err, "Expected an error for unsupported format")
				assert.Nil(t, cmd)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cmd)

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command, "Command string should not be empty")
			assert.NotEmpty(t, cmd.DumpPath, "DumpPath should not be empty")

			for _, expect := range tt.expectedContains {
				assert.Contains(t, cmd.Command, expect, "Expected command to contain: %s", expect)
			}

			assert.Contains(t, cmd.DumpPath, tt.expectedExt, "DumpPath should end with expected extension")
		})
	}
}
