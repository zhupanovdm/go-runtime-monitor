package main

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/handlers"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/storage/file"
	"github.com/zhupanovdm/go-runtime-monitor/storage/trivial"
)

func cli(cfg *config.Config, flag *flag.FlagSet) {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Monitor server address")
	flag.BoolVar(&cfg.Restore, "r", true, "Monitor will restore metrics at startup")
	flag.DurationVar(&cfg.StoreInterval, "i", 300*time.Second, "Monitor store interval")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "Monitor store file")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName("Monitor app"))
	logger.Info().Msg("starting runtime metrics monitor server")

	cfg, err := config.Load(cli)
	if err != nil {
		logger.Err(err).Msg("failed to load server config")
		return
	}

	mon := monitor.NewMonitor(cfg, file.NewStorage(cfg), trivial.NewGaugeStorage(), trivial.NewCounterStorage())
	if err := mon.Restore(ctx); err != nil {
		logger.Err(err).Msg("failed to restore metrics")
	}

	var wg sync.WaitGroup
	go mon.BackgroundTask().With(task.CompletionWait(&wg))(ctx)

	root := handlers.NewMetricsRouter(handlers.NewMetricsHandler(mon), handlers.NewMetricsAPIHandler(mon))
	server := monitor.NewServer(cfg, root)
	server.Start(ctx)

	logger.Info().Msgf("%v signal received", <-app.TerminationSignal())
	server.Stop(ctx)
	cancel()
	wg.Wait()
}
