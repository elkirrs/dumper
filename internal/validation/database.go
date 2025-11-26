package validation

import (
	"dumper/internal/domain/config"
	"fmt"
)

func validateDatabase(v *Validation, cfg *config.Config) error {
	for name, db := range cfg.Databases {

		db.Name = db.GetName()
		db.Port = db.GetPort(&cfg.Settings.DBPort)
		db.Driver = db.GetDriver(&cfg.Settings.Driver)
		db.Format = db.GetFormat(&cfg.Settings.DumpFormat)
		db.Storages = db.GetStorages(&cfg.Settings.Storages)

		if err := v.validator.Struct(db); err != nil {
			return fmt.Errorf("database '%s' invalid: %w", name, HumanError(err))
		}

		if _, ok := cfg.Servers[db.Server]; !ok {
			return fmt.Errorf("database '%s' has invalid server '%s'", name, db.Server)
		}

	}

	return nil
}
