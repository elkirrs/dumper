package command

import (
	"dumper/internal/config"
	cmdCfg "dumper/internal/domain/command-config"
)

type CmdGenerator interface {
	Generate(*cmdCfg.ConfigData, *config.Settings) (cmd string, remotePath string)
}

var generators = map[string]CmdGenerator{}

func Register(driver string, gen CmdGenerator) {
	generators[driver] = gen
}

func GetGenerator(driver string) (CmdGenerator, bool) {
	gen, ok := generators[driver]
	return gen, ok
}
