package main

import (
	"context"
	"flag"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/handlers"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
	"github.com/zhupanovdm/go-runtime-monitor/storage/file"
	"github.com/zhupanovdm/go-runtime-monitor/storage/sqldb"
	"github.com/zhupanovdm/go-runtime-monitor/storage/trivial"
)

func cli(cfg *config.Config, flag *flag.FlagSet) {
	flag.StringVar(&cfg.Address, "a", config.DefaultAddress, "Monitor server address")
	flag.BoolVar(&cfg.Restore, "r", config.DefaultRestore, "Monitor will restore metrics at startup")
	flag.DurationVar(&cfg.StoreInterval, "i", config.DefaultStoreInterval, "Monitor store interval")
	flag.StringVar(&cfg.StoreFile, "f", config.DefaultStoreFile, "Monitor store file")
	flag.StringVar(&cfg.Key, "k", "", "Packet signing key")
	flag.StringVar(&cfg.Database, "d", "", "Database connection string")
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

	dumper := storage.New(cfg, sqldb.New(sqldb.PGX{}), file.New)
	if dumper != nil {
		if err := dumper.Init(ctx); err != nil {
			logger.Err(err).Msg("failed to init storage")
			return
		}
	}

	mon := monitor.NewMonitor(cfg, dumper, trivial.NewGaugeStorage(), trivial.NewCounterStorage())
	if err := mon.Restore(ctx); err != nil {
		logger.Err(err).Msg("failed to restore metrics")
	}

	var wg sync.WaitGroup
	go mon.BackgroundTask().With(task.CompletionWait(&wg))(ctx)

	root := handlers.NewMetricsRouter(handlers.NewMetricsHandler(mon), handlers.NewMetricsAPIHandler(cfg, mon))
	server := monitor.NewServer(cfg, root)
	server.Start(ctx)

	logger.Info().Msgf("%v signal received", <-app.TerminationSignal())

	server.Stop(ctx)
	dumper.Close(ctx)

	cancel()
	wg.Wait()
}
