package influxdb_test

import (
	"dumper/internal/command/database/influxdb"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfluxDB2Generator_Generate_Variations(t *testing.T) {
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedArchive  string
	}{
		{
			name: "2.x minimal config",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:    "2.x",
						Source:     "influx",
						Host:       "localhost",
						SkipVerify: boolPtr(false),
					},
				},
				DumpName:         "dump",
				DumpDirRemote:    "/tmp",
				DumpNameTemplate: "dump",
			},
			expectedContains: []string{
				"influx backup",
				"--host localhost:8086",
				"--token token",
				"dump",
				"tar -czf dump.tar.gz",
			},
			expectedArchive: "dump.tar.gz",
		},
		{
			name: "2.x with bucket",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:    "2.x",
						Source:     "influx",
						Host:       "localhost",
						Bucket:     "metrics",
						SkipVerify: boolPtr(false),
					},
				},
				DumpName:         "bucket_dump",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "bucket_dump",
			},
			expectedContains: []string{
				"--bucket metrics",
			},
			expectedArchive: "bucket_dump.tar.gz",
		},
		{
			name: "2.x with bucket-id",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:    "2.x",
						Source:     "influx",
						Host:       "localhost",
						BucketId:   "bucket-id-123",
						SkipVerify: boolPtr(false),
					},
				},
				DumpName:         "bucket_id_dump",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "bucket_id_dump",
			},
			expectedContains: []string{
				"--bucket-id bucket-id-123",
			},
			expectedArchive: "bucket_id_dump.tar.gz",
		},
		{
			name: "2.x with org",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:      "2.x",
						Source:       "influx",
						Host:         "localhost",
						Organization: "my-org",
						SkipVerify:   boolPtr(false),
					},
				},
				DumpName:         "org_dump",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "org_dump",
			},
			expectedContains: []string{
				"--org my-org",
			},
			expectedArchive: "org_dump.tar.gz",
		},
		{
			name: "2.x with org-id",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:        "2.x",
						Source:         "influx",
						Host:           "localhost",
						OrganizationId: "org-id-999",
						SkipVerify:     boolPtr(false),
					},
				},
				DumpName:         "org_id_dump",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "org_id_dump",
			},
			expectedContains: []string{
				"--org-id org-id-999",
			},
			expectedArchive: "org_id_dump.tar.gz",
		},
		{
			name: "2.x with start and end time",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:    "2.x",
						Source:     "influx",
						Host:       "localhost",
						Start:      "2023-01-01T00:00:00Z",
						End:        "2023-01-31T23:59:59Z",
						SkipVerify: boolPtr(false),
					},
				},
				DumpName:         "timerange",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "timerange",
			},
			expectedContains: []string{
				"--start 2023-01-01T00:00:00Z",
				"--end 2023-01-31T23:59:59Z",
			},
			expectedArchive: "timerange.tar.gz",
		},
		{
			name: "2.x with filter",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:    "2.x",
						Source:     "influx",
						Host:       "localhost",
						Filter:     `_measurement="cpu"`,
						SkipVerify: boolPtr(false),
					},
				},
				DumpName:         "filter_dump",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "filter_dump",
			},
			expectedContains: []string{
				`--filter _measurement="cpu"`,
			},
			expectedArchive: "filter_dump.tar.gz",
		},
		{
			name: "2.x with skip-verify and remove backup",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port:  "8086",
					Token: "token",
					Options: option.Options{
						Version:    "2.x",
						Source:     "influx",
						Host:       "localhost",
						SkipVerify: boolPtr(true),
					},
				},
				DumpName:         "cleanup",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "cleanup",
				RemoveBackup:     true,
			},
			expectedContains: []string{
				"--skip-verify",
				"rm -rf cleanup",
			},
			expectedArchive: "cleanup.tar.gz",
		},
		{
			name: "3.x full filesystem backup",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Options: option.Options{
						Version: "3.x",
						DataDir: "/var/lib/influxdb3",
						NodeId:  "node-1",
					},
				},
				DumpName:         "backup3x",
				DumpDirRemote:    "/backups",
				DumpNameTemplate: "backup3x",
			},
			expectedContains: []string{
				"cp -r",
				"snapshots",
				"dbs",
				"wal",
				"catalog",
				"_catalog_checkpoint",
				"tar -czf backup3x.tar.gz",
			},
			expectedArchive: "backup3x.tar.gz",
		},
	}

	gen := influxdb.Generator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := gen.Generate(tt.config)

			require.NoError(t, err)
			require.NotNil(t, cmd)

			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command)
			assert.NotEmpty(t, cmd.DumpPath)

			for _, fragment := range tt.expectedContains {
				assert.Contains(t, cmd.Command, fragment)
			}

			assert.Equal(t, tt.expectedArchive, cmd.DumpPath)
		})
	}
}

func TestInfluxDB2Generator_CommandIntegrity_2x(t *testing.T) {
	cfg := &cmdCfg.Config{
		Database: cmdCfg.Database{
			Port:  "8086",
			Token: "secure",
			Options: option.Options{
				Version:    "2.x",
				Source:     "influx",
				Host:       "localhost",
				SkipVerify: boolPtr(false),
			},
		},
		DumpName:         "integrity",
		DumpDirRemote:    "/tmp",
		DumpNameTemplate: "integrity",
	}

	gen := influxdb.Generator{}
	cmd, err := gen.Generate(cfg)

	require.NoError(t, err)

	assert.Contains(
		t,
		cmd.Command,
		"influx backup --host localhost:8086 --token secure",
	)
	assert.Equal(t, "integrity.tar.gz", cmd.DumpPath)
}

func boolPtr(v bool) *bool {
	return &v
}
