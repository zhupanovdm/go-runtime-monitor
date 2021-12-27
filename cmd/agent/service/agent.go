package service

import (
	"log"
	"math/rand"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

var PollInterval time.Duration
var ReportInterval time.Duration
var ServerURL string

func StartAgent() {
	rand.Seed(time.Now().UnixNano())

	pipe := make(chan *metric.Metric, 1024)

	var pollCounter int64
	client := monitorClient(ServerURL)

	a := app.NewApp()
	a.NewTask(func(t *app.Task) {
		log.Println("fetch metrics")
		pollCounter++
		pollRuntimeMetrics(pipe, pollCounter)
	}).Periodic(PollInterval, func() {
		close(pipe)
		log.Println("poller stopped")
	})

	a.NewTask(func(t *app.Task) {
		log.Println("send to remote")

		count := len(pipe)
		for count > 0 {
			value, ok := <-pipe
			if !ok {
				log.Println("data pipe closed")
				return
			}
			if err := sendToMonitorServer(client, value.String()); err != nil {
				log.Printf("error occured while transmiting to server: %v", err)
				return
			}
			count--
		}
	}).Periodic(ReportInterval, func() {
		log.Println("sender stopped")
	})
	a.Immediate()
	log.Println("agent stopped")
}
