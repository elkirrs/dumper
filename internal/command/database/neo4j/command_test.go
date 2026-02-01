package neo4j_test

import (
	"dumper/internal/command/database/neo4j"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils/mapping"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeo4jGenerator_Generate_AllScenarios(t *testing.T) {
	source := mapping.GetDBSource("neo4j", "")
	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedExt      string
	}{
		{
			name: "Standard Neo4j dump, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name:    "mydb",
					Options: option.Options{Source: source},
				},
				DumpName:      "dump1",
				Archive:       false,
				DumpLocation:  "local",
				DumpDirRemote: "/var/lib/neo4j/dumps",
			},
			expectedContains: []string{
				"neo4j-admin database dump mydb",
				"--to-path=/var/lib/neo4j/dumps",
				"mv /var/lib/neo4j/dumps/mydb.dump dump1.dump",
			},
			expectedExt: "dump",
		},
		{
			name: "Neo4j dump with archive, local",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name:    "testdb",
					Options: option.Options{Source: source},
				},
				DumpName:      "archive1",
				Archive:       true,
				DumpLocation:  "local",
				DumpDirRemote: "/var/lib/neo4j/dumps",
			},
			expectedContains: []string{
				"neo4j-admin database dump testdb",
				"--to-path=/var/lib/neo4j/dumps",
				"gzip -c",
				"archive1.dump.gz",
			},
			expectedExt: "dump.gz",
		},
		{
			name: "Neo4j dump to server, no archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name:    "proddb",
					Options: option.Options{Source: source},
				},
				DumpName:      "serverDump",
				Archive:       false,
				DumpLocation:  "server",
				DumpDirRemote: "/var/lib/neo4j/dumps",
			},
			expectedContains: []string{
				"neo4j-admin database dump proddb",
				"--to-path=/var/lib/neo4j/dumps",
				"mv /var/lib/neo4j/dumps/proddb.dump serverDump.dump",
			},
			expectedExt: "dump",
		},
		{
			name: "Neo4j dump to server with archive",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Name:    "graphdb",
					Options: option.Options{Source: source},
				},
				DumpName:      "serverArchive",
				Archive:       true,
				DumpLocation:  "server",
				DumpDirRemote: "/var/lib/neo4j/dumps",
			},
			expectedContains: []string{
				"neo4j-admin database dump graphdb",
				"--to-path=/var/lib/neo4j/dumps",
				"gzip -c",
				"serverArchive.dump.gz",
			},
			expectedExt: "dump.gz",
		},
	}

	gen := neo4j.Generator{}

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

func TestNeo4jGenerator_CommandIntegrity(t *testing.T) {
	source := mapping.GetDBSource("neo4j", "")

	cfg := &cmdCfg.Config{
		Database: cmdCfg.Database{
			Name:    "mydb",
			Options: option.Options{Source: source},
		},
		DumpName:     "backup",
		Archive:      false,
		DumpLocation: "local",
	}

	gen := neo4j.Generator{}
	cmd, err := gen.Generate(cfg)
	require.NoError(t, err)
	require.NotNil(t, cmd)

	expectedPrefix := "neo4j-admin database dump mydb"
	assert.Contains(t, cmd.Command, expectedPrefix, "Command must be constructed correctly")
	assert.Equal(t, "backup.dump", cmd.DumpPath, "DumpPath should match expected filename")
}
