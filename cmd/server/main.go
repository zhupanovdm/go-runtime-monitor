package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/handlers"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/repo"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/service"
	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
	"log"
	"net/http"
)

var port = flag.Int("p", 8080, "Server port")

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	metrics := service.NewMetrics(repo.Gauges(), repo.Counters())
	router.Mount("/", handlers.NewMetricsHandler(metrics))

	app.Once(server(*port, router))
	app.Serve()
}

func server(port int, handler http.Handler) app.ExecutorHandler {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	return app.ExecutorHandler{
		OnStart: func() {
			log.Printf("starting server at %v", server.Addr)
		},
		OnExec: func(ctx context.Context, cancel context.CancelFunc) {
			log.Fatal(server.ListenAndServe())
		},
		OnEnd: func() {
			log.Println("closing server")
			if err := server.Close(); err != nil {
				log.Fatal(err)
			}
		},
	}
}
