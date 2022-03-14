package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	DefaultAddress        = "localhost:8080"
	DefaultReportInterval = 10 * time.Second
	DefaultPollInterval   = 2 * time.Second
	DefaultRestore        = true
	DefaultStoreInterval  = 300 * time.Second
	DefaultStoreFile      = "/tmp/devops-metrics-db.json"
	DefaultPProfAddress   = ":9000"
)

type CLIExport func(*Config, *flag.FlagSet)

type Config struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL"`
	StoreFile      string        `env:"STORE_FILE"`
	Restore        bool          `env:"RESTORE"`
	Key            string        `env:"KEY"`
	Database       string        `env:"DATABASE_DSN"`
	PProfAddress   string
}

func Load(cli CLIExport) (*Config, error) {
	cfg := &Config{PProfAddress: DefaultPProfAddress}

	if cli != nil {
		cli(cfg, flag.CommandLine)
		flag.Parse()
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
