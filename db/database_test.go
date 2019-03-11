package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/projector/config"
	"github.com/bobheadxi/projector/dev"
)

func TestNew(t *testing.T) {
	type args struct {
		opts config.Database
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"missing params", args{config.Database{}}, true},
		{"invalid authentication", args{config.Database{
			Host: "127.0.0.1",
			Port: "5431",
		}}, true},
		{"connect to devenv", args{dev.DatabaseOptions}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l = zaptest.NewLogger(t).Sugar()
			_, err := New(l, tt.name, tt.args.opts)
			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}
