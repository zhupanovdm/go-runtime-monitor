package main

import (
	"flag"

	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/handlers"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/repo"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/service"
)

func main() {
	flag.IntVar(&service.Port, "p", 8080, "Server port")

	router := handlers.NewRouter()
	metrics := service.NewMetrics(repo.Gauges(), repo.Counters())
	router.Mount("/", handlers.NewMetricsHandler(metrics))

	service.StartMonitor(router)
}
