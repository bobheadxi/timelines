package log

import (
	"context"
	"time"

	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const (
	// LogKeyRID is the key to use for request IDs in logging
	LogKeyRID = "request_id"
)

type loggerMiddleware struct {
	l *zap.Logger
}

// NewHTTPLogger instantiates a middleware function that logs all requests
// using the provided logger
func NewHTTPLogger(l *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return loggerMiddleware{l.Desugar()}.Handler
}

func (l loggerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		l.l.Info(r.Method+" "+r.URL.Path+": request completed",
			// request metadata
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("method", r.Method),
			zap.String("user_agent", r.UserAgent()),

			// response metadata
			zap.Int("status", ww.Status()),
			zap.Duration("took", time.Since(start)),

			// additional metadata
			zap.String("real_ip", r.RemoteAddr),
			zap.String(LogKeyRID, RequestID(r.Context())))
	})
}

// RequestID returns the request ID injected by the chi requestID middleware
func RequestID(ctx context.Context) string {
	if reqID := ctx.Value(middleware.RequestIDKey); reqID != nil {
		return reqID.(string)
	}
	return ""
}
