package main

import (
	"context"
	"flag"
	"net/http"
	"sync"

	_ "net/http/pprof"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
	client "github.com/zhupanovdm/go-runtime-monitor/providers/monitor/http/v2"
	"github.com/zhupanovdm/go-runtime-monitor/service/agent"
)

var (
	buildVersion app.BuildInfoString
	buildDate    app.BuildInfoString
	buildCommit  app.BuildInfoString
)

func cli(cfg *config.Config, flag *flag.FlagSet) {
	flag.StringVar(&cfg.Address, "a", config.DefaultAddress, "Monitor server address")
	flag.DurationVar(&cfg.ReportInterval, "r", config.DefaultReportInterval, "Agent reporting interval")
	flag.DurationVar(&cfg.PollInterval, "p", config.DefaultPollInterval, "Agent polling interval")
	flag.StringVar(&cfg.Key, "k", "", "Packet signing key")
}

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName("Agent app"))

	logger.Info().Msgf("Build version: %v", buildVersion)
	logger.Info().Msgf("Build date: %v", buildDate)
	logger.Info().Msgf("Build commit: %v", buildCommit)

	logger.Info().Msg("starting runtime metrics monitor agent")

	cfg, err := config.Load(cli)
	if err != nil {
		logger.Err(err).Msg("failed to load agent config")
		return
	}

	mon, err := client.NewClient(monitor.NewConfig(cfg))
	if err != nil {
		logger.Err(err).Msg("failed to create monitor client")
		return
	}

	froze := agent.NewFroze()
	reporterSvc := agent.NewMetricsReporter(cfg, froze, mon)
	collector := agent.NewMetricsCollector(cfg, froze, agent.MemStats(), agent.PS())

	go reporterSvc.BackgroundTask().With(task.CompletionWait(&wg))(ctx)
	go collector.BackgroundTask().With(task.CompletionWait(&wg))(ctx)

	srv := &http.Server{Addr: cfg.PProfAddress}
	defer func() {
		logger.Info().Msg("closing pprof server")
		if err := srv.Close(); err != nil {
			logger.Err(err).Msg("failed to stop pprof server")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info().Msg("starting pprof server")
		if err := srv.ListenAndServe(); err != nil {
			logger.Err(err).Msg("stopped to serve pprof")
		}
	}()

	logger.Info().Msgf("%v signal received", <-app.TerminationSignal())
}
