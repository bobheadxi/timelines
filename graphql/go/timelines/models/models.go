// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Burndown interface {
	IsBurndown()
}

type AuthorBurndown struct {
	RepoID int              `json:"repoID"`
	Type   BurndownType     `json:"type"`
	Author []*BurndownEntry `json:"author"`
}

func (AuthorBurndown) IsBurndown() {}

type BurndownAlert struct {
	RepoID int          `json:"repoID"`
	Type   BurndownType `json:"type"`
	Alert  string       `json:"alert"`
}

func (BurndownAlert) IsBurndown() {}

type BurndownEntry struct {
	Start time.Time `json:"start"`
	Bands []int     `json:"bands"`
}

type FileBurndown struct {
	RepoID int                  `json:"repoID"`
	Type   BurndownType         `json:"type"`
	File   []*FileBurndownEntry `json:"file"`
}

func (FileBurndown) IsBurndown() {}

type FileBurndownEntry struct {
	File  string         `json:"file"`
	Entry *BurndownEntry `json:"entry"`
}

type GlobalBurndown struct {
	RepoID  int              `json:"repoID"`
	Type    BurndownType     `json:"type"`
	Entries []*BurndownEntry `json:"entries"`
}

func (GlobalBurndown) IsBurndown() {}

type Repository struct {
	ID          int    `json:"id"`
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RepositoryAnalytics struct {
	Repository *Repository `json:"repository"`
	Burndown   Burndown    `json:"burndown"`
}

type ServiceStatus struct {
	Build    string    `json:"build"`
	Deployed time.Time `json:"deployed"`
}

type BurndownType string

const (
	BurndownTypeGlobal BurndownType = "GLOBAL"
	BurndownTypeFile   BurndownType = "FILE"
	BurndownTypeAuthor BurndownType = "AUTHOR"
	BurndownTypeAlert  BurndownType = "ALERT"
)

var AllBurndownType = []BurndownType{
	BurndownTypeGlobal,
	BurndownTypeFile,
	BurndownTypeAuthor,
	BurndownTypeAlert,
}

func (e BurndownType) IsValid() bool {
	switch e {
	case BurndownTypeGlobal, BurndownTypeFile, BurndownTypeAuthor, BurndownTypeAlert:
		return true
	}
	return false
}

func (e BurndownType) String() string {
	return string(e)
}

func (e *BurndownType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BurndownType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BurndownType", str)
	}
	return nil
}

func (e BurndownType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RepositoryHost string

const (
	RepositoryHostGithub    RepositoryHost = "GITHUB"
	RepositoryHostGitlab    RepositoryHost = "GITLAB"
	RepositoryHostBitbucket RepositoryHost = "BITBUCKET"
)

var AllRepositoryHost = []RepositoryHost{
	RepositoryHostGithub,
	RepositoryHostGitlab,
	RepositoryHostBitbucket,
}

func (e RepositoryHost) IsValid() bool {
	switch e {
	case RepositoryHostGithub, RepositoryHostGitlab, RepositoryHostBitbucket:
		return true
	}
	return false
}

func (e RepositoryHost) String() string {
	return string(e)
}

func (e *RepositoryHost) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RepositoryHost(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RepositoryHost", str)
	}
	return nil
}

func (e RepositoryHost) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
