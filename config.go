package robot

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"

	"github.com/supergiant/robot/pkg/api"
	"github.com/supergiant/robot/pkg/logger"
	"github.com/supergiant/robot/pkg/plugin"
)

// Config  struct represents configuration of robot service
type Config struct {
	Logging logger.Config   `mapstructure:"logging"`
	API     api.Config      `mapstructure:"api"`
	Plugin  plugin.Config   `mapstructure:"plugin"`
	ETCD    clientv3.Config `mapstructure:"etcd"`
}

// Validate checks configuration instance for correctness
func (c *Config) Validate() error {
	if err := c.Logging.Validate(); err != nil {
		return err
	}

	if err := c.Plugin.Validate(); err != nil {
		return err
	}

	if len(c.ETCD.Endpoints) == 0 {
		return errors.New("etcd endpoints where not configured")
	}

	return nil
}
