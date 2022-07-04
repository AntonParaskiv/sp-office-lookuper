package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"sp-office-lookuper/internal/logging"
)

const (
	StackTag    = "stack"
	RecoveryTag = "recovery"
)

type GracefulShutdown struct {
	sleepTime        time.Duration
	requests         *chan int
	maxRequestsCount int // healthchecker returns error when connections more than
	logger           *logging.Logger
}

func (g GracefulShutdown) GetMiddleware() func(http.Handler) http.Handler {
	gracefulMiddleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			*g.requests <- 1
			defer func() {
				<-*g.requests
				// recover after panic and return 500
				if rec := recover(); rec != nil {
					msg := "GracefullShutdownMiddleware: recover server from crash."
					trace := string(debug.Stack())

					if g.logger != nil {
						g.logger.CreateEntry().
							WithField(StackTag, trace).
							WithField(RecoveryTag, rec).
							Error(r.Context(), msg)
					} else {
						args := []interface{}{"error: " + msg, rec}
						log.Println(args...)
					}
					// return 500 after recovered from panic
					WriteError(w, r, http.StatusInternalServerError, "Internal Server Error",
						"GracefulShutdown-1")
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
	return gracefulMiddleware
}

func (g GracefulShutdown) HealthcheckConnections(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g.maxRequestsCount > 0 && len(*g.requests) > g.maxRequestsCount {
			WriteError(w, r, http.StatusTooManyRequests, "Too Many Requests", "TooManyRequests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (g GracefulShutdown) WaitUntilRequestsDone() {
	for len(*g.requests) > 0 {
		g.logger.CreateEntry().Infof(
			context.Background(),
			"waiting active connections: %d, sleep %d milliseconds",
			len(*g.requests),
			g.sleepTime,
		)
		time.Sleep(g.sleepTime * time.Millisecond)
	}
}

// NewGracefulShutdownWithConnLimiter returns new graceful shutdowner
func NewGracefulShutdownWithConnLimiter(connectionsLimit, maxConnections, sleepTime int,
	logger *logging.Logger) *GracefulShutdown {
	requests := make(chan int, connectionsLimit)
	return &GracefulShutdown{
		requests:         &requests,
		logger:           logger,
		maxRequestsCount: maxConnections,
		sleepTime:        time.Duration(sleepTime),
	}
}
