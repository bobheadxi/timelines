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
			zap.String("req.path", r.URL.Path),
			zap.String("req.query", r.URL.RawQuery),
			zap.String("req.method", r.Method),
			zap.String("req.ip", r.RemoteAddr),
			zap.String("req.user_agent", r.UserAgent()),

			// response metadata
			zap.Int("resp.status", ww.Status()),

			// additional metadata
			zap.Duration("duration", time.Since(start)),
			zap.String(LogKeyRID, HTTPRequestID(r.Context())))
	})
}

// HTTPRequestID returns the request ID injected by the chi requestID middleware
func HTTPRequestID(ctx context.Context) string { return ctxString(ctx, middleware.RequestIDKey) }
