package monitor

import (
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
)

const DefaultTimeout = 30 * time.Second

type Config struct {
	*config.Config
	Timeout time.Duration
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Config:  cfg,
		Timeout: DefaultTimeout,
	}
}
