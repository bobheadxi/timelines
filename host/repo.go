package host

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
)

// Repo represents a repository
type Repo interface {
	Hosted
	GetOwner() string
	GetName() string
	fmt.Stringer
}

type githubRepo struct {
	*github.Repository
}

// RepoFromGitHub wraps a github repository in a Repo
func RepoFromGitHub(r *github.Repository) Repo { return &githubRepo{r} }

// ReposFromGitHub wraps a slice of repositories in a slice of Repos
func ReposFromGitHub(rs []*github.Repository) []Repo {
	repos := make([]Repo, len(rs))
	for i, r := range rs {
		repos[i] = RepoFromGitHub(r)
	}
	return repos
}

func (g *githubRepo) GetHost() Host    { return HostGitHub }
func (g *githubRepo) GetID() string    { return strconv.Itoa(int(g.Repository.GetID())) }
func (g *githubRepo) GetOwner() string { return g.Repository.GetOwner().GetName() }
