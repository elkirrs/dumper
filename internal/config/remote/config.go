package remote_config

import (
	"context"
	"dumper/internal/connect"
	"dumper/internal/domain/config"
	"dumper/internal/domain/config/database"
	dbConnect "dumper/internal/domain/config/db-connect"
	"dumper/internal/domain/config/server"
	"dumper/internal/validation"
	"dumper/pkg/logging"
	"fmt"

	"gopkg.in/yaml.v3"
)

type RCfg struct {
	ctx        context.Context
	conn       *connect.Connect
	cfgPath    string
	connectDBs map[string]dbConnect.DBConnect
	cfg        *config.Config
}

type Config struct {
	Databases map[string]database.Database
	Servers   map[string]server.Server
}

func New(
	ctx context.Context,
	conn *connect.Connect,
	cfgPath string,
	cfg *config.Config,
) *RCfg {
	return &RCfg{
		ctx:     ctx,
		conn:    conn,
		cfgPath: cfgPath,
		cfg:     cfg,
	}
}

func (r *RCfg) Load() error {
	checkCmd := fmt.Sprintf("test -f %s", r.cfgPath)

	logging.L(r.ctx).Info(
		"Run command found config in server with name",
		logging.StringAttr("name", r.cfgPath),
	)

	msg, err := r.conn.RunCommand(checkCmd)
	if err != nil {
		logging.L(r.ctx).Error("Failed file not exist")
		return fmt.Errorf("failed file not exist : %v", err)
	}

	logging.L(r.ctx).Info(
		"File exists on server",
		logging.StringAttr("name", r.cfgPath),
		logging.StringAttr("msg", msg),
	)

	logging.L(r.ctx).Info(
		"Run command read config in server with name",
		logging.StringAttr("name", r.cfgPath),
	)

	readCmd := fmt.Sprintf("cat %s", r.cfgPath)
	msg, err = r.conn.RunCommand(readCmd)
	if err != nil {
		logging.L(r.ctx).Error("Failed to read config")
		return fmt.Errorf("failed to read config : %v", err)
	}

	logging.L(r.ctx).Info(
		"File red on server",
		logging.StringAttr("name", r.cfgPath),
	)

	err = r.loadFromString(msg)
	if err != nil {
		logging.L(r.ctx).Error("Failed to parse config")
		return fmt.Errorf("failed to parse config : %v", err)
	}

	return nil
}

func (r *RCfg) Config() map[string]dbConnect.DBConnect {
	return r.connectDBs
}

func (r *RCfg) loadFromString(strYml string) error {
	var cfg Config

	if err := yaml.Unmarshal([]byte(strYml), &cfg); err != nil {
		logging.L(r.ctx).Error(
			"Failed to parse config",
			logging.ErrAttr(err),
		)
		return err
	}

	if r.cfg.Servers == nil {
		r.cfg.Servers = make(map[string]server.Server)
	}
	for k, v := range cfg.Servers {
		newKeySrv := k + "_" + v.GetName()
		r.cfg.Servers[newKeySrv] = v
	}

	if r.cfg.Databases == nil {
		r.cfg.Databases = make(map[string]database.Database)
	}
	for k, v := range cfg.Databases {
		if srv, exists := r.cfg.Servers[k]; !exists {
			if _, ok := cfg.Servers[v.Server]; !ok {
				newKeyDB := k + "_" + v.GetName()
				v.Server = v.Server + "_" + srv.GetName()
				r.cfg.Databases[newKeyDB] = v
			}
		}
	}

	validate := validation.New()
	if err := validate.Handler(r.cfg); err != nil {
		logging.L(r.ctx).Error(
			"Failed to validations",
			logging.ErrAttr(err),
		)
		return fmt.Errorf("config validation failed: %w", err)
	}

	data := make(map[string]dbConnect.DBConnect, len(cfg.Databases))

	for idx, db := range cfg.Databases {
		data[idx] = dbConnect.DBConnect{
			Database: db,
			Server:   cfg.Servers[db.Server],
		}
	}

	r.connectDBs = data

	return nil
}
