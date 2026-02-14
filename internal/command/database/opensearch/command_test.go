package opensearch_test

import (
	"dumper/internal/command/database/opensearch"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenSearchGenerator_Generate_Variations(t *testing.T) {

	tests := []struct {
		name             string
		config           *cmdCfg.Config
		expectedContains []string
		expectedArchive  string
	}{
		{
			name: "minimal config",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port: "9200",
					Options: option.Options{
						Source: "curl",
						Host:   "http://localhost",
					},
				},
				DumpName:         "dump",
				DumpDirRemote:    "/backup",
				DumpNameTemplate: "repo",
			},
			expectedContains: []string{
				"curl -f -X GET http://localhost:9200/_snapshot/repo",
				"_snapshot/repo/",
				"tar -czf dump.tar.gz",
				"rm -Rf dump",
			},
			expectedArchive: "dump.tar.gz",
		},

		{
			name: "basic auth enabled",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					User:     "admin",
					Password: "secret",
					Port:     "9200",
					Options: option.Options{
						Source: "curl",
						Host:   "https://node",
					},
				},
				DumpName:         "secure",
				DumpDirRemote:    "/data",
				DumpNameTemplate: "repo_secure",
			},
			expectedContains: []string{
				"-u admin:secret",
			},
			expectedArchive: "secure.tar.gz",
		},

		{
			name: "token auth enabled",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Token: "Bearer abc123",
					Port:  "9200",
					Options: option.Options{
						Source: "curl",
						Host:   "https://cluster",
					},
				},
				DumpName:         "token_dump",
				DumpDirRemote:    "/snapshots",
				DumpNameTemplate: "repo_token",
			},
			expectedContains: []string{
				`Authorization: Bearer abc123`,
			},
			expectedArchive: "token_dump.tar.gz",
		},

		{
			name: "custom certificate",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port: "9200",
					Options: option.Options{
						Source:          "curl",
						Host:            "https://secure-node",
						PathCertificate: "/certs/ca.pem",
					},
				},
				DumpName:         "cert_dump",
				DumpDirRemote:    "/mnt",
				DumpNameTemplate: "repo_cert",
			},
			expectedContains: []string{
				"--cacert /certs/ca.pem",
			},
			expectedArchive: "cert_dump.tar.gz",
		},

		{
			name: "repository creation fallback present",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port: "9200",
					Options: option.Options{
						Source: "curl",
						Host:   "http://os",
					},
				},
				DumpName:         "fallback",
				DumpDirRemote:    "/opt",
				DumpNameTemplate: "repo_fb",
			},
			expectedContains: []string{
				"|| curl -X PUT",
			},
			expectedArchive: "fallback.tar.gz",
		},
	}

	gen := opensearch.Generator{}

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

func TestOpenSearchGenerator_CommandIntegrity(t *testing.T) {

	cfg := &cmdCfg.Config{
		Database: cmdCfg.Database{
			Port: "9200",
			Options: option.Options{
				Source: "curl",
				Host:   "http://localhost",
			},
		},
		DumpName:         "prod_backup",
		DumpDirRemote:    "/data",
		DumpNameTemplate: "repo_prod",
	}

	gen := opensearch.Generator{}

	cmd, err := gen.Generate(cfg)

	require.NoError(t, err)

	assert.Contains(t, cmd.Command,
		"_snapshot/repo_prod",
	)

	assert.Contains(t, cmd.Command,
		"wait_for_completion=true",
	)

	assert.Contains(t, cmd.Command,
		"tar -czf prod_backup.tar.gz",
	)

	assert.Equal(t, "prod_backup.tar.gz", cmd.DumpPath)
}
