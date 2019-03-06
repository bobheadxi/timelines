package git

import (
	"errors"

	gogit "gopkg.in/src-d/go-git.v4"
)

// Repository represents a managed repository, cloned for analysis
type Repository struct {
	dir string
	git *gogit.Repository
}

// Name gets the repository name
func (r *Repository) Name() (string, error) {
	remote, err := r.git.Remote("origin")
	if err != nil {
		return "", err
	}
	if len(remote.Config().URLs) < 1 {
		return "", errors.New("no URLs configured for remote 'origin'")
	}
	return getRepoFromRemote(remote.Config().URLs[0])
}

// Dir gets where this repository is stored
func (r *Repository) Dir() string { return r.dir }

// GitRepo returns the underlying go-git repository
func (r *Repository) GitRepo() *gogit.Repository { return r.git }
