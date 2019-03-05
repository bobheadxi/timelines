package github

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("-short enabled, skipping")
	}
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

	assert.NoError(t, ic.GetIssues(ctx, "bobheadxi", "calories", ItemFilter{
		State: IssueStateAll,
	}, issuesC, pullsC, wg))

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

func TestSyncer(t *testing.T) {
	if testing.Short() {
		t.Skip("-short enabled, skipping")
	}

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

	var wg = &sync.WaitGroup{}
	var s = NewSyncer(l.Named("syncer"), ic, SyncOptions{
		Repo: Repo{
			Owner: "bobheadxi",
			Name:  "calories",
		},
		Filter: ItemFilter{
			State: IssueStateAll,
		},
		DetailsFetchWorkers: 3,
		OutputBufferSize:    500,
	})

	outC, errC := s.Sync(ctx, wg)
	wg.Wait()
	select {
	case err := <-errC:
		t.Errorf("sync errored: %s", err.Error())
	default:
		break
	}
	go func() {
		for i := range outC {
			t.Logf("received #%d", i.Number)
		}
	}()
	time.Sleep(10 * time.Second)
}
