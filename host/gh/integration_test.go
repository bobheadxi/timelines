package gh

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/google/go-github/v25/github"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/timelines/dev"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("-short enabled, skipping")
	}

	godotenv.Load("../../.env")
	var l = zaptest.NewLogger(t).Sugar()
	var ctx = context.Background()

	signer, err := NewSigningClient(l, NewEnvAuth())
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	ic, err := signer.GetInstallationClient(ctx, dev.GetTestInstallationID())
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	var (
		issuesC = make(chan *github.Issue, 500)
		pullsC  = make(chan *github.PullRequest, 500)
		wg      = &sync.WaitGroup{}
	)

	wg.Add(2)
	assert.NoError(t, ic.GetIssues(ctx, "bobheadxi", "calories", ItemFilter{
		State: IssueStateAll,
	}, issuesC, wg))
	assert.NoError(t, ic.GetPullRequests(ctx, "bobheadxi", "calories", ItemFilter{
		State: IssueStateAll,
	}, pullsC, wg))

	var (
		count = 0
		max   = 10
		done  = make(chan bool, 1)
	)
	for {
		select {
		case i := <-issuesC:
			t.Logf("issue #%v", i.GetNumber())
			count++
			if count == max {
				done <- true
			}
		case pr := <-pullsC:
			t.Logf("pull request #%v", pr.GetNumber())
			count++
			if count == max {
				done <- true
			}
		case <-done:
			assert.Equal(t, max, count)
			t.Log("looks good, aborting!")
			return
		}
	}
}

func TestSyncer(t *testing.T) {
	if testing.Short() {
		t.Skip("-short enabled, skipping")
	}

	godotenv.Load("../../.env")

	var (
		ctx     = context.Background()
		l       = zaptest.NewLogger(t).Sugar()
		auth    = NewEnvAuth()
		install = os.Getenv("GITHUB_TEST_INSTALLTION")
	)

	signer, err := NewSigningClient(l, auth)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	ic, err := signer.GetInstallationClient(ctx, install)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	var wg = &sync.WaitGroup{}
	var s = NewSyncer(l.Named("syncer").Named(install), ic, SyncOptions{
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

	// execute sync and wait
	outC, errC := s.Sync(ctx, wg)
	wg.Wait()
	select {
	case err := <-errC:
		t.Errorf("sync errored: %s", err.Error())
	default:
		break
	}

	// process output
	for i := range outC {
		t.Logf("received #%d", i.Number)
	}
}
