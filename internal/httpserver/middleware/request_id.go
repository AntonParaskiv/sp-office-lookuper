package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	RequestIDHeader = "X-Request-Id"
	RequestIDKey    = "rid"
)

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID(next http.Handler) http.Handler {
	nextWithContext := RequestIDToOutgoingContext(next)
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Add(RequestIDHeader, requestID)
		}
		w.Header().Set(RequestIDHeader, requestID)
		//nolint:revive,staticcheck // it's ok
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		nextWithContext.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func RequestIDToOutgoingContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rid, ok := ctx.Value(RequestIDKey).(string)
		if !ok {
			next.ServeHTTP(w, r)
		}
		ctx = metadata.AppendToOutgoingContext(ctx, RequestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
