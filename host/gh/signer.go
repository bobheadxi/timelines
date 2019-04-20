package gh

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// SigningClient is a client dedicated to generating installation clients and
// doing other app-management stuff
type SigningClient struct {
	c   *Client
	l   *zap.SugaredLogger
	app *AppAuth
}

// NewSigningClient instantiates a new SigningClient
func NewSigningClient(l *zap.SugaredLogger, auth oauth2.TokenSource) (*SigningClient, error) {
	if auth == nil {
		return nil, errors.New("auth required")
	}
	app, ok := auth.(*AppAuth)
	if !ok {
		return nil, errors.New("auth type AppAuth is required")
	}
	client, err := newClient(l, auth)
	if err != nil {
		return nil, err
	}
	return &SigningClient{
		c:   client,
		l:   l,
		app: app,
	}, nil
}

// GetInstallationClient gets an installation-specific client
func (c *SigningClient) GetInstallationClient(ctx context.Context, id string) (*Client, error) {
	auth, err := NewInstallationAuth(ctx, c.c.gh, c.l.Named("auth"), id)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate installation: %s", err.Error())
	}
	return NewClient(ctx, c.l.Named("installation-"+id), auth)
}
