package service

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

var PollInterval time.Duration
var ReportInterval time.Duration
var ServerURL string

func Start() {
	rand.Seed(time.Now().UnixNano())

	data := make(chan *metric.Metric, 1024)

	var pollCounter int64
	app.Periodic(PollInterval, publish(data, func(pipe chan<- *metric.Metric) {
		pollCounter++
		pollRuntimeMetrics(pipe, pollCounter)
	}))

	client := monitorClient(ServerURL)
	app.Periodic(ReportInterval, subscribe(data, func(val *metric.Metric) error {
		return sendToMonitorServer(client, val.String())
	}))

	app.Serve()
}

func publish(data chan<- *metric.Metric, produce func(chan<- *metric.Metric)) app.Executor {
	return app.ExecutorHandler{
		OnStart: func() {
			log.Println("poller started")
		},
		OnExec: func(context.Context, context.CancelFunc) {
			log.Println("fetch metrics")
			produce(data)
		},
		OnEnd: func() {
			close(data)
			log.Println("poller completed")
		},
	}
}

func subscribe(data <-chan *metric.Metric, consume func(*metric.Metric) error) app.Executor {
	return app.ExecutorHandler{
		OnStart: func() {
			log.Println("send started")
		},
		OnExec: func(context.Context, context.CancelFunc) {
			log.Println("send to remote")

			count := len(data)
			for count > 0 {
				value, ok := <-data
				if !ok {
					log.Println("data pipe closed")
					return
				}
				if err := consume(value); err != nil {
					log.Printf("error occured while transmiting to server: %v", err)
					return
				}
				count--
			}
		},
		OnEnd: func() {
			log.Println("send complete")
		},
	}
}
