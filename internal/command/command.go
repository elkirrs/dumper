package command

import (
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type Settings struct {
	Config *cmdCfg.ConfigData
	AppCfg *setting.Settings
}

func NewApp(appCfg *setting.Settings, config *cmdCfg.ConfigData) *Settings {
	return &Settings{
		Config: config,
		AppCfg: appCfg,
	}
}

func (s *Settings) GetCommand() (string, string, error) {
	gen, ok := GetGenerator(s.Config.Driver)
	if !ok {
		return "", "", fmt.Errorf("unsupported driver: %s ", s.Config.Driver)
	}

	cmd, remotePath := gen.Generate(s.Config, s.AppCfg)
	return cmd, remotePath, nil
}
