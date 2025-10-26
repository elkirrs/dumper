package mssql

import (
	"dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/setting"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestMSSQLGenerator_Generate(t *testing.T) {
	gen := MSSQLGenerator{}

	tests := []struct {
		name       string
		data       *command_config.ConfigData
		settings   *setting.Settings
		wantCmd    string
		wantRemote string
	}{
		{
			name: "bac, archive=false, dumpLocation=client",
			data: &command_config.ConfigData{
				Name:       "testdb",
				User:       "sa",
				Password:   "pass",
				DumpName:   "backup",
				DumpFormat: "bac",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    "sqlcmd -S testdb -U sa -P pass -Q \"BACKUP DATABASE [testdb] TO DISK='./backup.bak' WITH FORMAT, INIT\"",
			wantRemote: "./backup.bak",
		},
		{
			name: "bac, archive=true, dumpLocation=client",
			data: &command_config.ConfigData{
				Name:       "testdb",
				User:       "sa",
				Password:   "pass",
				DumpName:   "backup",
				DumpFormat: "bac",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "client",
			},
			wantCmd:    "sqlcmd -S testdb -U sa -P pass -Q \"BACKUP DATABASE [testdb] TO DISK='./backup.bak' WITH FORMAT, INIT\" && powershell Compress-Archive -Path './backup.bak' -DestinationPath './backup.bak.gz'",
			wantRemote: "./backup.bak.gz",
		},
		{
			name: "bacpac, SSL=false, archive=false",
			data: &command_config.ConfigData{
				Host:       "localhost",
				Name:       "db1",
				User:       "admin",
				Password:   "pw",
				DumpName:   "db1dump",
				DumpFormat: "bacpac",
				Options:    option.Options{SSL: new(bool)}, // false
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    "sqlpackage /Action:Export /SourceServerName:localhost /SourceDatabaseName:db1 /SourceUser:admin /SourcePassword:pw /TargetFile:./db1dump.bacpac",
			wantRemote: "./db1dump.bacpac",
		},
		{
			name: "bacpac, SSL=true, archive=true",
			data: &command_config.ConfigData{
				Host:       "localhost",
				Name:       "db2",
				User:       "sa",
				Password:   "secret",
				DumpName:   "db2dump",
				DumpFormat: "bacpac",
				Options:    option.Options{SSL: boolPtr(true)},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "server",
			},
			wantCmd:    "sqlpackage /Action:Export /SourceServerName:localhost /SourceDatabaseName:db2 /SourceUser:sa /SourcePassword:secret /TargetFile:./db2dump.bacpac /SourceTrustServerCertificate:True && powershell Compress-Archive -Path './db2dump.bacpac' -DestinationPath './db2dump.bacpac.gz'",
			wantRemote: "./db2dump.bacpac.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.data)
			require.NotNil(t, tt.settings)

			gotCmd, gotRemote := gen.Generate(tt.data, tt.settings)

			assert.Equal(t, tt.wantCmd, gotCmd, "Command mismatch in test '%s'", tt.name)
			assert.Equal(t, tt.wantRemote, gotRemote, "Remote path mismatch in test '%s'", tt.name)
		})
	}
}
