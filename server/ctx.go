package server

import (
	"context"

	"github.com/go-chi/chi/middleware"
)

func requestID(ctx context.Context) string {
	return ctxString(ctx, middleware.RequestIDKey)
}

// ctxString gets a string with the given key from the given context's values
func ctxString(ctx context.Context, key interface{}) string {
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	return ""
}
