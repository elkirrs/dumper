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
	source := utils.GetDBSource("mongo", "")
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
			wantCmd:  `mongodump --uri "mongodb://user:pass@127.0.0.1:27017/?authSource=admin&ssl=true" --db mydb --archive=backup.gz --gzip`,
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
			wantCmd:  `mongodump --uri "mongodb://root:qwerty@127.0.0.1:27018/?authSource=admin&ssl=true" --db mydb --out ./ && tar -czf dump1.tar.gz mydb`,
			wantDump: "dump1.tar.gz",
		},
		{
			name: "no archive + archive format",
			cfg: &cmdCfg.Config{
				DumpName:     "ddd",
				Archive:      false,
				DumpLocation: "server",
				Database: cmdCfg.Database{
					Name:     "mydb",
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
			wantCmd:  `mongodump --uri "mongodb://u:p@127.0.0.1:27017/?authSource=admin&ssl=true" --db mydb --archive=ddd.archive`,
			wantDump: "ddd.archive",
		},
		{
			name: "no archive + directory format",
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
			wantCmd:  `mongodump --uri "mongodb://u:pp@127.0.0.1:27017/?ssl=true" --db mydb --out ./ && tar -czf test2.tar.gz mydb`,
			wantDump: "test2.tar.gz",
		},
		{
			name: "special characters in username/password must be escaped",
			cfg: &cmdCfg.Config{
				DumpName:     "sp",
				Archive:      false,
				DumpLocation: "server",
				Database: cmdCfg.Database{
					Name:     "mydb",
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
			wantCmd:  `mongodump --uri "mongodb://us+er:p%40ss@127.0.0.1:27017/?authSource=adm+in&ssl=true" --db mydb --archive=sp.archive`,
			wantDump: "sp.archive",
		},
		{
			name: "include tables in dump",
			cfg: &cmdCfg.Config{
				DumpName:     "sp",
				Archive:      false,
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
						IncTables:  []string{"t1"},
					},
				},
			},
			wantCmd:  `mongodump --uri "mongodb://user:pass@127.0.0.1:27017/?authSource=admin&ssl=true" --db mydb  --collection t1 --archive=sp.archive`,
			wantDump: "sp.archive",
		},
		{
			name: "exclude tables in dump",
			cfg: &cmdCfg.Config{
				DumpName:     "sp",
				Archive:      false,
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
						ExcTables:  []string{"t1"},
					},
				},
			},
			wantCmd:  `mongodump --uri "mongodb://user:pass@127.0.0.1:27017/?authSource=admin&ssl=true" --db mydb  --excludeCollection t1 --archive=sp.archive`,
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
