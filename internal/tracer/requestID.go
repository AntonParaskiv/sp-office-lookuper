package tracer

import (
	"context"

	"github.com/google/uuid"
)

func GetOrCreateRequestID(ctx context.Context) string {
	rid := getRequestID(ctx)
	if rid == "" {
		rid = uuid.New().String()
	}
	return rid
}

func getRequestID(ctx context.Context) string {
	if rid, ok := ctx.Value(RequestIDKey).(string); ok {
		return rid
	}
	return ""
}
