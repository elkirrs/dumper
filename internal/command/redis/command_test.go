package redis

import (
	"dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/internal/domain/config/setting"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestRedisGenerator_Generate(t *testing.T) {
	gen := RedisGenerator{}

	tests := []struct {
		name       string
		data       *command_config.ConfigData
		settings   *setting.Settings
		wantCmd    string
		wantRemote string
	}{
		{
			name: "archive=false, mode=default, dumpLocation=client",
			data: &command_config.ConfigData{
				Port:     "6379",
				Password: "pass",
				DumpName: "backup",
				Options:  option.Options{},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "client",
			},
			wantCmd:    "redis-cli -h 127.0.0.1 -p 6379 -a pass --rdb ./backup.rdb",
			wantRemote: "./backup.rdb",
		},
		{
			name: "archive=false, mode=save, dumpLocation=server",
			data: &command_config.ConfigData{
				Port:     "6380",
				Password: "secret",
				DumpName: "dumpfile",
				Options: option.Options{
					Mode: "save",
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(false),
				DumpLocation: "server",
			},
			wantCmd:    "redis-cli -h 127.0.0.1 -p 6380 -a secret SAVE && redis-cli -h 127.0.0.1 -p 6380 -a secret --rdb ./dumpfile.rdb",
			wantRemote: "./dumpfile.rdb",
		},
		{
			name: "archive=true, mode=default, dumpLocation=client",
			data: &command_config.ConfigData{
				Port:     "6379",
				Password: "p123",
				DumpName: "data",
				Options:  option.Options{},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "client",
			},
			wantCmd:    "redis-cli -h 127.0.0.1 -p 6379 -a p123 --rdb - | gzip > ./data.rdb.gz",
			wantRemote: "./data.rdb.gz",
		},
		{
			name: "archive=true, mode=save, dumpLocation=server",
			data: &command_config.ConfigData{
				Port:     "6381",
				Password: "pw",
				DumpName: "redis_save",
				Options: option.Options{
					Mode: "save",
				},
			},
			settings: &setting.Settings{
				Archive:      boolPtr(true),
				DumpLocation: "server",
			},
			wantCmd:    "redis-cli -h 127.0.0.1 -p 6381 -a pw SAVE && redis-cli -h 127.0.0.1 -p 6381 -a pw --rdb - | gzip > ./redis_save.rdb.gz",
			wantRemote: "./redis_save.rdb.gz",
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
