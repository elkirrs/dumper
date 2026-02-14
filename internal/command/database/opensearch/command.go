package opensearch

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Generator struct{}

func (g *Generator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := ".tar.gz"

	authData := curlAuth(data)

	baseCmd := fmt.Sprintf(
		"%s -f -X GET %s:%s/_snapshot/%s %s",
		data.Database.Options.Source,
		data.Database.Options.Host,
		data.Database.Port,
		data.DumpNameTemplate,
		authData,
	)

	baseCmd += fmt.Sprintf(
		" || %s -X PUT %s:%s/_snapshot/%s -H \"Content-Type: application/json\" -d '%s' %s",
		data.Database.Options.Source,
		data.Database.Options.Host,
		data.Database.Port,
		data.DumpNameTemplate,
		createRepository(data),
		authData,
	)

	baseCmd += fmt.Sprintf(
		" && %s -X PUT %s:%s/_snapshot/%s/%s?wait_for_completion=true -H \"Content-Type: application/json\" %s",
		data.Database.Options.Source,
		data.Database.Options.Host,
		data.Database.Port,
		data.DumpNameTemplate,
		strconv.FormatInt(time.Now().Unix(), 10),
		authData,
	)

	archiveNamePath := fmt.Sprintf("%s%s", data.DumpName, ext)

	baseCmd += fmt.Sprintf(
		" && tar -czf %s -C %s %s",
		archiveNamePath,
		data.DumpDirRemote,
		data.DumpNameTemplate,
	)

	baseCmd += fmt.Sprintf(
		" && %s -X DELETE %s:%s/_snapshot/%s -H \"Content-Type: application/json\" %s",
		data.Database.Options.Source,
		data.Database.Options.Host,
		data.Database.Port,
		data.DumpNameTemplate,
		authData,
	)

	baseCmd += fmt.Sprintf(
		"&& rm -Rf %s",
		data.DumpName,
	)

	return &commandDomain.DBCommand{
		Command:  baseCmd,
		DumpPath: archiveNamePath,
	}, nil
}

func curlAuth(data *cmdCfg.Config) string {

	auth := ""

	if data.Database.User != "" {
		auth += fmt.Sprintf("-u %s:%s ", data.Database.User, data.Database.Password)
	}

	if data.Database.Token != "" {
		auth += fmt.Sprintf(`-H "Authorization: %s" `, data.Database.Token)
	}

	if data.Database.Options.PathCertificate != "" {
		auth += fmt.Sprintf("--cacert %s ", data.Database.Options.PathCertificate)
	} else {
		auth += "-k "
	}

	return auth
}

func createRepository(data *cmdCfg.Config) string {

	repository := make(map[string]any)
	repository["type"] = "fs"

	settings := make(map[string]any)

	settings["location"] = data.DumpName

	repository["settings"] = settings

	jsonData, _ := json.Marshal(repository)

	return string(jsonData)
}
