package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
)

var PollInterval time.Duration
var ReportInterval time.Duration
var Server string

func Start() {
	URL, err := url.Parse(Server)
	if err != nil {
		log.Fatal(fmt.Sprintf("cant parse srv parameter %s. must be correct url: %v", Server, err))
	}

	data := make(chan string, 1024)

	rand.Seed(time.Now().UnixNano())

	var pollCounter int64
	app.Periodic(PollInterval, publish(data, func(pipe chan<- string) {
		pollCounter++
		pollRuntimeMetrics(pipe, pollCounter)
	}))

	app.Periodic(ReportInterval, subscribe(data, func(val string) error {
		return sendToMonitorServer(URL, val)
	}))

	app.Serve()
}

func publish(data chan<- string, produce func(chan<- string)) app.Executor {
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

func subscribe(data <-chan string, consume func(string) error) app.Executor {
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
