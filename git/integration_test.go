package git

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/timelines/host/gh"
)

func TestManager(t *testing.T) {
	var (
		ctx = context.Background()
		l   = zaptest.NewLogger(t).Sugar()
		m   = NewManager(l.Named("manager"), ManagerOpts{Workdir: "tmp"})
	)

	// try without auth
	repo, err := m.Download(ctx, "https://github.com/bobheadxi/calories.git", DownloadOpts{})
	assert.NoError(t, err)
	name, err := repo.Name()
	assert.NoError(t, err)
	assert.Equal(t, "bobheadxi/calories", name)
	os.RemoveAll(repo.Dir())

	// try with auth
	if testing.Short() {
		t.Log("skipping authenticated test")
		return
	}
	godotenv.Load("../.env")
	s, err := gh.NewSigningClient(l, gh.NewEnvAuth())
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	ic, err := s.GetInstallationClient(ctx, os.Getenv("GITHUB_TEST_INSTALLTION"))
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	repo, err = m.Download(ctx, "https://github.com/bobheadxi/calories.git", DownloadOpts{
		AccessToken: ic.InstallationToken(),
	})
	assert.NoError(t, err)
	name, err = repo.Name()
	assert.NoError(t, err)
	assert.Equal(t, "bobheadxi/calories", name)
	os.RemoveAll(repo.Dir())
}
