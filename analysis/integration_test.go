package analysis

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/projector/git"
)

func TestGitRepoAnalyser(t *testing.T) {
	// get repo
	l := zaptest.NewLogger(t).Sugar()
	m := git.NewManager(l.Named("git-manager"), git.ManagerOpts{Workdir: "./tmp"})
	repo, err := m.Download(
		context.Background(),
		"https://github.com/bobheadxi/calories.git",
		git.DownloadOpts{})
	assert.NoError(t, err)
	defer os.RemoveAll("./tmp")

	// execute analysis
	a := NewGitAnalyser(l.Named("analysis"), repo.GitRepo(), GitRepoAnalyserOptions{})
	report, err := a.Analyze()
	assert.NoError(t, err)

	// print burndown
	b, _ := json.Marshal(report.Burndown)
	t.Log("\n== burndown:", string(b))

	// print coupling
	b, _ = json.Marshal(report.Coupling)
	t.Log("\n== coupling:", string(b))
}
