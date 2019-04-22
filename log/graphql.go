package log

import (
	"context"
	"net/http"
	"strings"
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
// graphql.FieldMiddleware is super verbose, and things like introspection
// seem to cause the log to blow up. not sure how useful this is, though on the
// other hand graphql.RequestMiddleware isn't very informative
func NewGraphLogger(l *zap.Logger) graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		// call handler
		res := graphql.GetResolverContext(ctx)
		if strings.HasPrefix(res.Object, "__") || !res.IsMethod {
			return next(ctx)
		}

		start := time.Now()
		response, err := next(ctx)

		// log request
		// TODO: could implement more advanced tracing, hm, via RequestContext.Trace
		req := graphql.GetRequestContext(ctx)
		// TODO: log message not very informative
		// TODO: evaluate usefulness of logged fields
		l.Info(res.Object+": graph query completed",
			zap.Bool("ismethod", res.IsMethod),
			zap.Any("path", res.Path()),
			// request metadata
			zap.Int("req.complexity", req.OperationComplexity),
			zap.Any("req.variables", req.Variables),
			zap.Any("req.errors", req.Errors),
			zap.Any("req.extensions", req.Extensions),
			zap.String("req.ip", ctxString(ctx, httpCtxKeyRemoteAddr)),
			zap.String("req.user_agent", ctxString(ctx, httpCtxKeyUserAgent)),

			// response metadata
			zap.NamedError("resp.err", err),

			// additional metadata
			zap.Duration("duration", time.Since(start)),
			zap.String(LogKeyRID, HTTPRequestID(ctx)))

		return response, err
	}
}

func ctxString(ctx context.Context, key interface{}) string {
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	return ""
}
