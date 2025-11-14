package mongodb

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"fmt"
	"net/url"
)

type MongoGenerator struct{}

func (g MongoGenerator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := ""
	formatFlag := ""

	switch data.Database.Format {
	case "archive":
		formatFlag = "--archive"
	}

	uri := fmt.Sprintf(
		"mongodb://%s:%s@127.0.0.1:%s/%s",
		url.QueryEscape(data.Database.User),
		url.QueryEscape(data.Database.Password),
		data.Database.Port,
		data.Database.Name,
	)

	params := ""
	if data.Database.Options.AuthSource != "" {
		params += fmt.Sprintf("authSource=%s", url.QueryEscape(data.Database.Options.AuthSource))
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

	baseCmd := fmt.Sprintf("mongodump --uri \"%s\"", uri)

	if data.Archive {
		if formatFlag == "--archive" {
			ext = "gz"
			baseCmd = fmt.Sprintf("%s --archive=%s.%s --gzip", baseCmd, data.DumpName, ext)
		} else {
			baseCmd = fmt.Sprintf("%s && tar -czf %s.tar.gz %s", baseCmd, data.DumpName, data.Database.Name)
			ext = "tar.gz"
		}
	} else {
		if formatFlag == "--archive" {
			ext = "archive"
			baseCmd = fmt.Sprintf("%s --archive=%s.%s", baseCmd, data.DumpName, ext)
		} else {
			baseCmd = fmt.Sprintf("--out ./ %s && tar -czf %s.tar.gz %s", baseCmd, data.DumpName, data.Database.Name)
			ext = "tar.gz"
		}
	}

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("%s", fileName)

	if data.DumpLocation == "server" {
		return &commandDomain.DBCommand{
			Command:  baseCmd,
			DumpPath: remotePath,
		}, nil
	}

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: remotePath,
	}, nil
}
