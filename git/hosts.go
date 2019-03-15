package git

// Host denotes supported git hosts
type Host string

const (
	// HostGitHub is https://github.com
	HostGitHub Host = "github.com"

	// HostGitLab is https://gitlab.com
	HostGitLab Host = "gitlab.com"
)
