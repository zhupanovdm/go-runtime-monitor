package main

import (
	"context"
	"flag"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/http"
	"github.com/zhupanovdm/go-runtime-monitor/service/agent"
)

func main() {
	var wg sync.WaitGroup

	flags := flag.NewFlagSet("agent", flag.ExitOnError)
	cfg := config.New().FromCLI(flags)

	client := http.NewClient(http.NewConfig().FromCLI(flags))
	reporterSvc := agent.NewMetricsReporter(cfg, client)
	collectorSvc := agent.NewRuntimeMetricsCollector(cfg, reporterSvc)

	ctx, cancel := context.WithCancel(context.Background())
	go reporterSvc.BackgroundTask().With(task.CompletionWait(&wg))(ctx)
	go collectorSvc.BackgroundTask().With(task.CompletionWait(&wg))(ctx)

	<-app.TerminationSignal()
	cancel()
	wg.Wait()
}
