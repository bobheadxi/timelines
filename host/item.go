package host

import (
	"time"

	"github.com/google/go-github/github"
)

// ItemType denotes supported host item types
type ItemType string

const (
	// ItemTypeIssue is an issue
	ItemTypeIssue ItemType = "issue"
	// ItemTypePR is a pull request
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
	Reactions ItemReactions

	Details map[string]interface{}
}

// ItemReactions denotes reactions to an item
type ItemReactions struct {
	Total    int
	Positive int
	Negative int
}

// WithGitHubLabels attaches given labels as strings to Item
func (i *Item) WithGitHubLabels(labels []github.Label) {
	i.Labels = make([]string, len(labels))
	for x, l := range labels {
		i.Labels[x] = l.GetName()
	}
}

// WithGitHubReactions sets reactions after stripping out unecessary things
func (i *Item) WithGitHubReactions(r *github.Reactions) {
	if r == nil {
		return
	}
	i.Reactions.Positive = r.GetHeart() + r.GetHooray() + r.GetLaugh() + r.GetPlusOne()
	i.Reactions.Negative = r.GetConfused() + r.GetMinusOne()
	i.Reactions.Total = r.GetTotalCount()
}
