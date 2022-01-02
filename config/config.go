package config

import (
	"flag"
	"time"
)

type Config struct {
	ServerPort         int
	PollInterval       time.Duration
	ReportInterval     time.Duration
	ReporterBufferSize int
}

func New() *Config {
	return &Config{}
}

func (c *Config) FromCLI(flag *flag.FlagSet) *Config {
	flag.IntVar(&c.ServerPort, "server-port", 8080, "Monitor server port")

	flag.DurationVar(&c.PollInterval, "poll-interval", 2*time.Second, "Metrics polling interval")
	flag.DurationVar(&c.ReportInterval, "report-interval", 10*time.Second, "Metrics reporting interval")
	flag.IntVar(&c.ReporterBufferSize, "reporter-buffer", 1024, "Reporter buffer size")
	return c
}
