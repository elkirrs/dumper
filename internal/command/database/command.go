package database

import (
	"context"
	"dumper/internal/command/database/mariadb"
	"dumper/internal/command/database/mongodb"
	"dumper/internal/command/database/mssql"
	"dumper/internal/command/database/mysql"
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

func (s *Settings) GetCommand() (*commandDomain.DBCommand, error) {

	var gen Generator

	switch s.Config.Database.Driver {
	case "psql":
		gen = postgres.PSQLGenerator{}
	case "mysql":
		gen = mysql.MySQLGenerator{}
	case "mongo":
		gen = mongodb.MongoGenerator{}
	case "sqlite":
		gen = sqlite.SQLiteGenerator{}
	case "mariadb":
		gen = mariadb.MariaDbGenerator{}
	case "redis":
		gen = redis.RedisGenerator{}
	case "mssql":
		gen = mssql.MSQLGenerator{}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", s.Config.Database.Driver)
	}

	cmdData, err := gen.Generate(s.Config)

	if err != nil {
		return nil, err
	}

	if s.Config.Database.Docker.Enabled {
		docker := docker.NewApp(s.ctx, cmdData, s.Config)
		docker.Prepare()
	}

	return cmdData, nil
}
