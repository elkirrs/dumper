package mysql

import (
	"dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestMySQLGenerator_Generate(t *testing.T) {
	gen := MySQLGenerator{}

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
				User:     "root",
				Password: "pass",
				Port:     "3306",
				Name:     "mydb",
				DumpName: "backup",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    "/usr/bin/mysqldump -h 127.0.0.1 -P 3306 -u root -ppass mydb",
			wantRemote: "./backup.sql",
		},
		{
			name: "archive=false, dumpLocation=server",
			data: &command_config.ConfigData{
				User:     "admin",
				Password: "secret",
				Port:     "3307",
				Name:     "testdb",
				DumpName: "dumpfile",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "server",
			},
			wantCmd:    "/usr/bin/mysqldump -h 127.0.0.1 -P 3307 -u admin -psecret testdb > ./dumpfile.sql",
			wantRemote: "./dumpfile.sql",
		},
		{
			name: "archive=true, dumpLocation=client",
			data: &command_config.ConfigData{
				User:     "user1",
				Password: "p123",
				Port:     "3306",
				Name:     "db1",
				DumpName: "data",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "client",
			},
			wantCmd:    "/usr/bin/mysqldump -h 127.0.0.1 -P 3306 -u user1 -pp123 db1 | gzip",
			wantRemote: "./data.sql.gz",
		},
		{
			name: "archive=true, dumpLocation=server",
			data: &command_config.ConfigData{
				User:     "alice",
				Password: "pw",
				Port:     "3308",
				Name:     "sales",
				DumpName: "sales_dump",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "server",
			},
			wantCmd:    "/usr/bin/mysqldump -h 127.0.0.1 -P 3308 -u alice -ppw sales | gzip > ./sales_dump.sql.gz",
			wantRemote: "./sales_dump.sql.gz",
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
