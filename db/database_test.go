package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

var devEnvOptions = Options{
	Address:  "127.0.0.1:5431",
	Database: "projector_dev",
	User:     "bobheadxi",
	Password: "bobheadxi",
}

func TestNew(t *testing.T) {
	type args struct {
		opts Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"missing params", args{Options{}}, true},
		{"invalid authentication", args{Options{
			Address: "127.0.0.1:5431",
		}}, true},
		{"connect to devenv", args{devEnvOptions}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l = zaptest.NewLogger(t).Sugar()
			_, err := New(l, tt.args.opts)
			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}
