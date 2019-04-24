package host

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/v25/github"
)

// Installation represents an app installation on a code host
type Installation interface {
	Hosted
	fmt.Stringer
}

type githubInstall struct {
	*github.Installation
}

// InstallationFromGitHub wraps a github installation in an Installation
func InstallationFromGitHub(i *github.Installation) Installation { return &githubInstall{i} }

func (g *githubInstall) GetHost() Host { return HostGitHub }
func (g *githubInstall) GetID() string { return strconv.Itoa(int(g.Installation.GetID())) }
