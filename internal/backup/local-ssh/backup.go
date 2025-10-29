package local_ssh

import (
	"context"
	"dumper/internal/connect"
	commandConfig "dumper/internal/domain/command-config"
)

type BackupServer struct {
	ctx    context.Context
	conn   *connect.Connect
	config *commandConfig.Config
}

func NewApp(
	ctx context.Context,
	conn *connect.Connect,
	config *commandConfig.Config,
) *BackupServer {
	return &BackupServer{
		ctx:    ctx,
		conn:   conn,
		config: config,
	}
}

func (b *BackupServer) Run() error {
	panic("implement me")
}
