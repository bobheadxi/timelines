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
	remote, err := r.origin()
	if err != nil {
		return "", err
	}
	name, _, err := getRepoFromRemote(remote)
	return name, err
}

// Host gets the code host of this repository
func (r *Repository) Host() (Host, error) {
	remote, err := r.origin()
	if err != nil {
		return "", err
	}
	_, host, err := getRepoFromRemote(remote)
	return host, err
}

// Dir gets where this repository is stored
func (r *Repository) Dir() string { return r.dir }

// GitRepo returns the underlying go-git repository
func (r *Repository) GitRepo() *gogit.Repository { return r.git }

func (r *Repository) origin() (string, error) {
	remote, err := r.git.Remote("origin")
	if err != nil {
		return "", err
	}
	if len(remote.Config().URLs) < 1 {
		return "", errors.New("no URLs configured for remote 'origin'")
	}
	return remote.Config().URLs[0], nil
}
