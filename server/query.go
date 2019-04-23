package server

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/graphql/go/timelines"
	"github.com/bobheadxi/timelines/graphql/go/timelines/models"
	"github.com/bobheadxi/timelines/host"
	"github.com/bobheadxi/timelines/log"
)

type queryResolver struct {
	db *db.Database

	l *zap.SugaredLogger
}

func newQueryResolver(
	l *zap.SugaredLogger,
	database *db.Database,
) timelines.QueryResolver {
	return &queryResolver{database, l.Named("query")}
}

func (q *queryResolver) Repo(
	ctx context.Context,
	owner, name string, h *models.RepositoryHost,
) (*models.RepositoryAnalytics, error) {
	var l = q.l.With(log.LogKeyRID, log.HTTPRequestID(ctx),
		"owner", owner, "name", name)
	hostService, err := modelToHost(h)
	if err != nil {
		return nil, err
	}

	repo, err := q.db.Repos().GetRepository(ctx, hostService, owner, name)
	if err != nil {
		if !db.IsNotFound(err) {
			l.Errorw(err.Error())
		}
		return nil, fmt.Errorf("could not find repository for '%s/%s'", owner, name)
	}

	return &models.RepositoryAnalytics{
		Repository: *repo,
	}, nil
}

func (q *queryResolver) Repos(
	ctx context.Context,
	owner string, h *models.RepositoryHost,
) ([]models.Repository, error) {
	var l = q.l.With(log.LogKeyRID, log.HTTPRequestID(ctx),
		"owner", owner)
	hostService, err := modelToHost(h)
	if err != nil {
		return nil, err
	}

	repos, err := q.db.Repos().GetRepositories(ctx, hostService, owner)
	if err != nil {
		if !db.IsNotFound(err) {
			l.Errorw(err.Error())
		}
		return nil, fmt.Errorf("could not find repositories for '%s'", owner)
	}
	return repos, nil
}

func modelToHost(h *models.RepositoryHost) (host.Host, error) {
	if h == nil {
		return host.HostGitHub, nil
	}
	hostService := host.Host(h.String())
	switch hostService {
	case host.HostGitHub, host.HostGitLab, host.HostBitbucket:
		return hostService, nil
	default:
		return "", fmt.Errorf("unknown host '%v'", h)
	}
}
