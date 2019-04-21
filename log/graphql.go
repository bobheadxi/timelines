package log

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"
)

// httpCtx are context keys used for injected HTTP variables, mostly for the
// convenience of the GraphQL logger
type httpCtx int

const (
	httpCtxKeyUserAgent httpCtx = iota + 1
	httpCtxKeyRemoteAddr
)

// GraphCtxHandler injects request fields into context for use with the GraphQL
// request logger
func GraphCtxHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), httpCtxKeyUserAgent, r.UserAgent())
		ctx = context.WithValue(ctx, httpCtxKeyRemoteAddr, r.RemoteAddr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewGraphLogger returns a logger for use with GraphQL queries
func NewGraphLogger(l *zap.Logger) graphql.RequestMiddleware {
	return func(ctx context.Context, next func(context.Context) []byte) []byte {
		// call handler
		start := time.Now()
		response := next(ctx)

		// log request
		// TODO: could implement more advanced tracing, hm, via RequestContext.Trace
		req := graphql.GetRequestContext(ctx)
		// TODO: log message not very informative
		// TODO: evaluate usefulness of logged fields
		l.Info("graph query completed",
			// request metadata
			zap.Int("req.complexity", req.OperationComplexity),
			zap.Any("req.variables", req.Variables),
			zap.Any("req.errors", req.Errors),
			zap.Any("req.extensions", req.Extensions),
			zap.String("req.ip", ctxString(ctx, httpCtxKeyRemoteAddr)),
			zap.String("req.user_agent", ctxString(ctx, httpCtxKeyUserAgent)),

			// response metadata
			zap.Int("resp.size", len(response)),

			// additional metadata
			zap.Duration("duration", time.Since(start)),
			zap.String(LogKeyRID, HTTPRequestID(ctx)))

		return response
	}
}

func ctxString(ctx context.Context, key interface{}) string {
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	return ""
}
