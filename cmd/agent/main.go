package main

import (
	"flag"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/cmd/agent/service"
)

func main() {
	flag.DurationVar(&service.PollInterval, "pi", 2*time.Second, "Poll interval")
	flag.DurationVar(&service.ReportInterval, "ri", 10*time.Second, "Poll interval")
	flag.StringVar(&service.ServerURL, "srv", "http://127.0.0.1:8080", "Base url of agent server")

	service.StartAgent()
}
