package gh

import (
	"net/http"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Client struct {
	gh *github.Client

	l *zap.Logger
}

func NewClient(token string) *Client {
	var tc *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		})
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}
	return &Client{gh: github.NewClient(tc)}
}
