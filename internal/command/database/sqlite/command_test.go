package sqlite

import (
	"dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/setting"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestSQLiteGenerator_Generate(t *testing.T) {
	gen := SQLiteGenerator{}

	tests := []struct {
		name       string
		data       *command_config.ConfigData
		settings   *setting.Settings
		wantCmd    string
		wantRemote string
	}{
		{
			name: "archive=false, dumpLocation=client",
			data: &command_config.ConfigData{
				DumpName: "backup",
				Options: option.Options{
					Path: "/tmp/db.sqlite",
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    "sqlite3 /tmp/db.sqlite .dump > ./backup.sql",
			wantRemote: "./backup.sql",
		},
		{
			name: "archive=false, dumpLocation=server",
			data: &command_config.ConfigData{
				DumpName: "dumpfile",
				Options: option.Options{
					Path: "/var/db/app.db",
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "server",
			},
			wantCmd:    "sqlite3 /var/db/app.db .dump > ./dumpfile.sql",
			wantRemote: "./dumpfile.sql",
		},
		{
			name: "archive=true, dumpLocation=client",
			data: &command_config.ConfigData{
				DumpName: "data",
				Options: option.Options{
					Path: "/data/db.sqlite3",
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "client",
			},
			wantCmd:    "sqlite3 /data/db.sqlite3 .dump | gzip > ./data.sql.gz",
			wantRemote: "./data.sql.gz",
		},
		{
			name: "archive=true, dumpLocation=server",
			data: &command_config.ConfigData{
				DumpName: "app_backup",
				Options: option.Options{
					Path: "/opt/app/app.db",
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "server",
			},
			wantCmd:    "sqlite3 /opt/app/app.db .dump | gzip > ./app_backup.sql.gz",
			wantRemote: "./app_backup.sql.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.NotNil(t, tt.data)
			require.NotNil(t, tt.settings)

			gotCmd, gotRemote := gen.Generate(tt.data, tt.settings)

			assert.Equal(t, tt.wantCmd, gotCmd, "Command mismatch in test '%s'", tt.name)
			assert.Equal(t, tt.wantRemote, gotRemote, "Remote path mismatch in test '%s'", tt.name)
		})
	}
}
