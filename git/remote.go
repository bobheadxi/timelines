package git

import (
	"fmt"
	"strings"
)

func getRepoFromRemote(remote string) (string, error) {
	var parts = strings.Split(remote, "/")
	if len(parts) != 5 {
		return "", fmt.Errorf("unexpected number of components in remote '%s'", remote)
	}
	parts[len(parts)-1] = strings.TrimSuffix(parts[len(parts)-1], ".git")
	return strings.Join(parts[3:], "/"), nil
}
