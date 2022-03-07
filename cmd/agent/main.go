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

	reporterSvc := agent.NewMetricsReporter(cfg, mon)
	collector := agent.NewMetricsCollector(cfg, reporterSvc, agent.MemStats(), agent.PS())

	go reporterSvc.BackgroundTask().With(task.CompletionWait(&wg))(ctx)
	go collector.BackgroundTask().With(task.CompletionWait(&wg))(ctx)

	logger.Info().Msgf("%v signal received", <-app.TerminationSignal())
}
