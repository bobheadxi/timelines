package git

import (
	"fmt"
	"strings"
)

func getRepoFromRemote(remote string) (string, Host, error) {
	var parts = strings.Split(remote, "/")
	if len(parts) != 5 {
		return "", "", fmt.Errorf("unexpected number of components in remote '%s'", remote)
	}

	var h Host
	switch parts[2] {
	case string(HostGitHub):
		h = HostGitHub
	case string(HostGitLab):
		h = HostGitLab
	default:
		return "", "", fmt.Errorf("unknown host '%s' in remote '%s'", parts[2], remote)
	}
	parts[len(parts)-1] = strings.TrimSuffix(parts[len(parts)-1], ".git")
	return strings.Join(parts[3:], "/"), h, nil
}
