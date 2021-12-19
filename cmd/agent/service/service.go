package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
	"github.com/zhupanovdm/go-runtime-monitor/internal/encoder"
	"github.com/zhupanovdm/go-runtime-monitor/internal/metrics"
)

func MetricsCollector(data chan<- encoder.Encoder) app.Executor {
	var reader func(chan<- encoder.Encoder)

	mc := app.ExecutorHandler{}
	mc.OnStart = func() {
		reader = metrics.Reader()
	}
	mc.OnExec = func(ctx context.Context, cancel context.CancelFunc) {
		log.Println("metrics fetch")
		reader(data)
	}
	mc.OnEnd = func() {
		close(data)
		log.Println("metrics reader completed")
	}
	return mc
}

func MetricsTransporter(data <-chan encoder.Encoder) app.Executor {
	mt := app.ExecutorHandler{}
	mt.OnExec = func(ctx context.Context, cancel context.CancelFunc) {
		log.Println("metrics transit")
		count := len(data)
		for i := 0; i < count; i++ {
			v, ok := <-data
			if !ok {
				log.Println("metrics pipe closed")
				return
			}
			if err := sendToServer(v); err != nil {
				log.Printf("error occured on server send: %v", err)
				return
			}
		}
	}
	mt.OnEnd = func() {
		log.Println("metrics transporter completed")
	}
	return mt
}

func sendToServer(val encoder.Encoder) error {
	resp, err := http.Post(getURL(val), "text/plain", nil)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to connect to server: %v", err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("unable to close response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("server responded: %d", resp.StatusCode))
	}

	return nil
}

func getURL(val encoder.Encoder) string {
	u, _ := url.Parse(URL.String())
	u.Path = path.Join(u.Path, "update", val.Encode())
	return u.String()
}
