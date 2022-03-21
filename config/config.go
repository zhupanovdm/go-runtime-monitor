// Package config is used to describe application runtime operation parameters.
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

type (
	// Config describes application parameters.
	Config struct {
		// Address is Monitor server address.
		Address string `env:"ADDRESS"`

		// PollInterval specifies runtime metrics polling period.
		PollInterval time.Duration `env:"POLL_INTERVAL"`

		// ReportInterval specifies collected metrics send period.
		ReportInterval time.Duration `env:"REPORT_INTERVAL"`

		// StoreInterval specifies dumping period. Dumps on every update if not set.
		StoreInterval time.Duration `env:"STORE_INTERVAL"`

		// StoreFile sets file to dump gathered metrics.
		StoreFile string `env:"STORE_FILE"`

		// Restore enables metrics restore from dump on monitor server startup. Disables dumping if not set.
		Restore bool `env:"RESTORE"`

		// Key is a secret key for signing metrics that will be transmitted to monitor server.
		Key string `env:"KEY"`

		// Database describes database connection which will be used to persist gathered metrics.
		Database string `env:"DATABASE_DSN"`

		// PProfAddress is address for pprof utility
		PProfAddress string
	}

	// CLIExport is used to expose CLI options to Config
	CLIExport func(*Config, *flag.FlagSet)
)

// Load returns initialized config. Config parameters wil be searched in such priority:
// 1. environment
// 2. CLI (if CLIExport is specified)
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
