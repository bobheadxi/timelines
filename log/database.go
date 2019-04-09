package log

import (
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type databaseLogger struct {
	l *zap.Logger
}

// NewDatabaseLogger wraps the given logger in pgx.Logger
func NewDatabaseLogger(l *zap.SugaredLogger) pgx.Logger {
	return &databaseLogger{
		// don't take stacktrace of wrapper class
		l: l.Desugar().WithOptions(zap.AddCallerSkip(1)),
	}
}

func (d *databaseLogger) Log(lv pgx.LogLevel, msg string, context map[string]interface{}) {
	var zapData = zap.Any("context", context)
	switch lv {
	case pgx.LogLevelDebug:
		d.l.Debug(msg, zapData)
	case pgx.LogLevelWarn:
		d.l.Warn(msg, zapData)
	case pgx.LogLevelError:
		d.l.Error(msg, zapData)
	default:
		d.l.Info(msg, zapData)
	}
}
