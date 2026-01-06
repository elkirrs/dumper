package local_config

import (
	"dumper/internal/crypt"
	"dumper/internal/domain/config"
	"dumper/internal/domain/config/docker"
	"dumper/internal/domain/config/encrypt"
	"dumper/internal/domain/config/setting"
	"dumper/internal/domain/config/shell"
	sshConfig "dumper/internal/domain/config/ssh-config"
	"dumper/internal/validation"
	crypt2 "dumper/pkg/utils/crypt"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

func Load(filename, appSecret string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if crypt2.IsEncrypted(data) && crypt2.LooksEncrypted(data) {
		cFile, err := crypt2.ReadEncryptedFile(filename)

		if err != nil {
			return nil, err
		}

		scope := crypt2.GetScope(cFile.Header)
		data, err = crypt.DecryptInApp(cFile.Data, scope, appSecret)
		if err != nil {
			return nil, err
		}
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Settings == nil {
		cfg.Settings = &setting.Settings{}
	}
	if cfg.Settings.SSH == nil {
		cfg.Settings.SSH = &sshConfig.SSHConfig{}
	}
	if cfg.Settings.Encrypt == nil {
		cfg.Settings.Encrypt = &encrypt.Encrypt{}
	}
	if cfg.Settings.Docker == nil {
		cfg.Settings.Docker = &docker.Docker{}
	}
	if cfg.Settings.Shell == nil {
		cfg.Settings.Shell = &shell.Shell{}
	}

	_ = defaults.Set(cfg.Settings.SSH)
	_ = defaults.Set(cfg.Settings.Encrypt)
	_ = defaults.Set(cfg.Settings.Docker)
	_ = defaults.Set(cfg.Settings.Shell)
	_ = defaults.Set(cfg.Settings)
	_ = defaults.Set(&cfg)

	v := validation.New()
	if err := v.Handler(&cfg); err != nil {
		return nil, validation.HumanError(err)
	}

	return &cfg, nil
}
