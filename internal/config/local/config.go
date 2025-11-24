package local_config

import (
	"dumper/internal/crypt"
	"dumper/internal/domain/config"
	"dumper/internal/domain/config/storage"
	"dumper/pkg/utils"
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

func Load(filename, appSecret string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if utils.IsEncrypted(data) && utils.LooksEncrypted(data) {
		data, err = utils.ReadEncryptedFile(filename)
		if err != nil {
			return nil, err
		}

		data, err = crypt.DecryptInApp(data, appSecret)
		if err != nil {
			return nil, err
		}
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := defaults.Set(&cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	if err := validateStorages(&cfg); err != nil {
		return nil, fmt.Errorf("storage validation failed: %w", err)
	}

	for k, srv := range cfg.Servers {
		if srv.Name == "" {
			srv.Name = srv.Host
			cfg.Servers[k] = srv
		}
	}

	return &cfg, nil
}

func validateStorages(cfg *config.Config) error {
	validate := validator.New()
	for name, s := range cfg.Storages {
		switch s.Type {
		case "local":
			local := storage.Local{
				Type: s.Type,
				Dir:  s.Dir,
			}
			if err := validate.Struct(local); err != nil {
				return fmt.Errorf("storage '%s' (local) invalid: %w", name, err)
			}
		case "azure":
			switch s.AuthType {
			case "SharedKey":
				azure := storage.AzureSharedKey{
					Type:      s.Type,
					Endpoint:  s.Endpoint,
					Container: s.Container,
					Name:      s.Name,
					SharedKey: s.SharedKey,
				}
				if err := validate.Struct(azure); err != nil {
					return fmt.Errorf("storage '%s' (azure-shared-key) invalid: %w", name, err)
				}

			case "AzureAD":
				azure := storage.AzureAD{
					Type:         s.Type,
					Endpoint:     s.Endpoint,
					Container:    s.Container,
					TenantID:     s.TenantID,
					ClientID:     s.ClientID,
					ClientSecret: s.ClientSecret,
				}

				if err := validate.Struct(azure); err != nil {
					return fmt.Errorf("storage '%s' (azure-ad) invalid: %w", name, err)
				}

			default:
				return fmt.Errorf("storage '%s': unknown azure auth_type '%s'", name, s.AuthType)
			}

		case "ftp":
			ftp := storage.FTP{
				Dir:      s.Dir,
				Host:     s.Host,
				Port:     s.Port,
				Username: s.Username,
				Password: s.Password,
			}

			if err := validate.Struct(ftp); err != nil {
				return fmt.Errorf("storage '%s' (ftp) invalid: %w", name, err)
			}

		case "sftp":
			sftp := storage.SFTP{
				Dir:        s.Dir,
				Host:       s.Host,
				Port:       s.Port,
				Username:   s.Username,
				PrivateKey: s.PrivateKey,
				Passphrase: s.Passphrase,
			}

			if err := validate.Struct(sftp); err != nil {
				return fmt.Errorf("storage '%s' (sftp) invalid: %w", name, err)
			}

		default:
			if s.Type == "" {
				return fmt.Errorf("storage '%s' missing required field 'type'", name)
			}
			return fmt.Errorf("storage '%s': unknown type '%s'", name, s.Type)
		}
	}

	return nil
}
