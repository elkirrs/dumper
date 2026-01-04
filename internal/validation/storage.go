package validation

import (
	"dumper/internal/domain/config"
	"dumper/internal/domain/config/storage"
	"fmt"
)

func validateStorages(v *Validation, cfg *config.Config) error {
	validate := v.validator

	for name, s := range cfg.Storages {
		switch s.Type {
		case "local":
			local := storage.Local{
				Type: s.Type,
				Dir:  s.Dir,
			}
			if err := validate.Struct(local); err != nil {
				return fmt.Errorf("storage '%s' (local) invalid: %w", name, HumanError(err))
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
				return fmt.Errorf("storage '%s' (ftp) invalid: %w", name, HumanError(err))
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
				return fmt.Errorf("storage '%s' (sftp) invalid: %w", name, HumanError(err))
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
					return fmt.Errorf("storage '%s' (azure-shared-key) invalid: %w", name, HumanError(err))
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
					return fmt.Errorf("storage '%s' (azure-ad) invalid: %w", name, HumanError(err))
				}

			default:
				return fmt.Errorf("storage '%s': unknown azure auth_type '%s'", name, s.AuthType)
			}

		case "s3":
			s3 := storage.S3{
				Type:      s.Type,
				Region:    s.Region,
				Bucket:    s.Bucket,
				AccessKey: s.AccessKey,
				SecretKey: s.SecretKey,
			}

			if err := validate.Struct(s3); err != nil {
				return fmt.Errorf("storage '%s' (s3) invalid: %w", name, HumanError(err))
			}

		case "minio":
			minio := storage.MinIO{
				Type:      s.Type,
				Region:    s.Region,
				Bucket:    s.Bucket,
				AccessKey: s.AccessKey,
				SecretKey: s.SecretKey,
				Endpoint:  s.Endpoint,
			}

			if err := validate.Struct(minio); err != nil {
				return fmt.Errorf("storage '%s' (MinIO) invalid: %w", name, HumanError(err))
			}

		case "r2":
			r2 := storage.Cloudflare{
				Type:      s.Type,
				Bucket:    s.Bucket,
				AccessKey: s.AccessKey,
				SecretKey: s.SecretKey,
				Endpoint:  s.Endpoint,
				AccountID: s.AccountID,
			}

			if err := validate.Struct(r2); err != nil {
				return fmt.Errorf("storage '%s' (Cloudflare R2) invalid: %w", name, HumanError(err))
			}

		case "b2":
			b2 := storage.Backblaze{
				Type:      s.Type,
				Bucket:    s.Bucket,
				AccessKey: s.AccessKey,
				SecretKey: s.SecretKey,
				Endpoint:  s.Endpoint,
				Region:    s.Region,
			}

			if err := validate.Struct(b2); err != nil {
				return fmt.Errorf("storage '%s' (Backblaze) invalid: %w", name, HumanError(err))
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
