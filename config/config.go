package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	ReportBuffer   int
}

func New() *Config {
	return &Config{
		Address:        "localhost:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		ReportBuffer:   1024,
	}
}

func (c *Config) LoadFromEnv() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) FromCLI(flag *flag.FlagSet) error {
	flag.StringVar(&c.Address, "address", c.Address, "Monitor server address")
	flag.DurationVar(&c.PollInterval, "poll-interval", c.PollInterval, "Metrics polling interval")
	flag.DurationVar(&c.ReportInterval, "report-interval", c.ReportInterval, "Metrics reporting interval")
	return nil
}
