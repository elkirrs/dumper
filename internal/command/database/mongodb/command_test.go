package mongodb_test

import (
	"dumper/pkg/utils"
	"testing"

	"dumper/internal/command/database/mongodb"
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/option"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongoGenerator_Generate(t *testing.T) {
	source := utils.GetDBSource("mongodb", "")
	sslTrue := true

	tests := []struct {
		name       string
		cfg        *cmdCfg.Config
		wantCmd    string
		wantDump   string
		shouldFail bool
	}{
		{
			name: "archive format + --archive flag + gzip",
			cfg: &cmdCfg.Config{
				DumpName:     "backup",
				Archive:      true,
				DumpLocation: "server",
				Database: cmdCfg.Database{
					Name:     "mydb",
					User:     "user",
					Password: "pass",
					Port:     "27017",
					Format:   "archive",
					Options: option.Options{
						AuthSource: "admin",
						SSL:        &sslTrue,
						Source:     source,
					},
				},
			},
			wantCmd:  `mongodump --uri "mongodb://user:pass@127.0.0.1:27017/mydb?authSource=admin&ssl=true" --archive=backup.gz --gzip`,
			wantDump: "backup.gz",
		},
		{
			name: "non-archive format + tar.gz after dump",
			cfg: &cmdCfg.Config{
				DumpName:     "dump1",
				Archive:      true,
				DumpLocation: "server",
				Database: cmdCfg.Database{
					Name:     "mydb",
					User:     "root",
					Password: "qwerty",
					Port:     "27018",
					Format:   "",
					Options: option.Options{
						AuthSource: "admin",
						SSL:        &sslTrue,
						Source:     source,
					},
				},
			},
			wantCmd:  `mongodump --uri "mongodb://root:qwerty@127.0.0.1:27018/mydb?authSource=admin&ssl=true" --out ./ && tar -czf dump1.tar.gz mydb`,
			wantDump: "dump1.tar.gz",
		},
		{
			name: "no archive, format=archive, plain archive output",
			cfg: &cmdCfg.Config{
				DumpName:     "ddd",
				Archive:      false,
				DumpLocation: "server",
				Database: cmdCfg.Database{
					Name:     "db",
					User:     "u",
					Password: "p",
					Port:     "27017",
					Format:   "archive",
					Options: option.Options{
						AuthSource: "admin",
						SSL:        &sslTrue,
						Source:     source,
					},
				},
			},
			wantCmd:  `mongodump --uri "mongodb://u:p@127.0.0.1:27017/db?authSource=admin&ssl=true" --archive=ddd.archive`,
			wantDump: "ddd.archive",
		},
		{
			name: "no archive, format=dir ⇒ tar.gz",
			cfg: &cmdCfg.Config{
				DumpName:     "test2",
				Archive:      false,
				DumpLocation: "server",
				Database: cmdCfg.Database{
					Name:     "mydb",
					User:     "u",
					Password: "pp",
					Port:     "27017",
					Format:   "",
					Options: option.Options{
						AuthSource: "",
						SSL:        &sslTrue,
						Source:     source,
					},
				},
			},
			wantCmd:  `--out ./ mongodump --uri "mongodb://u:pp@127.0.0.1:27017/mydb?ssl=true" && tar -czf test2.tar.gz mydb`,
			wantDump: "test2.tar.gz",
		},
		{
			name: "special characters in username/password must be escaped",
			cfg: &cmdCfg.Config{
				DumpName:     "sp",
				Archive:      false,
				DumpLocation: "server",
				Database: cmdCfg.Database{
					Name:     "db",
					User:     "us er",
					Password: "p@ss",
					Port:     "27017",
					Format:   "archive",
					Options: option.Options{
						AuthSource: "adm in",
						SSL:        &sslTrue,
						Source:     source,
					},
				},
			},
			wantCmd:  `mongodump --uri "mongodb://us+er:p%40ss@127.0.0.1:27017/db?authSource=adm+in&ssl=true" --archive=sp.archive`,
			wantDump: "sp.archive",
		},
	}

	gen := mongodb.MongoGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := gen.Generate(tt.cfg)

			if tt.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)

			assert.Equal(t, tt.wantCmd, res.Command)
			assert.Equal(t, tt.wantDump, res.DumpPath)

			_, ok := interface{}(res).(*commandDomain.DBCommand)
			assert.True(t, ok)
		})
	}
}
