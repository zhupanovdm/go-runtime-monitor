package monitor

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

const serverName = "Monitor HTTP Server"

type Server struct {
	*http.Server
	wg sync.WaitGroup
}

func (srv *Server) Start(ctx context.Context) {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(serverName))
	logger.Info().Msgf("running server on %v", srv.Addr)

	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		if err := srv.ListenAndServe(); err != nil {
			logger.Err(err).Msg("server stopped")
		}
	}()
}

func (srv *Server) Stop(ctx context.Context) {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(serverName))
	if err := srv.Close(); err != nil {
		logger.Err(err).Msg("server close failed")
	}
	srv.wg.Wait()
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: root(handler, CID, middleware.RealIP, Logger, middleware.Recoverer),
	}
	return &Server{Server: srv}
}

func root(h http.Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mw := range middlewares {
		router.Use(mw)
	}
	router.Mount("/", h)
	return router
}

func CID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get(logging.CorrelationIDHeader)
		if cid == "" {
			cid = logging.NewCID()
		}
		ctx, _ := logging.SetCID(r.Context(), cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := logging.SetIfAbsentCID(r.Context(), logging.NewCID())
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(serverName), logging.WithCID(ctx))
		logger.Info().Msgf("%s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
