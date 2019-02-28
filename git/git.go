package git

import (
	// TODO: explore gitbase
	// gitbase "gopkg.in/src-d/gitbase.v0"
	gogit "gopkg.in/src-d/go-git.v4"
)

type Repository struct {
	git *gogit.Repository
}
