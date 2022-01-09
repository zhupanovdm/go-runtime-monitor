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
	StoreInterval  time.Duration `env:"STORE_INTERVAL"`
	StoreFile      string        `env:"STORE_FILE"`
	Restore        bool          `env:"RESTORE"`
}

func New() *Config {
	return &Config{
		Address:        "localhost:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		ReportBuffer:   1024,
		StoreInterval:  10 * time.Second,
		StoreFile:      "/tmp/devops-metrics-db.json",
		Restore:        true,
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
	flag.DurationVar(&c.ReportInterval, "store-interval", c.ReportInterval, "Monitor store interval")
	flag.DurationVar(&c.ReportInterval, "store-file", c.ReportInterval, "Monitor store file")
	flag.DurationVar(&c.ReportInterval, "restore", c.ReportInterval, "Monitor will restore metrics at startup")

	flag.DurationVar(&c.PollInterval, "poll-interval", c.PollInterval, "Agent polling interval")
	flag.DurationVar(&c.ReportInterval, "report-interval", c.ReportInterval, "Agent reporting interval")

	return nil
}
