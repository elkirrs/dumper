package elasticsearch

import (
	commandDomain "dumper/internal/domain/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/pkg/utils/template"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Generator struct{}

func (g *Generator) Generate(data *cmdCfg.Config) (*commandDomain.DBCommand, error) {
	ext := ".tar.gz"

	authData := curlAuth(data)
	snapshotFullDir := template.GetFullPath(data.Database.Options.SnapPath, data.DumpNameTemplate)

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
		createRepository(snapshotFullDir),
		authData,
	)

	baseCmd += fmt.Sprintf(
		" && %s -X PUT %s:%s/_snapshot/%s/%s?wait_for_completion=true -H \"Content-Type: application/json\" -d '%s' %s",
		data.Database.Options.Source,
		data.Database.Options.Host,
		data.Database.Port,
		data.DumpNameTemplate,
		strconv.FormatInt(time.Now().Unix(), 10),
		createBodySnapshot(data),
		authData,
	)

	archiveNamePath := fmt.Sprintf("%s%s", data.DumpName, ext)

	baseCmd += fmt.Sprintf(
		" && tar -czf %s -C %s %s",
		archiveNamePath,
		data.Database.Options.SnapPath,
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
		snapshotFullDir,
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

	if data.Database.Options.CACertPath != "" {
		auth += fmt.Sprintf("--cacert %s ", data.Database.Options.CACertPath)
	}

	if data.Database.Options.KeyPath != "" {
		auth += fmt.Sprintf("--key %s ", data.Database.Options.KeyPath)
	}

	if data.Database.Options.CertPath != "" {
		if data.Database.Options.KeyPass != "" {
			auth += fmt.Sprintf(
				"--cert %s:%s ",
				data.Database.Options.CertPath,
				data.Database.Options.KeyPass,
			)
		} else {
			auth += fmt.Sprintf("--cert %s ", data.Database.Options.CertPath)
		}
	}

	if len(auth) == 0 {
		auth += "-k "
	}

	return auth
}

func createRepository(locationSnapshot string) string {

	repository := make(map[string]any)
	repository["type"] = "fs"

	settings := make(map[string]string)

	settings["location"] = locationSnapshot

	repository["settings"] = settings

	jsonData, _ := json.Marshal(repository)
	return string(jsonData)
}

func createBodySnapshot(data *cmdCfg.Config) string {
	body := make(map[string]any)

	body["indices"] = "*"
	if len(data.Database.Options.Indices) != 0 {
		body["indices"] = strings.Join(data.Database.Options.Indices, ",")
	}

	body["ignore_unavailable"] = data.Database.Options.IgnoreUnavailable
	body["include_global_state"] = data.Database.Options.IncludeGlobalState

	jsonData, _ := json.Marshal(body)
	return string(jsonData)
}
