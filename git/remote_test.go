package git

import "testing"

func Test_getRepoFromRemote(t *testing.T) {
	type args struct {
		remote string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"ok", args{"https://github.com/src-d/go-git.git"}, "src-d/go-git", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRepoFromRemote(tt.args.remote)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRepoFromRemote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRepoFromRemote() = %v, want %v", got, tt.want)
			}
		})
	}
}
