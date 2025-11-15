package docker

import (
	"context"
	commandDomain "dumper/internal/domain/command"
	commandConfig "dumper/internal/domain/command-config"
	"dumper/pkg/logging"
	"fmt"
)

type Docker struct {
	ctx     context.Context
	cmdData *commandDomain.DBCommand
	config  *commandConfig.Config
}

func NewApp(
	ctx context.Context,
	cmdData *commandDomain.DBCommand,
	config *commandConfig.Config,
) *Docker {
	return &Docker{
		ctx:     ctx,
		cmdData: cmdData,
		config:  config,
	}
}

func (d *Docker) Prepare() {
	logging.L(d.ctx).Info("Prepare docker command")

	d.cmdData.Command = fmt.Sprintf("%s %s", d.config.Database.Docker.Command, d.cmdData.Command)
}
