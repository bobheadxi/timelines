package store

import (
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/projector/config"
	"github.com/bobheadxi/projector/dev"
)

func TestNewClient(t *testing.T) {
	type args struct {
		opts config.Store
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"missing params", args{config.Store{}}, true},
		{"invalid authentication", args{config.Store{
			Address:  "127.0.0.1:6379",
			Password: "i_love_chicken_rice",
		}}, true},
		{"connect to devenv", args{dev.StoreOptions}, false},
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
