package git

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	gogit "gopkg.in/src-d/go-git.v4"
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
type DownloadOpts struct{}

// Download downloads the given repository
func (m *Manager) Download(ctx context.Context, remote string, opts DownloadOpts) (*Repository, error) {
	var l = m.l.With("remote", remote)
	repoDir, err := m.repoDir(remote)
	if err != nil {
		return nil, err
	}
	l.Infow("cloning repository", "dir", repoDir)

	gitrepo, err := gogit.PlainCloneContext(
		ctx,
		repoDir,
		true,
		&gogit.CloneOptions{
			URL: remote,

			// TODO: private repos?
			Auth: nil,
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

// Load loads a downloaded repository
func (m *Manager) Load(remote string) (*Repository, error) {
	return nil, nil
}

func (m *Manager) repoDir(remote string) (string, error) {
	repo, err := getRepoFromRemote(remote)
	if err != nil {
		return "", err
	}
	var path = filepath.Join(m.opts.Workdir, repo)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}
	return path, nil
}
