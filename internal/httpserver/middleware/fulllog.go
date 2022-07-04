package middleware

import (
	"context"
	"net/http"
)

type ShowFullLogCtxKeyType string

const ShowFullLogCtxKey = ShowFullLogCtxKeyType("show-full-log")

// ShowFullLog handler enable full log
func ShowFullLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ShowFullLogCtxKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
