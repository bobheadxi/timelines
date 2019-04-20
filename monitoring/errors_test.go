package monitoring

import (
	"testing"

	"cloud.google.com/go/errorreporting"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bobheadxi/timelines/log"
)

type fakeGCPReporter struct {
	lastEntry errorreporting.Entry
	flushed   bool
}

func (f *fakeGCPReporter) Report(e errorreporting.Entry) { f.lastEntry = e }
func (f *fakeGCPReporter) Flush()                        { f.flushed = true }

func Test_gcpErrorReportingZapCore_Enabled(t *testing.T) {
	var f fakeGCPReporter
	core := gcpErrorsWrapCore(&f)
	assert.False(t, core.Enabled(zapcore.DebugLevel))
	assert.False(t, core.Enabled(zapcore.InfoLevel))
	assert.True(t, core.Enabled(zapcore.WarnLevel))
	assert.True(t, core.Enabled(zapcore.ErrorLevel))
}

func Test_stacktrace(t *testing.T) {
	stack := stacktrace()
	assert.NotNil(t, stack)
}

func Test_gcpErrorReportingZapCore_Write(t *testing.T) {
	type args struct {
		entry  zapcore.Entry
		fields []zapcore.Field
	}
	type want struct {
		user string
		err  string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{"all empty", args{
			zapcore.Entry{}, nil,
		}, want{}, false},
		{"message", args{
			zapcore.Entry{Message: "hello world"}, nil,
		}, want{err: "hello world"}, false},
		{"request ID", args{
			zapcore.Entry{},
			[]zapcore.Field{
				zap.String(log.LogKeyRID, "1234"),
			},
		}, want{user: "1234"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f fakeGCPReporter
			core := gcpErrorsWrapCore(&f)
			if tt.wantErr {
				assert.Error(t, core.Write(tt.args.entry, tt.args.fields))
			} else {
				assert.NoError(t, core.Write(tt.args.entry, tt.args.fields))
			}
			assert.Equal(t, tt.want.user, f.lastEntry.User)
			assert.Contains(t, f.lastEntry.Error.Error(), tt.want.err)
		})
	}
}

func Test_gcpErrorReportingZapCore_Sync(t *testing.T) {
	var f fakeGCPReporter
	core := gcpErrorsWrapCore(&f)
	assert.False(t, f.flushed)
	assert.NoError(t, core.Sync())
	assert.True(t, f.flushed)
}
