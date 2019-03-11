package store

import (
	"errors"
	"fmt"

	"github.com/bobheadxi/projector/config"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// Client is the interface to the redis store
type Client struct {
	redis *redis.Client

	l *zap.SugaredLogger
}

// NewClient sets up a new client for the redis store
func NewClient(l *zap.SugaredLogger, opts config.Store) (*Client, error) {
	if opts.Address == "" {
		return nil, errors.New("no address provided")
	}

	var c = redis.NewClient(&redis.Options{
		Addr:      opts.Address,
		Password:  opts.Password,
		TLSConfig: opts.TLS,

		DB: 0, // use default DB
	})

	l.Infow("pinging redis...")
	var s = c.Ping()
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis store: %s", err.Error())
	}
	l.Infow("ping successful", "response", s.Val())

	return &Client{
		redis: c,
		l:     l,
	}, nil
}

// RepoJobs returns a client for managing repo job entries
func (c *Client) RepoJobs() *RepoJobsClient {
	return &RepoJobsClient{c: c, l: c.l.Named(repoJobsName)}
}

// Reset drops all keys
func (c *Client) Reset() {
	var (
		repoJobs, _   = c.redis.Keys(queueRepoJobs + "*").Result()
		repoStates, _ = c.redis.Keys(statesRepoJobs + "*").Result()

		keys = append(repoJobs, repoStates...)
	)
	if len(keys) > 0 {
		if err := c.redis.Del(keys...).Err(); err != nil {
			c.l.Errorw("failed to remove keys", "error", err)
		}
	}
}

// Close disconnects the client from Redis
func (c *Client) Close() error { return c.redis.Close() }
