package host

// Host denotes supported hosts
type Host string

const (
	// HostGitHub is GitHub
	HostGitHub Host = "github"
)

// Hosted defines a common interface for all host struct types
type Hosted interface {
	GetID() string
	GetHost() Host
}
