package mariadb

import (
	"dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestMariaDBGenerator_Generate(t *testing.T) {
	gen := MariaDBGenerator{}

	tests := []struct {
		name       string
		data       *command_config.ConfigData
		settings   *setting.Settings
		wantCmd    string
		wantRemote string
	}{
		{
			name: "port is empty, archive=false, dumpLocation=client",
			data: &command_config.ConfigData{
				User:     "root",
				Password: "pass",
				Port:     "",
				Name:     "mydb",
				DumpName: "backup",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    "/usr/bin/mariadb-dump -uroot -ppass -h127.0.0.1 -P3306 mydb",
			wantRemote: "./backup.sql",
		},
		{
			name: "port is set, archive=false, dumpLocation=server",
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
			wantCmd:    "/usr/bin/mariadb-dump -uadmin -psecret -h127.0.0.1 -P3307 testdb > ./dumpfile.sql",
			wantRemote: "./dumpfile.sql",
		},
		{
			name: "port empty, archive=true, dumpLocation=client",
			data: &command_config.ConfigData{
				User:     "user1",
				Password: "p123",
				Name:     "db1",
				DumpName: "data",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "client",
			},
			wantCmd:    "/usr/bin/mariadb-dump -uuser1 -pp123 -h127.0.0.1 -P3306 db1 | gzip",
			wantRemote: "./data.sql.gz",
		},
		{
			name: "port set, archive=true, dumpLocation=server",
			data: &command_config.ConfigData{
				User:     "alice",
				Password: "pw",
				Port:     "3310",
				Name:     "sales",
				DumpName: "sales_dump",
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "server",
			},
			wantCmd:    "/usr/bin/mariadb-dump -ualice -ppw -h127.0.0.1 -P3310 sales | gzip > ./sales_dump.sql.gz",
			wantRemote: "./sales_dump.sql.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.data)
			require.NotNil(t, tt.settings)

			gotCmd, gotRemote := gen.Generate(tt.data, tt.settings)

			assert.Equal(t, tt.wantCmd, gotCmd,
				"Command mismatch in test '%s'", tt.name)
			assert.Equal(t, tt.wantRemote, gotRemote,
				"Remote path mismatch in test '%s'", tt.name)
		})
	}
}
