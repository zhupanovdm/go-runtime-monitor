package monitor

import (
	"compress/gzip"
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
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
		Addr: cfg.Address,
		Handler: entryHandler(handler,
			middleware.RealIP,
			cid,
			serverLogger,
			compress,
			decompress,
			middleware.Recoverer),
	}
	return &Server{Server: srv}
}

func entryHandler(h http.Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mw := range middlewares {
		router.Use(mw)
	}
	router.Mount("/", h)
	return router
}

func decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			handleInternalError(w, r, err, "decompressor: failed to create")
			return
		}
		defer func() {
			if err := gz.Close(); err != nil {
				handleInternalError(w, r, err, "decompressor: failed to close")
			}
		}()
		r.Body = gz
		next.ServeHTTP(w, r)
	})
}

func compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			handleInternalError(w, r, err, "compressor: failed to create")
			return
		}
		defer func() {
			if err := gz.Close(); err != nil {
				handleInternalError(w, r, err, "compressor: failed to close")
			}
		}()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(httplib.ResponseCustomWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func cid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get(logging.CorrelationIDHeader)
		if cid == "" {
			cid = logging.NewCID()
		}
		ctx, _ := logging.SetCID(r.Context(), cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func serverLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(serverName), logging.WithCID(ctx))

		logger = logger.With().
			Stringer("header", httplib.Header(r.Header)).
			Str("remote_addr", r.RemoteAddr).
			Logger()

		logger.Info().Msgf("%s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func handleInternalError(w http.ResponseWriter, r *http.Request, err error, msg string) {
	ctx := r.Context()
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(serverName), logging.WithCID(ctx))
	logger.Err(err).Msg(msg)
	httplib.Error(w, http.StatusInternalServerError, nil)
}
