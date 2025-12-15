package redis_test

import (
	"dumper/internal/command/database/redis"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisGenerator_Generate_AllScenarios(t *testing.T) {
	source := utils.GetDBSource("redis", "")
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Standard RDB dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "6379",
					Password: "pass",
					Options:  option.Options{Mode: "", Source: source},
				},
				DumpName:     "dump1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"redis-cli",
				"--rdb",
				"dump1.rdb",
			},
			expectedExt: "rdb",
		},
		{
			name: "RDB dump with SAVE mode, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "6379",
					Password: "pass",
					Options:  option.Options{Mode: "save", Source: source},
				},
				DumpName:     "dump2",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"SAVE",
				"--rdb",
				"dump2.rdb",
			},
			expectedExt: "rdb",
		},
		{
			name: "RDB dump with archive, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "6380",
					Password: "secret",
					Options:  option.Options{Mode: "", Source: source},
				},
				DumpName:     "archive1",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"--rdb",
				"| gzip",
				"archive1.rdb.gz",
			},
			expectedExt: "rdb.gz",
		},
		{
			name: "RDB dump with SAVE mode and archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "6380",
					Password: "secret",
					Options:  option.Options{Mode: "save", Source: source},
				},
				DumpName:     "archive2",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"SAVE",
				"--rdb",
				"| gzip",
				"archive2.rdb.gz",
			},
			expectedExt: "rdb.gz",
		},
		{
			name: "RDB dump to server, no archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "6379",
					Password: "pass",
					Options:  option.Options{Mode: "", Source: source},
				},
				DumpName:     "serverDump",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"--rdb",
				"serverDump.rdb",
			},
			expectedExt: "rdb",
		},
		{
			name: "RDB dump to server with archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:     "6379",
					Password: "pass",
					Options:  option.Options{Mode: "save", Source: source},
				},
				DumpName:     "serverArchive",
				Archive:      true,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"SAVE",
				"--rdb",
				"| gzip",
				"serverArchive.rdb.gz",
			},
			expectedExt: "rdb.gz",
		},
	}

	gen := redis.RedisGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)
			require.NoError(t, err, "Generate() should not return an error")
			require.NotNil(t, cmd, "DBCommand should not be nil")

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command, "Command string should not be empty")
			assert.NotEmpty(t, cmd.DumpPath, "DumpPath should not be empty")

			for _, expected := range tt.expectedContains {
				assert.Contains(t, cmd.Command, expected, "Command should contain %s", expected)
			}

			assert.Contains(t, cmd.DumpPath, tt.expectedExt, "DumpPath should contain extension %s", tt.expectedExt)
		})
	}
}
