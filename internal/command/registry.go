package command

import (
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
)

type CmdGenerator interface {
	Generate(*cmdCfg.ConfigData, *setting.Settings) (cmd string, remotePath string)
}

var generators = map[string]CmdGenerator{}

func Register(driver string, gen CmdGenerator) {
	generators[driver] = gen
}

func GetGenerator(driver string) (CmdGenerator, bool) {
	gen, ok := generators[driver]
	return gen, ok
}
