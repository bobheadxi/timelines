package host

// Host denotes supported hosts
type Host string

const (
	// HostGitHub is github.com
	HostGitHub Host = "GITHUB"
	// HostGitLab is gitlab.com
	HostGitLab Host = "GITLAB"
	// HostBitbucket is bitbucket.org
	HostBitbucket Host = "BITBUCKET"
)

// Hosted defines a common interface for all host struct types
type Hosted interface {
	GetID() string
	GetHost() Host
}
