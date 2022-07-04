package httpserver

import (
	"context"
	"net/http"

	"sp-office-lookuper/internal/config"
	"sp-office-lookuper/internal/httpserver/middleware"
	"sp-office-lookuper/internal/logging"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// operations
	OperationOffice = "/office"
)

type HTTPServer struct {
	cfg        *config.APIConfig
	server     *http.Server
	storage    Storage
	logger     *logging.Logger
	shutDowner *middleware.GracefulShutdown
}

func (s *HTTPServer) prepareHandlers(ctx context.Context) http.Handler {
	handlers := NewHTTPHandlers(ctx, s.cfg, s.storage, s.logger)

	r := chi.NewRouter()

	r.Use(
		middleware.ShowFullLog,
		chiMiddleware.NoCache,
		middleware.RequestID,
		s.shutDowner.GetMiddleware(),
	)

	r.Group(func(r chi.Router) {
		metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			promhttp.Handler().ServeHTTP(w, r)
		})

		r.Get("/v1/metrics", metricsHandler)
		r.Get("/healthcheck", handlers.Healthcheck)

		r.Post("/office", handlers.OfficeHandler)
	})

	return r
}

func (s *HTTPServer) Close(ctx context.Context) error {
	// waiting until all requests are done
	s.shutDowner.WaitUntilRequestsDone()
	err := s.server.Close()
	if err != nil {
		return err
	}

	s.logger.Info(ctx, "HTTP server shutting down...")
	return nil
}

func (s *HTTPServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func NewHTTPServer(ctx context.Context,
	cfg *config.APIConfig,
	storage Storage,
	logger *logging.Logger,
) (*HTTPServer, error) {
	srv := &HTTPServer{
		cfg:     cfg,
		storage: storage,
		logger:  logger,
	}

	srv.shutDowner = middleware.NewGracefulShutdownWithConnLimiter(
		cfg.HTTPMaxConnection, cfg.HTTPMaxConnection-1, 100, logger)

	srv.server = &http.Server{
		Addr:    srv.cfg.HTTPListenAddress,
		Handler: srv.prepareHandlers(ctx),
	}

	return srv, nil
}
