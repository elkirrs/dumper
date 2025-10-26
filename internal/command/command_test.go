package command

import (
	"dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGenerator struct{}

func (m mockGenerator) Generate(data *command_config.ConfigData, settings *setting.Settings) (string, string) {
	return "mock-cmd", "mock-remote"
}

func TestSettings_GetCommand(t *testing.T) {
	Register("mockdriver", mockGenerator{})

	tests := []struct {
		name       string
		driver     string
		wantCmd    string
		wantRemote string
		wantErr    error
	}{
		{
			name:       "supported driver",
			driver:     "mockdriver",
			wantCmd:    "mock-cmd",
			wantRemote: "mock-remote",
			wantErr:    nil,
		},
		{
			name:       "unsupported driver",
			driver:     "unknown",
			wantCmd:    "",
			wantRemote: "",
			wantErr:    errors.New("unsupported driver: unknown "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &command_config.ConfigData{
				Driver: tt.driver,
			}
			appCfg := &setting.Settings{}

			s := NewApp(appCfg, config)
			require.NotNil(t, s)

			cmd, remote, err := s.GetCommand()

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.Equal(t, "", cmd)
				assert.Equal(t, "", remote)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCmd, cmd)
				assert.Equal(t, tt.wantRemote, remote)
			}
		})
	}
}
