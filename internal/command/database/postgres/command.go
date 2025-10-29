package postgres

import (
	commandDomain "dumper/internal/domain/command"
	commandConfig "dumper/internal/domain/command-config"
	"fmt"
)

type PSQLGenerator struct{}

func (g PSQLGenerator) Generate(data *commandConfig.Config) (*commandDomain.DBCommand, error) {
	formatFlag := "-Fp" // plain SQL
	ext := "sql"

	switch data.Database.Format {
	case "dump":
		formatFlag = "-Fc"
		ext = "dump"
	case "tar":
		formatFlag = "-Ft"
		ext = "tar"
	}

	baseCmd := fmt.Sprintf(
		"/usr/bin/pg_dump --dbname=postgresql://%s:%s@127.0.0.1:%s/%s --clean --if-exists --no-owner %s",
		data.Database.User,
		data.Database.Password,
		data.Database.Port,
		data.Database.Name,
		formatFlag,
	)

	if data.Archive && formatFlag == "-Fp" { // gzip only for plain
		baseCmd += " | gzip"
		ext += ".gz"
	}

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	if data.DumpLocation == "server" {
		return &commandDomain.DBCommand{
			Command:  fmt.Sprintf("%s > %s", baseCmd, remotePath),
			DumpPath: remotePath,
		}, nil
	}
	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: remotePath,
	}, nil

}
