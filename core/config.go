package core

import (
	"context"
	"errors"

	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server *server `json:"server,omitempty"`
	Redis  *redis  `json:"redis,omitempty"`
}

type server struct {
	Listen string `json:"listen"`
}

type redis struct {
	Addr     string `json:"addr"`
	Protocol string `json:"protocol"`
	Database int    `json:"database"`
	PoolSize int    `json:"pool_size"`
}

func NewConfig(ctx context.Context, path string) (*Config, error) {
	conf := &Config{}
	err := configor.Load(conf, path)
	if err != nil {
		log.Errorf("load config err %+v", err)
		return nil, err
	}
	err = conf.check()
	if err != nil {
		log.Errorf("check config err %+v", err)
		return nil, err
	}
	return conf, nil
}

func (c *Config) check() error {
	if c.Server == nil {
		return errors.New("invalid server config")
	}

	if c.Redis == nil {
		return errors.New("invalid redis config")
	}
	return nil
}
