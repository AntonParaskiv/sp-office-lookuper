package httpserver

import (
	"context"
	"net/http"
	"net/http/pprof"

	"sp-office-lookuper/internal/config"
	"sp-office-lookuper/internal/logging"

	"github.com/go-chi/chi"
)

type PprofServer struct {
	server *http.Server
	logger *logging.Logger
}

func (s *PprofServer) prepareHandlers() http.Handler {
	router := chi.NewRouter()

	router.Get("/debug/pprof/", pprof.Index)
	router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	router.Get("/debug/pprof/profile", pprof.Profile)
	router.Get("/debug/pprof/symbol", pprof.Symbol)
	router.Get("/debug/pprof/trace", pprof.Trace)
	router.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	router.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	router.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	router.Get("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	router.Get("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)

	return router
}

func (s *PprofServer) Close(ctx context.Context) error {
	s.logger.Info(ctx, "HTTP pprof server shutting down...")
	return s.server.Close()
}

func (s *PprofServer) ListenAndServe() (err error) {
	return s.server.ListenAndServe()
}

func NewHTTPPPprofServer(cfg *config.APIConfig, logger *logging.Logger) (*PprofServer, error) {
	srv := &PprofServer{
		logger: logger,
	}

	srv.server = &http.Server{
		Addr:    cfg.HTTPPprofAddress,
		Handler: srv.prepareHandlers(),
	}

	return srv, nil
}
