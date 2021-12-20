package service

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
	"github.com/zhupanovdm/go-runtime-monitor/internal/encoder"
)

var PollInterval time.Duration
var ReportInterval time.Duration
var Server string

func Start() {
	URL, err := url.Parse(Server)
	if err != nil {
		log.Fatal(fmt.Sprintf("cant parse srv parameter %s. must be correct url: %v", Server, err))
	}

	data := make(chan encoder.Encoder, 1024)
	app.Periodic(PollInterval, Publisher(data, MetricsReader()))
	app.Periodic(ReportInterval, Subscriber(data, ServerTransmitter(URL)))
	app.Serve()
}

func Publisher(data chan<- encoder.Encoder, produce func(chan<- encoder.Encoder)) app.Executor {
	return app.ExecutorHandler{
		OnExec: func(context.Context, context.CancelFunc) {
			log.Println("fetch metrics")
			produce(data)
		},
		OnEnd: func() {
			close(data)
			log.Println("publisher completed")
		},
	}
}

func Subscriber(data <-chan encoder.Encoder, consume func(val encoder.Encoder) error) app.Executor {
	return app.ExecutorHandler{
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
