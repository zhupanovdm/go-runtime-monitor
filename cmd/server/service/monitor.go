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
		log.Println("starting server at", server.Addr)
		log.Println(server.ListenAndServe())
	}).Serve(func() {
		log.Println("closing server")
		if err := server.Close(); err != nil {
			log.Fatal(err)
		}
	})
	a.Immediate()
	log.Println("monitor stopped")
}
