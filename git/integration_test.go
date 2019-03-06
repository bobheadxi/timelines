package git

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestManager(t *testing.T) {
	var (
		ctx = context.Background()
		l   = zaptest.NewLogger(t).Sugar()
		m   = NewManager(l.Named("manager"), ManagerOpts{Workdir: "tmp"})
	)
	repo, err := m.Download(ctx, "https://github.com/bobheadxi/calories.git", DownloadOpts{})
	assert.NoError(t, err)
	defer os.RemoveAll(repo.Dir())

	name, err := repo.Name()
	assert.NoError(t, err)
	assert.Equal(t, "bobheadxi/calories", name)
}
