package main

import (
	"flag"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor/handlers"
	"github.com/zhupanovdm/go-runtime-monitor/storage/trivial"
)

func main() {
	flags := flag.NewFlagSet("monitor", flag.ExitOnError)
	cfg := config.New().FromCLI(flags)

	mon := monitor.NewMetricsMonitor(trivial.NewGaugeStorage(), trivial.NewCounterStorage())
	server := monitor.RunServer(cfg, handlers.NewMetricsHandler(mon))

	<-app.TerminationSignal()
	if err := server.Close(); err != nil {
	}
	server.Wait()

}
