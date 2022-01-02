package monitor

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
)

type Server struct {
	*http.Server
	*sync.WaitGroup
}

func RunServer(cfg *config.Config, handler http.Handler) *Server {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: httplib.NewRouter(handler,
			middleware.RequestID,
			middleware.RealIP,
			middleware.Logger,
			middleware.Recoverer),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil {
		}
	}()

	return &Server{
		Server:    srv,
		WaitGroup: &wg,
	}
}
