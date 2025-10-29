package database_test

import (
	"dumper/internal/command/database"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildConfig creates a safe default config for all drivers
func buildConfig(driver string, format ...string) *cmdCfg.Config {
	ssl := false
	dbFormat := "default"
	if len(format) > 0 {
		dbFormat = format[0]
	}

	return &cmdCfg.Config{
		Database: cmdCfg.Database{
			Driver:   driver,
			Format:   dbFormat,
			User:     "user",
			Password: "pass",
			Name:     "db",
			Port:     "1234",
			Options: option.Options{
				SSL:  &ssl,
				Path: "db.sqlite",
				Mode: "",
			},
		},
		DumpName:     "dumpfile",
		Archive:      false,
		DumpLocation: "local",
		Server:       cmdCfg.Server{Host: "localhost"},
	}
}

func TestSettings_GetCommand_AllDrivers(t *testing.T) {
	tests := []struct {
		name      string
		driver    string
		expectErr bool
		format    string
	}{
		{"Postgres driver", "psql", false, ""},
		{"MySQL driver", "mysql", false, ""},
		{"MongoDB driver", "mongo", false, ""},
		{"SQLite driver", "sqlite", false, ""},
		{"MariaDB driver", "mariadb", false, ""},
		{"Redis driver", "redis", false, ""},
		{"MSSQL driver", "mssql", false, "bac"}, // valid MSSQL format
		{"Unsupported driver", "unknown", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := buildConfig(tt.driver)
			if tt.driver == "mssql" && tt.format != "" {
				cfg.Database.Format = tt.format
			}

			app := database.NewApp(cfg)
			cmd, err := app.GetCommand()

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, cmd)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cmd)
			assert.IsType(t, &commandDomain.DBCommand{}, cmd)
			assert.NotEmpty(t, cmd.Command)
			assert.NotEmpty(t, cmd.DumpPath)
		})
	}
}

func TestSettings_GetCommand_PropagatesGeneratorError(t *testing.T) {
	// Force MSSQL unsupported format to generate an error
	cfg := buildConfig("mssql", "unsupported_format")

	app := database.NewApp(cfg)
	cmd, err := app.GetCommand()

	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "unsupported database format")
}

func TestSettings_GetCommand_UnsupportedDriver(t *testing.T) {
	cfg := buildConfig("nope")
	app := database.NewApp(cfg)

	cmd, err := app.GetCommand()
	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "unsupported database driver")
}
