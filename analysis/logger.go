package analysis

import (
	"go.uber.org/zap"
	"gopkg.in/src-d/hercules.v10"
)

type herculesLogger struct {
	*zap.SugaredLogger
}

// newHerculesLogger wrapps the given logger in hercules.Logger
func newHerculesLogger(l *zap.SugaredLogger) hercules.Logger {
	return &herculesLogger{
		// don't take stacktrace of wrapper class
		l.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
	}
}

func (h *herculesLogger) Critical(v ...interface{}) {
	h.With("critical", true).Error(v...)
}

func (h *herculesLogger) Criticalf(f string, v ...interface{}) {
	h.With("critical", true).Errorf(f, v...)
}
