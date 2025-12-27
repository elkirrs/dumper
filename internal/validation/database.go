package validation

import (
	"dumper/internal/domain/config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils"
	"fmt"

	"github.com/creasty/defaults"
)

func validateDatabase(v *Validation, cfg *config.Config) error {
	for name, db := range cfg.Databases {

		db.Name = db.GetName()
		db.Port = db.GetPort(&cfg.Settings.DBPort)
		db.Driver = db.GetDriver(&cfg.Settings.Driver)
		db.Format = db.GetFormat(&cfg.Settings.DumpFormat)
		db.Storages = db.GetStorages(&cfg.Settings.Storages)
		db.DirRemote = db.GetDirRemote(&cfg.Settings.DirRemote)
		isArchive := db.IsArchive(*cfg.Settings.Archive)
		db.Archive = &isArchive

		if db.Options == nil {
			db.Options = &option.Options{}
			_ = defaults.Set(db.Options)
		}

		if db.Options.Source == "" {
			db.Options.Source = utils.GetDBSource(db.Driver, db.Format)
		}

		cfg.Databases[name] = db

		if err := v.validator.Struct(db); err != nil {
			return fmt.Errorf("database '%s' invalid: %w", name, HumanError(err))
		}

		if _, ok := cfg.Servers[db.Server]; !ok {
			return fmt.Errorf("database '%s' has invalid server '%s'", name, db.Server)
		}

	}

	return nil
}
