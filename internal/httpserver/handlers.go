package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"sp-office-lookuper/internal/config"
	"sp-office-lookuper/internal/logging"
)

const (
	HealthCheckResponse = "sp-office-lookuper is alive"
)

type HTTPHandlers struct {
	cfg     *config.APIConfig
	storage Storage
	logger  *logging.Logger
}

func (h *HTTPHandlers) Healthcheck(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, HealthCheckResponse)
}

func NewHTTPHandlers(_ context.Context, cfg *config.APIConfig, storage Storage, logger *logging.Logger) *HTTPHandlers {
	return &HTTPHandlers{
		cfg:     cfg,
		storage: storage,
		logger:  logger,
	}
}
