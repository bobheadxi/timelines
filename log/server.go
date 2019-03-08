package log

import (
	"go.uber.org/zap"
	"gocloud.dev/requestlog"
)

type serverLogger struct {
	l       *zap.Logger
	loggers []requestlog.Logger
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

// NewRequestLogger wraps the given logger in a request logger
func NewRequestLogger(l *zap.SugaredLogger, loggers ...requestlog.Logger) requestlog.Logger {
	return &serverLogger{l.Desugar(), loggers}
}
