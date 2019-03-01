package github

import "go.uber.org/zap"

// Syncer manages all GitHub synchronization tasks
type Syncer struct {
	c *Client
	l *zap.SugaredLogger
}

// NewSyncer instantiates a new GitHub Syncer
func NewSyncer(l *zap.SugaredLogger, client *Client) *Syncer {
	return nil
}
