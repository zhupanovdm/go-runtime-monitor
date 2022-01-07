package main

import (
	"context"
	"flag"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
	client "github.com/zhupanovdm/go-runtime-monitor/providers/monitor/http/v2"
	"github.com/zhupanovdm/go-runtime-monitor/service/agent"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName("Agent app"))
	logger.Info().Msg("starting runtime metrics monitor agent")

	cfg := config.New()
	if err := cfg.LoadFromEnv(); err != nil {
		logger.Err(err).Msg("failed to load app config")
	}
	if err := cfg.FromCLI(flag.NewFlagSet("agent", flag.ExitOnError)); err != nil {
		logger.Err(err).Msg("failed to load app config")
	}

	mon := client.NewClient(monitor.NewConfig(cfg))
	reporterSvc := agent.NewMetricsReporter(cfg, mon)
	collectorSvc := agent.NewRuntimeMetricsCollector(cfg, reporterSvc)

	var wg sync.WaitGroup
	go reporterSvc.BackgroundTask().With(task.CompletionWait(&wg))(ctx)
	go collectorSvc.BackgroundTask().With(task.CompletionWait(&wg))(ctx)

	logger.Info().Msgf("%v signal received", <-app.TerminationSignal())
	cancel()
	wg.Wait()
}
