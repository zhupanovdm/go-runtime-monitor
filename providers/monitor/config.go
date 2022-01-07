package monitor

import (
	"flag"
	"time"
)

type Config struct {
	Server  string
	Timeout time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) FromCLI(flag *flag.FlagSet) *Config {
	flag.StringVar(&c.Server, "monitor-srv", "http://localhost:8080", "monitor server URL")
	flag.DurationVar(&c.Timeout, "client-timeout", 30*time.Second, "client timeout")
	return c
}
