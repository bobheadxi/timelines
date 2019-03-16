package analysis

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/timelines/git"
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
	a, err := NewGitAnalyser(l.Named("analysis"), repo.GitRepo(), GitRepoAnalyserOptions{})
	assert.NoError(t, err)
	report, err := a.Analyze()
	assert.NoError(t, err)

	// print burndown
	b, _ := json.MarshalIndent(report.Burndown, "", "  ")
	os.Remove("test.burndown.json")
	ioutil.WriteFile("test.burndown.json", b, os.ModePerm)

	// print coupling
	b, _ = json.MarshalIndent(report.Coupling, "", "  ")
	os.Remove("test.coupling.json")
	ioutil.WriteFile("test.coupling.json", b, os.ModePerm)
}
