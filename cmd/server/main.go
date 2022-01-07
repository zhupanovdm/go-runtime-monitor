package main

import (
	"context"
	"flag"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/handlers"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/storage/trivial"
)

func main() {
	ctx := context.Background()
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName("Monitor app"))
	logger.Info().Msg("starting runtime metrics monitor server")

	cfg := config.New()
	if err := cfg.LoadFromEnv(); err != nil {
		logger.Err(err).Msg("failed to load app config")
	}
	if err := cfg.FromCLI(flag.NewFlagSet("monitor", flag.ExitOnError)); err != nil {
		logger.Err(err).Msg("failed to load app config")
	}

	mon := monitor.NewMonitor(trivial.NewGaugeStorage(), trivial.NewCounterStorage())
	root := handlers.NewMetricsRouter(handlers.NewMetricsHandler(mon), handlers.NewMetricsApiHandler(mon))
	server := monitor.NewServer(cfg, root)
	server.Start(ctx)

	logger.Info().Msgf("%v signal received", <-app.TerminationSignal())
	server.Stop(ctx)
}
