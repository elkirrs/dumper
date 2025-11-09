package mongodb_test

import (
	"dumper/internal/command/database/mongodb"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongoGenerator_Generate_AllScenarios(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Default BSON dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "mongoUser",
					Password: "secret",
					Port:     "27017",
					Name:     "testdb",
					Format:   "dump",
					Options: option.Options{
						SSL:        &falseVal,
						AuthSource: "",
					},
				},
				DumpName:     "dump1",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"mongodump",
				"--uri",
				"mongodb://mongoUser:secret@127.0.0.1:27017/testdb",
				"--out ./",
			},
			expectedExt: "bson",
		},
		{
			name: "Archive format, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "mongoUser",
					Password: "secret",
					Port:     "27017",
					Name:     "testdb",
					Format:   "archive",
					Options: option.Options{
						SSL:        &falseVal,
						AuthSource: "",
					},
				},
				DumpName:     "dump2",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"--archive",
			},
			expectedExt: "archive",
		},
		{
			name: "BSON dump with gzip (Archive=true)",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "mongoUser",
					Password: "secret",
					Port:     "27017",
					Name:     "testdb",
					Format:   "dump",
					Options: option.Options{
						SSL:        &falseVal,
						AuthSource: "",
					},
				},
				DumpName:     "dump3",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"tar -czf dump3.tar.gz testdb",
			},
			expectedExt: "tar.gz",
		},
		{
			name: "Archive format + gzip",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "mongoUser",
					Password: "secret",
					Port:     "27017",
					Name:     "testdb",
					Format:   "archive",
					Options: option.Options{
						SSL:        &falseVal,
						AuthSource: "",
					},
				},
				DumpName:     "dump4",
				Archive:      true,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"--archive",
				"--gzip",
			},
			expectedExt: "archive.gz",
		},
		{
			name: "Server dump output redirect",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "mongoUser",
					Password: "secret",
					Port:     "27017",
					Name:     "testdb",
					Format:   "dump",
					Options: option.Options{
						SSL:        &falseVal,
						AuthSource: "",
					},
				},
				DumpName:     "serverDump",
				Archive:      false,
				DumpLocation: "server",
			},
			expectedContains: []string{
				"> ./serverDump.bson",
			},
			expectedExt: "bson",
		},
		{
			name: "SSL and AuthSource enabled",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "mongoUser",
					Password: "secret",
					Port:     "27017",
					Name:     "testdb",
					Format:   "dump",
					Options: option.Options{
						SSL:        &trueVal,
						AuthSource: "admin",
					},
				},
				DumpName:     "secureDump",
				Archive:      false,
				DumpLocation: "local",
			},
			expectedContains: []string{
				"authSource=admin",
				"ssl=true",
			},
			expectedExt: "bson",
		},
	}

	gen := mongodb.MongoGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)
			require.NoError(t, err, "Generate() should not return an error")
			require.NotNil(t, cmd, "Returned command must not be nil")

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command, "Command string should not be empty")
			assert.NotEmpty(t, cmd.DumpPath, "DumpPath should not be empty")

			for _, expect := range tt.expectedContains {
				assert.Contains(t, cmd.Command, expect, "Command should contain expected fragment: %s", expect)
			}

			assert.Contains(t, cmd.DumpPath, tt.expectedExt, "DumpPath should have the correct extension")
		})
	}
}
