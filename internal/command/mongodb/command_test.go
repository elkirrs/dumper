package mongodb

import (
	cmdConfig "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/setting"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestMongoGenerator_Generate(t *testing.T) {
	gen := MongoGenerator{}

	tests := []struct {
		name       string
		data       *cmdConfig.ConfigData
		settings   *setting.Settings
		wantCmd    string
		wantRemote string
	}{
		{
			name: "default port, dump format=bson, archive=false, dumpLocation=client",
			data: &cmdConfig.ConfigData{
				User:       "root",
				Password:   "pass",
				Port:       "",
				Name:       "testdb",
				DumpName:   "backup",
				DumpFormat: "dump",
				Options: option.Options{
					AuthSource: "",
					SSL:        boolPtr(false),
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    `/usr/bin/mongodump --uri "mongodb://root:pass@127.0.0.1:27017/testdb" --out ./`,
			wantRemote: "./backup.bson",
		},
		{
			name: "port specified, dump format=archive, archive=false, dumpLocation=server",
			data: &cmdConfig.ConfigData{
				User:       "admin",
				Password:   "secret",
				Port:       "27018",
				Name:       "mydb",
				DumpName:   "dumpfile",
				DumpFormat: "archive",
				Options: option.Options{
					AuthSource: "",
					SSL:        boolPtr(false),
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "server",
			},
			wantCmd:    `/usr/bin/mongodump --uri "mongodb://admin:secret@127.0.0.1:27018/mydb" --archive > ./dumpfile.archive`,
			wantRemote: "./dumpfile.archive",
		},
		{
			name: "dump format=archive, archive=true, dumpLocation=client",
			data: &cmdConfig.ConfigData{
				User:       "user1",
				Password:   "p123",
				Name:       "db1",
				DumpName:   "data",
				DumpFormat: "archive",
				Options: option.Options{
					AuthSource: "",
					SSL:        boolPtr(false),
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "client",
			},
			wantCmd:    `/usr/bin/mongodump --uri "mongodb://user1:p123@127.0.0.1:27017/db1" --archive --gzip`,
			wantRemote: "./data.archive.gz",
		},
		{
			name: "dump format=dump, archive=true, dumpLocation=server",
			data: &cmdConfig.ConfigData{
				User:       "alice",
				Password:   "pw",
				Port:       "27019",
				Name:       "sales",
				DumpName:   "sales_dump",
				DumpFormat: "dump",
				Options: option.Options{
					AuthSource: "admin",
					SSL:        boolPtr(true),
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "server",
			},
			wantCmd:    `/usr/bin/mongodump --uri "mongodb://alice:pw@127.0.0.1:27019/sales?authSource=admin&ssl=true" --out ./ && tar -czf sales_dump.tar.gz sales > ./sales_dump.tar.gz`,
			wantRemote: "./sales_dump.tar.gz",
		},
		{
			name: "dump format=dump, archive=false, ssl and authSource set, dumpLocation=client",
			data: &cmdConfig.ConfigData{
				User:       "bob",
				Password:   "pw2",
				Port:       "27020",
				Name:       "app",
				DumpName:   "app_backup",
				DumpFormat: "dump",
				Options: option.Options{
					AuthSource: "admin",
					SSL:        boolPtr(true),
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    `/usr/bin/mongodump --uri "mongodb://bob:pw2@127.0.0.1:27020/app?authSource=admin&ssl=true" --out ./`,
			wantRemote: "./app_backup.bson",
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
