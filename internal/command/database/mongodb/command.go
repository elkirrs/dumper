package mongodb

import (
	command "dumper/internal/command/database"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type MongoGenerator struct{}

func (g MongoGenerator) Generate(data *cmdCfg.ConfigData, settings *setting.Settings) (string, string) {
	if data.Port == "" {
		data.Port = "27017"
	}

	ext := "bson"

	formatFlag := ""
	switch data.DumpFormat {
	case "archive":
		formatFlag = "--archive"
		ext = "archive"
	default:
		formatFlag = "--dump"
		ext = "bson"
	}

	uri := fmt.Sprintf("mongodb://%s:%s@127.0.0.1:%s/%s",
		data.User, data.Password, data.Port, data.Name)

	params := ""
	if data.Options.AuthSource != "" {
		params += fmt.Sprintf("authSource=%s", data.Options.AuthSource)
	}

	if *data.Options.SSL {
		if len(params) > 0 {
			params += "&"
		}
		params += fmt.Sprintf("ssl=%t", *data.Options.SSL)
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

	if *settings.Archive {
		if formatFlag == "--archive" {
			baseCmd += " --gzip"
			ext += ".gz"
		} else {
			baseCmd = fmt.Sprintf("%s && tar -czf %s.tar.gz %s", baseCmd, data.DumpName, data.Name)
			ext = "tar.gz"
		}
	}

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	if settings.DumpLocation == "server" {
		return fmt.Sprintf("%s > %s", baseCmd, remotePath), remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("mongo", MongoGenerator{})
}
