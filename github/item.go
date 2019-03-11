package github

import (
	"time"

	"github.com/google/go-github/github"
)

// ItemType denotes supported GitHub item types
type ItemType string

const (
	// ItemTypeIssue is a GitHub issue
	ItemTypeIssue ItemType = "issue"
	// ItemTypePR is a GitHub pull request
	ItemTypePR ItemType = "pull_request"
)

// Item is a GitHub item due for indexing
// TODO: this needs to be better
type Item struct {
	GitHubID int
	Number   int

	Author string
	Opened time.Time
	Closed *time.Time

	Type      ItemType
	Title     string
	Body      string
	Labels    []string
	Reactions *github.Reactions

	Details map[string]interface{}
}

func (i *Item) WithLabels(labels []github.Label) {
	i.Labels = make([]string, len(labels))
	for x, l := range labels {
		i.Labels[x] = l.GetName()
	}
}

func (i *Item) WithReactions(reacs *github.Reactions) {
	if reacs == nil {
		return
	}
	reacs.URL = nil
	i.Reactions = reacs
}
