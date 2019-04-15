package server

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/graphql/go/timelines/models"
)

type queryResolver struct {
	db *db.Database

	l *zap.SugaredLogger
}

func (q *queryResolver) Repo(ctx context.Context, owner, name string, h *models.RepositoryHost) (*models.Repository, error) {
	if h == nil {
		*h = models.RepositoryHostGithub
	}

	repo, err := q.db.Repos().GetRepository(ctx, owner, name)
	if err != nil {
		if !db.IsNotFound(err) {
			q.l.Errorw(err.Error(),
				"owner", owner,
				"name", name)
		}
		return nil, fmt.Errorf("could not find repository for '%s/%s'", owner, name)
	}
	return repo, nil
}

func (q *queryResolver) Repos(ctx context.Context, owner string, h *models.RepositoryHost) ([]models.Repository, error) {
	if h == nil {
		*h = models.RepositoryHostGithub
	}

	repos, err := q.db.Repos().GetRepositories(ctx, owner)
	if err != nil {
		q.l.Errorw(err.Error(),
			"owner", owner)
		return nil, fmt.Errorf("could not find repositories for '%s'", owner)
	}
	return repos, nil
}

func (q *queryResolver) Burndown(ctx context.Context, id int, t *models.BurndownType) (*models.Burndown, error) {
	if t == nil {
		*t = models.BurndownTypeGlobal
	}
	l := q.l.With("repo", id, "type", *t)

	switch *t {
	case models.BurndownTypeFile:
		return nil, errors.New("unimplemented")
	case models.BurndownTypeAuthor:
		return nil, errors.New("unimplemented")
	case models.BurndownTypeGlobal:
		deltas, err := q.db.Repos().GetGlobalBurndown(ctx, id)
		if err != nil {
			l.Error(err.Error())
			return nil, fmt.Errorf("could not find '%s' burndowns for repo '%d'", t, id)
		}
		return &models.Burndown{
			Type:    t,
			Entries: deltas,
		}, nil
	default:
		return nil, fmt.Errorf("invalid burndown type '%v'", t)
	}
}
