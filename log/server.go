package log

import (
	"go.uber.org/zap"
	"gocloud.dev/requestlog"
)

type serverLogger struct {
	l       *zap.Logger
	loggers []requestlog.Logger
}

// NewRequestLogger wraps the given logger in gocloud's requestlog.Logger
func NewRequestLogger(l *zap.SugaredLogger, loggers ...requestlog.Logger) requestlog.Logger {
	return &serverLogger{
		// don't take stacktrace of wrapper class
		l.Desugar().WithOptions(zap.AddCallerSkip(1)),
		loggers,
	}
}

func (s *serverLogger) Log(e *requestlog.Entry) {
	if e == nil {
		return
	}

	s.l.Info(e.RequestMethod+" "+e.RequestURL,
		zap.Duration("latency", e.Latency),

		zap.String("request.method", e.RequestMethod),
		zap.String("request.url", e.RequestURL),
		zap.Int64("request.header_size", e.RequestHeaderSize),
		zap.Int64("request.body_size", e.RequestBodySize),
		zap.String("request.user_agent", e.UserAgent),
		zap.String("request.referer", e.Referer),

		zap.Int("response.status", e.Status),
		zap.Int64("response.header_size", e.ResponseHeaderSize),
		zap.Int64("response.body_size", e.ResponseBodySize))

	for _, l := range s.loggers {
		l.Log(e)
	}
}
