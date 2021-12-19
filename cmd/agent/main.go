package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/cmd/agent/service"
)

func main() {
	var baseURL string
	var err error

	flag.DurationVar(&service.PollInterval, "pi", 2*time.Second, "Poll interval")
	flag.DurationVar(&service.ReportInterval, "ri", 10*time.Second, "Poll interval")
	flag.StringVar(&baseURL, "srv", "http://127.0.0.1:8080", "Base url of agent server")

	service.URL, err = url.Parse(baseURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("cant parse srv parameter %s. must be correct url: %v", baseURL, err))
	}

	service.StartAgent()
	log.Println("agent completed")
}
