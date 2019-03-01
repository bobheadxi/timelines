package github

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestGitHubClient(t *testing.T) {
	godotenv.Load("../.env")
	var l = zaptest.NewLogger(t).Sugar()
	var ctx = context.Background()

	signer, err := NewSigningClient(l, NewEnvAuth())
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	ic, err := signer.GetInstallationClient(ctx, os.Getenv("GITHUB_TEST_INSTALLTION"))
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	var (
		issuesC = make(chan *github.Issue, 500)
		pullsC  = make(chan *github.Issue, 500)
		wg      = &sync.WaitGroup{}
	)

	ic.GetIssues(ctx, "bobheadxi", "calories", IssuesFilter{
		State: IssueStateAll,
	}, issuesC, pullsC, wg)

	unauth, err := NewClient(ctx, l, nil)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	var (
		count = 0
		max   = 5
		done  = make(chan bool, 1)
	)
	for {
		select {
		case i := <-issuesC:
			t.Logf("issue # %v", i.GetNumber())
			count++
			if count == max {
				done <- true
			}
		case pr := <-pullsC:
			t.Logf("pull request # %v", pr.GetNumber())
			pull, err := unauth.GetPullRequest(ctx, pr)
			if assert.NoError(t, err) {
				t.Logf("pull request has %v lines of addition", pull.GetAdditions())
			}
			count++
			if count == max {
				done <- true
			}
		case <-done:
			assert.NotZero(t, count)
			t.Log("looks good, aborting!")
			return
		}
	}
}
