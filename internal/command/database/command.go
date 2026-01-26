package database

import (
	"context"
	"dumper/internal/command/database/db2"
	"dumper/internal/command/database/dynamodb"
	"dumper/internal/command/database/firebird"
	"dumper/internal/command/database/influxdb"
	"dumper/internal/command/database/mariadb"
	"dumper/internal/command/database/mongodb"
	"dumper/internal/command/database/mssql"
	"dumper/internal/command/database/mysql"
	"dumper/internal/command/database/neo4j"
	"dumper/internal/command/database/postgres"
	"dumper/internal/command/database/redis"
	"dumper/internal/command/database/sqlite"
	"dumper/internal/docker"
	commandDomain "dumper/internal/domain/command"
	commandConfig "dumper/internal/domain/command-config"
	"fmt"
)

type Settings struct {
	ctx    context.Context
	Config *commandConfig.Config
}

func NewApp(
	ctx context.Context,
	config *commandConfig.Config,
) *Settings {
	return &Settings{
		ctx:    ctx,
		Config: config,
	}
}

type Generator interface {
	Generate(*commandConfig.Config) (*commandDomain.DBCommand, error)
}

var dataBaseGeneratorList = map[string]Generator{
	"psql":     postgres.PSQLGenerator{},
	"mysql":    mysql.MySQLGenerator{},
	"mongo":    mongodb.MongoGenerator{},
	"sqlite":   sqlite.SQLiteGenerator{},
	"mariadb":  mariadb.MariaDbGenerator{},
	"redis":    redis.RedisGenerator{},
	"mssql":    mssql.MSQLGenerator{},
	"neo4j":    neo4j.Neo4jGenerator{},
	"dynamodb": dynamodb.DynamoDBGenerator{},
	"influxdb": influxdb.InfluxDB2Generator{},
	"db2":      db2.DB2Generator{},
	"firebird": firebird.FirebirdDbGenerator{},
}

func (s *Settings) GetCommand() (*commandDomain.DBCommand, error) {

	generator, ok := dataBaseGeneratorList[s.Config.Database.Driver]

	if !ok {
		return nil, fmt.Errorf("unsupported database driver: %s", s.Config.Database.Driver)
	}

	cmdData, err := generator.Generate(s.Config)

	if err != nil {
		return nil, err
	}

	if *s.Config.Database.Docker.Enabled {
		dockerApp := docker.NewApp(s.ctx, cmdData, s.Config)
		dockerApp.Prepare()
	}

	return cmdData, nil
}
