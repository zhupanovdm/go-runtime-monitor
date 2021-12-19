package service

import (
	"context"
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

			resp, err := http.Post(getUrl(v), "text/plain", nil)
			if err != nil {
				log.Printf("unable to server: %v", err)
				return
			}
			if resp.StatusCode != 200 {
				log.Printf("server responded: %d", resp.StatusCode)
			}
		}
	}
	mt.OnEnd = func() {
		log.Println("metrics transporter completed")
	}
	return mt
}

func getUrl(val encoder.Encoder) string {
	u, _ := url.Parse(Url.String())
	u.Path = path.Join(u.Path, "update", val.Encode())
	return u.String()
}
