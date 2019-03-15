package git

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Manager handles repo management
type Manager struct {
	opts ManagerOpts

	l *zap.SugaredLogger
}

// ManagerOpts denotes options for the manager
type ManagerOpts struct {
	Workdir string
}

// NewManager instantiates a new manager
func NewManager(l *zap.SugaredLogger, opts ManagerOpts) *Manager {
	if opts.Workdir == "" {
		opts.Workdir = "./"
	}
	return &Manager{
		opts: opts,
		l:    l,
	}
}

// DownloadOpts denotes options for downloading a repo
type DownloadOpts struct {
	AccessToken string
	Depth       int
}

// Download downloads the given repository
func (m *Manager) Download(ctx context.Context, remote string, opts DownloadOpts) (*Repository, error) {
	var l = m.l.With("remote", remote)
	repoDir, err := m.repoDir(remote)
	if err != nil {
		return nil, err
	}
	os.RemoveAll(repoDir)

	// set auth if provided
	var auth *http.BasicAuth
	if opts.AccessToken != "" {
		auth = &http.BasicAuth{
			Username: "x-access-token",
			Password: opts.AccessToken,
		}
	}

	// grab repo
	l.Infow("cloning repository",
		"dir", repoDir,
		"auth", opts.AccessToken != "")
	gitrepo, err := gogit.PlainCloneContext(
		ctx,
		repoDir,
		true,
		&gogit.CloneOptions{
			URL:          remote,
			SingleBranch: true,

			Depth: opts.Depth,
			Auth:  auth,
		})
	if err != nil {
		l.Errorw("failed to clone repo", "error", err)
		return nil, err
	}
	l.Info("repo cloned")
	return &Repository{
		git: gitrepo,
		dir: repoDir,
	}, nil
}

// Load loads a downloaded repository - TODO: deprecate in favour of just using
// Donwload(), to unify depth management
func (m *Manager) Load(ctx context.Context, remote string) (*Repository, error) {
	repodir, err := m.repoDir(remote)
	if err != nil {
		return nil, err
	}
	gitrepo, err := gogit.PlainOpen(repodir)
	if err != nil {
		return nil, err
	}
	gitrepo.FetchContext(ctx, &gogit.FetchOptions{
		Force: true,

		// TODO: private repos
		Auth: nil,
	})
	return &Repository{
		dir: repodir,
		git: gitrepo,
	}, nil
}

func (m *Manager) repoDir(remote string) (string, error) {
	repo, host, err := getRepoFromRemote(remote)
	if err != nil {
		return "", err
	}
	var path = filepath.Join(m.opts.Workdir, string(host), repo)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}
	return path, nil
}
