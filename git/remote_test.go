package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getRepoFromRemote(t *testing.T) {
	type args struct {
		remote string
	}
	tests := []struct {
		name     string
		args     args
		wantRepo string
		wantHost Host
		wantErr  bool
	}{
		{"not ok - invalid format", args{"asdfasdf"}, "", "", true},
		{"not ok - invalid host", args{"https://bobheadxi.dev/src-d/go-git.git"}, "", "", true},
		{"ok - github",
			args{"https://github.com/src-d/go-git.git"},
			"src-d/go-git",
			HostGitHub, false},
		{"ok - gitlab",
			args{"https://gitlab.com/gitlab-org/gitlab-ce.git"},
			"gitlab-org/gitlab-ce",
			HostGitLab, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, host, err := getRepoFromRemote(tt.args.remote)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantRepo, repo)
			assert.Equal(t, tt.wantHost, host)
		})
	}
}
