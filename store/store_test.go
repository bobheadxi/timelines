package store

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

var devEnvOptions = Options{
	Address:  "127.0.0.1:6379",
	Password: "",
}

func TestNewClient(t *testing.T) {
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
			Address:  "127.0.0.1:6379",
			Password: "i_love_chicken_rice",
		}}, true},
		{"connect to devenv", args{devEnvOptions}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l = zaptest.NewLogger(t).Sugar()
			_, err := NewClient(l, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
