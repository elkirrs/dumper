package shell

import (
	"context"
	"dumper/internal/connect"
	commandConfig "dumper/internal/domain/command-config"
	"dumper/pkg/logging"
	"fmt"
)

type Shell struct {
	ctx    context.Context
	config *commandConfig.Config
	conn   *connect.Connect
}

func NewApp(
	ctx context.Context,
	config *commandConfig.Config,
	conn *connect.Connect,
) *Shell {
	return &Shell{
		ctx:    ctx,
		config: config,
		conn:   conn,
	}
}

func (s *Shell) RunScriptBefore() error {
	if !s.config.Server.Shell.Enabled || s.config.Server.Shell.Before == "" {
		return nil
	}

	fmt.Println("Run shell script before start backup")
	if s.config.Server.Shell.Before != "" {
		logging.L(s.ctx).Info("Run shell script before start backup")
		return s.runScript("before", s.config.Server.Shell.Before)
	}

	return nil
}

func (s *Shell) RunScriptAfter() error {
	if !s.config.Server.Shell.Enabled || s.config.Server.Shell.After == "" {
		return nil
	}

	fmt.Println("Run shell script after finished backup")
	if s.config.Server.Shell.After != "" {
		logging.L(s.ctx).Info("Run shell script after finished backup")
		return s.runScript("after", s.config.Server.Shell.Before)
	}
	return nil
}

func (s *Shell) runScript(
	shellType string,
	script string,
) error {
	scriptLog := fmt.Sprintf("Run shell script '%s' backup", shellType)
	msg, err := s.conn.RunCommand(script)
	if err != nil {
		logging.L(s.ctx).Error(
			scriptLog,
			logging.StringAttr("msg", msg),
			logging.ErrAttr(err),
		)
		return err
	}
	logging.L(s.ctx).Info(scriptLog, logging.StringAttr("msg", msg))
	return nil
}
