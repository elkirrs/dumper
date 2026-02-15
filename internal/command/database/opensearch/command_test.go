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
	trueVal := true
	falseVal := false
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
						Source:   "curl",
						Host:     "http://localhost",
						SnapPath: "/backup",
					},
				},
				DumpName:         "dump",
				DumpNameTemplate: "repo",
			},
			expectedContains: []string{
				"curl -f -X GET http://localhost:9200/_snapshot/repo",
				"_snapshot/repo/",
				"tar -czf dump.tar.gz",
				"rm -Rf /backup/repo",
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
						Source:   "curl",
						Host:     "https://node",
						SnapPath: "/data",
					},
				},
				DumpName:         "secure",
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
						Source:   "curl",
						Host:     "https://cluster",
						SnapPath: "/snapshots",
					},
				},
				DumpName:         "token_dump",
				DumpNameTemplate: "repo_token",
			},
			expectedContains: []string{
				`Authorization: Bearer abc123`,
			},
			expectedArchive: "token_dump.tar.gz",
		},
		{
			name: "custom certificate with key pass",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port: "9200",
					Options: option.Options{
						Source:     "curl",
						Host:       "https://secure-node",
						CACertPath: "/certs/ca.crt",
						CertPath:   "/certs/crt.crt",
						KeyPath:    "/certs/private.key",
						KeyPass:    "keypass",
						SnapPath:   "/mnt",
					},
				},
				DumpName:         "cert_dump",
				DumpNameTemplate: "repo_cert",
			},
			expectedContains: []string{
				"--cacert /certs/ca.crt",
				"--key /certs/private.key",
				"--cert /certs/crt.crt:keypass",
			},
			expectedArchive: "cert_dump.tar.gz",
		},
		{
			name: "repository creation fallback present",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port: "9200",
					Options: option.Options{
						Source:   "curl",
						Host:     "http://os",
						SnapPath: "/opt",
					},
				},
				DumpName:         "fallback",
				DumpNameTemplate: "repo_fb",
			},
			expectedContains: []string{
				"|| curl -X PUT",
			},
			expectedArchive: "fallback.tar.gz",
		},
		{
			name: "with custom indices and global state options",
			config: &cmdCfg.Config{
				Database: cmdCfg.Database{
					Port: "9200",
					Options: option.Options{
						Source:             "curl",
						Host:               "http://localhost",
						SnapPath:           "/snapshots",
						Indices:            []string{"logs", "metrics"},
						IgnoreUnavailable:  &trueVal,
						IncludeGlobalState: &falseVal,
					},
				},
				DumpName:         "indices_backup",
				DumpNameTemplate: "repo_indices",
			},
			expectedContains: []string{
				`"indices":"logs,metrics"`,
				`"ignore_unavailable":true`,
				`"include_global_state":false`,
			},
			expectedArchive: "indices_backup.tar.gz",
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
				Source:   "curl",
				Host:     "http://localhost",
				SnapPath: "/data",
			},
		},
		DumpName:         "prod_backup",
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
