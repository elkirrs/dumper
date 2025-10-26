package postgres

import (
	"dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestPSQLGenerator_Generate(t *testing.T) {
	gen := PSQLGenerator{}

	tests := []struct {
		name       string
		data       *command_config.ConfigData
		settings   *setting.Settings
		wantCmd    string
		wantRemote string
	}{
		{
			name: "format=plain, archive=false, dumpLocation=client",
			data: &command_config.ConfigData{
				User:       "root",
				Password:   "pass",
				Port:       "5432",
				Name:       "mydb",
				DumpName:   "backup",
				DumpFormat: "plain",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    "/usr/bin/pg_dump --dbname=postgresql://root:pass@127.0.0.1:5432/mydb --clean --if-exists --no-owner -Fp",
			wantRemote: "./backup.sql",
		},
		{
			name: "format=plain, archive=true, dumpLocation=server",
			data: &command_config.ConfigData{
				User:       "admin",
				Password:   "secret",
				Port:       "5433",
				Name:       "testdb",
				DumpName:   "dumpfile",
				DumpFormat: "plain",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "server",
			},
			wantCmd:    "/usr/bin/pg_dump --dbname=postgresql://admin:secret@127.0.0.1:5433/testdb --clean --if-exists --no-owner -Fp | gzip > ./dumpfile.sql.gz",
			wantRemote: "./dumpfile.sql.gz",
		},
		{
			name: "format=dump, archive=true, dumpLocation=client",
			data: &command_config.ConfigData{
				User:       "user1",
				Password:   "p123",
				Port:       "5432",
				Name:       "db1",
				DumpName:   "data",
				DumpFormat: "dump",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "client",
			},
			wantCmd:    "/usr/bin/pg_dump --dbname=postgresql://user1:p123@127.0.0.1:5432/db1 --clean --if-exists --no-owner -Fc",
			wantRemote: "./data.dump",
		},
		{
			name: "format=tar, archive=false, dumpLocation=server",
			data: &command_config.ConfigData{
				User:       "alice",
				Password:   "pw",
				Port:       "5434",
				Name:       "sales",
				DumpName:   "sales_dump",
				DumpFormat: "tar",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "server",
			},
			wantCmd:    "/usr/bin/pg_dump --dbname=postgresql://alice:pw@127.0.0.1:5434/sales --clean --if-exists --no-owner -Ft > ./sales_dump.tar",
			wantRemote: "./sales_dump.tar",
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
