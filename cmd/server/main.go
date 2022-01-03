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
	logger.Info().Msg("runtime metrics monitor server starting")

	flags := flag.NewFlagSet("monitor", flag.ExitOnError)
	cfg := config.New().FromCLI(flags)

	mon := monitor.NewMetricsMonitor(trivial.NewGaugeStorage(), trivial.NewCounterStorage())
	server := monitor.NewServer(cfg, handlers.NewMetricsHandler(mon))
	server.Start(ctx)

	logger.Info().Msgf("%v signal received", <-app.TerminationSignal())
	server.Stop(ctx)
}
