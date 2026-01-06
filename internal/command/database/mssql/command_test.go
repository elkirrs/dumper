package mssql_test

import (
	"dumper/internal/command/database/mssql"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils/mapping"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMSQLGenerator_Generate_AllScenarios(t *testing.T) {
	trueVal := true
	falseVal := false

	sourceBac := mapping.GetDBSource("mssql", "bac")
	sourceBacpac := mapping.GetDBSource("mssql", "bacpac")
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
					Options:  option.Options{SSL: &falseVal, Source: sourceBac},
				},
				DumpName:     "backup1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlcmd",
				"BACKUP DATABASE [MyDB]",
				"TO DISK=\\\"backup1.bak\\\"",
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
					Options:  option.Options{SSL: &falseVal, Source: sourceBacpac},
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
					Options:  option.Options{SSL: &trueVal, Source: sourceBacpac},
				},
				DumpName:     "exportSSL",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlpackage",
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
					Options:  option.Options{SSL: &falseVal, Source: sourceBac},
				},
				DumpName:     "backup2",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"tar -czf",
				"backup2.bak.gz",
				"backup2.bak",
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
					Options:  option.Options{SSL: &falseVal, Source: sourceBacpac},
				},
				DumpName:     "export2",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"tar -czf",
				"export2.bacpac.gz",
				"export2.bacpac",
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
					Options:  option.Options{SSL: &falseVal, Source: sourceBac},
				},
				DumpName:     "dumpServer",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"sqlcmd",
				"TO DISK=\\\"dumpServer.bak\\\"",
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
					Options:  option.Options{SSL: &falseVal, Source: sourceBac},
				},
				DumpName:     "badDump",
				Archive:      false,
				DumpLocation: "local",
			},
			expectErr: true,
		},
		{
			name: "Generate .bacpac with include tables",
			config: &cmdCfg.Config{
				Server: cmdCfg.Server{Host: "localhost"},
				Database: cmdCfg.Database{
					Format:   "bacpac",
					User:     "sa",
					Password: "Passw0rd!",
					Name:     "MyDB",
					Options: option.Options{
						SSL:       &falseVal,
						Source:    sourceBacpac,
						IncTables: []string{"Users", "Rules"},
					},
				},
				DumpName:     "export1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"sqlpackage",
				"/Action:Export",
				"/SourceServerName:localhost",
				"/p:TableData=\"dbo.Users\"",
				"/p:TableData=\"dbo.Rules\"",
				"/TargetFile:export1.bacpac",
			},
			expectedExt: ".bacpac",
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
