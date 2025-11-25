package validation

import (
	"dumper/internal/domain/config"
	"fmt"
)

func validateServer(v *Validation, cfg *config.Config) error {
	for name, srv := range cfg.Servers {

		srv.Name = srv.GetName()
		srv.Port = srv.GetPort(&cfg.Settings.SrvPost)
		srv.PrivateKey = srv.GetPrivateKey(&cfg.Settings.SSH.PrivateKey)
		srv.Password = srv.GetPrivateKey(&cfg.Settings.SSH.Password)

		if srv.PrivateKey == "" && srv.Password == "" {
			return fmt.Errorf("server %s invalid: Private key or Password are required if not set in global", name)
		}

		if err := v.validator.Struct(srv); err != nil {
			return fmt.Errorf("server '%s' invalid: %w", name, HumanError(err))
		}
	}

	return nil
}
