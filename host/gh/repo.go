package gh

import (
	"strconv"
	"strings"

	"github.com/bobheadxi/timelines/host"

	"github.com/google/go-github/github"
)

type githubRepo struct {
	*github.Repository
}

// RepoFromGitHub wraps a github repository in a Repo
func RepoFromGitHub(r *github.Repository) host.Repo { return &githubRepo{r} }

// ReposFromGitHub wraps a slice of repositories in a slice of Repos
func ReposFromGitHub(rs []*github.Repository) []host.Repo {
	repos := make([]host.Repo, len(rs))
	for i, r := range rs {
		repos[i] = RepoFromGitHub(r)
	}
	return repos
}

func (g *githubRepo) GetHost() host.Host { return host.HostGitHub }
func (g *githubRepo) GetID() string      { return strconv.Itoa(int(g.Repository.GetID())) }
func (g *githubRepo) GetOwner() string {
	if g.GetOwner() != "" {
		return g.GetOwner()
	}
	return strings.Split(g.Repository.GetFullName(), "/")[0]
}
func (g *githubRepo) IsPrivate() bool { return g.Repository.GetPrivate() }
