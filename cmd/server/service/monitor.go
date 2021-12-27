package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
)

var Port int

func StartMonitor(handler http.Handler) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: handler,
	}

	a := app.NewApp()
	a.NewTask(func(t *app.Task) {
		log.Printf("starting server at %v", server.Addr)
		log.Fatal(server.ListenAndServe())
	}).Serve(func() {
		log.Println("closing server")
		if err := server.Close(); err != nil {
			log.Fatal(err)
		}
	})
	a.Immediate()
	log.Println("monitor stopped")
}
