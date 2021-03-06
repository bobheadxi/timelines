package store

import (
	"errors"
	"fmt"

	"github.com/bobheadxi/timelines/config"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// Client is the interface to the redis store
type Client struct {
	name  string
	redis *redis.Client

	l *zap.SugaredLogger
}

// NewClient sets up a new client for the redis store
func NewClient(l *zap.SugaredLogger, name string, opts config.Store) (*Client, error) {
	// set up redis client
	var c *redis.Client
	if opts.RedisConnURL != "" {
		l.Info("connecting using connection URL")
		cfg, err := redis.ParseURL(opts.RedisConnURL)
		if err != nil {
			return nil, err
		}
		c = redis.NewClient(cfg)
		l.Infow("client successfully set up",
			"address", cfg.Addr,
			"db", cfg.DB)
	} else {
		l.Info("connecting using parameters")
		if opts.Address == "" {
			return nil, errors.New("store: no address provided")
		}
		c = redis.NewClient(&redis.Options{
			Addr:      opts.Address,
			Password:  opts.Password,
			TLSConfig: opts.TLS,

			DB: 0, // use default DB
		})
		l.Infow("client successfully set up",
			"address", opts.Address,
			"db", 0)
	}

	// make sure all's good
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
	return &RepoJobsClient{name: c.name, c: c, l: c.l.Named(repoJobsName)}
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
