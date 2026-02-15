package validation

import (
	"dumper/internal/domain/config"
	"dumper/internal/domain/config/option"
	"dumper/pkg/utils/mapping"
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
		docker := db.GetDocker(cfg.Settings.Docker)
		db.Docker = &docker
		removeDump := db.GetRemoveDump(cfg.Settings.RemoveDump)
		db.RemoveDump = &removeDump
		isArchive := db.IsArchive(*cfg.Settings.Archive)
		db.Archive = &isArchive

		if db.Port == "" {
			db.Port = mapping.GetDefaultDBPort(db.Driver)
		}

		if db.Options == nil {
			db.Options = &option.Options{}
		}

		_ = defaults.Set(db.Options)

		if db.Options.Source == "" {
			db.Options.Source = mapping.GetDBSource(db.Driver, db.Format)
		}

		if db.Options.SnapPath == "" {
			db.Options.SnapPath = db.DirRemote
		}

		cfg.Databases[name] = db

		if ok := mapping.IsValidFormatDump(db.Driver, db.Format); !ok {
			return fmt.Errorf("database '%s' invalid driver: '%s' or invalid format: '%s'", name, db.Driver, db.Format)
		}

		if err := v.validator.Struct(db); err != nil {
			return fmt.Errorf("database '%s' invalid: %w", name, HumanError(err))
		}

		if _, ok := cfg.Servers[db.Server]; !ok {
			return fmt.Errorf("database '%s' has invalid server '%s'", name, db.Server)
		}

	}

	return nil
}
