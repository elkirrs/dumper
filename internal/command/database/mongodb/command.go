package mongodb

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
)

type MongoGenerator struct{}

func (g MongoGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := "bson"

	formatFlag := ""
	switch data.Database.Format {
	case "archive":
		formatFlag = "--archive"
		ext = "archive"
	default:
		formatFlag = "--dump"
		ext = "bson"
	}

	uri := fmt.Sprintf("mongodb://%s:%s@127.0.0.1:%s/%s",
		data.Database.User, data.Database.Password, data.Database.Port, data.Database.Name)

	params := ""
	if data.Database.Options.AuthSource != "" {
		params += fmt.Sprintf("authSource=%s", data.Database.Options.AuthSource)
	}

	if *data.Database.Options.SSL {
		if len(params) > 0 {
			params += "&"
		}
		params += fmt.Sprintf("ssl=%t", *data.Database.Options.SSL)
	}

	if len(params) > 0 {
		uri += "?" + params
	}

	baseCmd := fmt.Sprintf("/usr/bin/mongodump --uri \"%s\"", uri)

	if formatFlag == "--archive" {
		baseCmd += " --archive"
	} else {
		baseCmd += " --out ./"
	}

	if data.Archive {
		if formatFlag == "--archive" {
			baseCmd += " --gzip"
			ext += ".gz"
		} else {
			baseCmd = fmt.Sprintf("%s && tar -czf %s.tar.gz %s", baseCmd, data.DumpName, data.Database.Name)
			ext = "tar.gz"
		}
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
