package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// AppAuth is a container for authorization configuration
type AppAuth struct {
	AppID          string
	ClientID       string
	ClientSecret   string
	SigningKeyPath string
}

// NewEnvAuth instantiates authentication from environment variables
func NewEnvAuth() *AppAuth {
	return &AppAuth{
		ClientID:       os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		SigningKeyPath: os.Getenv("GITHUB_APP_KEY"),
		AppID:          os.Getenv("GITHUB_APP_ID"),
	}
}

// GenerateJWT signs a new JWT for use with the GitHub API
func (a *AppAuth) GenerateJWT() (string, *time.Time, error) {
	priv, err := ioutil.ReadFile(a.SigningKeyPath)
	if err != nil {
		return "", nil, fmt.Errorf("could not read singing key: %s", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse signing key: %s", err)
	}

	var expiry = time.Now().Add(time.Minute)
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expiry.Unix(),
		Issuer:    a.AppID,
	}).SignedString(key)
	if err != nil {
		return "", nil, fmt.Errorf("could not sign token: %s", err)
	}

	return token, &expiry, nil
}

// Token implements oauth2.TokenSource, and is used as an autogenerating token
// source
func (a *AppAuth) Token() (*oauth2.Token, error) {
	t, exp, err := a.GenerateJWT()
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: t,
		Expiry:      *exp,
	}, nil
}

// InstallationAuth contains parameters for a specific installation
type InstallationAuth struct {
	ID int64
	gh *github.Client
	l  *zap.SugaredLogger
}

// NewInstallationAuth instantiates a new installation-specific token generator
func NewInstallationAuth(ctx context.Context, gh *github.Client, logger *zap.SugaredLogger, id string) (*InstallationAuth, error) {
	var l = logger.With("installation", id)

	install, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		l.Warnw("invalid installation ID", "error", err)
		return nil, fmt.Errorf("invalid installation '%s': %s", id, err.Error())
	}

	return &InstallationAuth{
		ID: install,
		gh: gh,
		l:  l,
	}, nil
}

// Token implements oauth2.TokenSource, and is used as an autogenerating token
// source. It queries GitHub for an installation-specific token.
func (i *InstallationAuth) Token() (*oauth2.Token, error) {
	token, _, err := i.gh.Apps.CreateInstallationToken(context.Background(), i.ID)
	if err != nil {
		i.l.Warnw("could not get token for installation", "error", err)
		return nil, fmt.Errorf("invalid installation '%d': %s", i.ID, err.Error())
	}
	i.l.Infow("generated token", "token.expiry", token.GetExpiresAt())

	return &oauth2.Token{
		AccessToken: token.GetToken(),
		Expiry:      token.GetExpiresAt(),
	}, nil
}
