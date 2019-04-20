package host

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
)

// Repo represents a repository
type Repo interface {
	Hosted
	GetOwner() string
	GetName() string
	GetDescription() string
	IsPrivate() bool
	fmt.Stringer
}

// BaseRepo is an implementation of Repo, should probably be only used for testing
type BaseRepo struct {
	Host        Host
	ID          int
	Owner       string
	Name        string
	Description string
	Private     bool
}

// GetHost returns the repo's host
func (b *BaseRepo) GetHost() Host { return b.Host }

// GetID returns the repo's ID
func (b *BaseRepo) GetID() string { return strconv.Itoa(b.ID) }

// GetOwner returns the repo's owner
func (b *BaseRepo) GetOwner() string { return b.Owner }

// GetName returns the repo's name
func (b *BaseRepo) GetName() string { return b.Name }

// GetDescription returns the repo's description
func (b *BaseRepo) GetDescription() string { return b.Description }

// IsPrivate indicates whether the repo is private of not
func (b *BaseRepo) IsPrivate() bool { return b.Private }

func (b *BaseRepo) String() string { return fmt.Sprintf("%s:%s/%s", b.Host, b.Owner, b.Name) }

// doIstillImplementRepo is a dumb way of checking that BaseRepo is up to date
// with the Repo interface
func (b *BaseRepo) doIstillImplementRepo() Repo { return b }

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
func (g *githubRepo) GetOwner() string { return strings.Split(g.Repository.GetFullName(), "/")[0] }
func (g *githubRepo) IsPrivate() bool  { return g.Repository.GetPrivate() }
