package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/handlers"
	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
	"log"
	"net/http"
)

var port = flag.Int("p", 8080, "Server port")

func main() {
	app.Once(server(*port, handlers.NewRouter()))
	app.Serve()
}

func server(port int, handler http.Handler) app.ExecutorHandler {
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
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
