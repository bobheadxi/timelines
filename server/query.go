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

func (q *queryResolver) Repo(ctx context.Context, owner string, repo string) (*models.Repository, error) {
	return nil, nil
}

func (q *queryResolver) Burndown(ctx context.Context, id int, t *models.BurndownType) (*models.Burndown, error) {
	if t == nil {
		return nil, errors.New("burndown type is required")
	}

	switch *t {
	case models.BurndownTypeFile:
		return nil, errors.New("unimplemented")
	case models.BurndownTypeAuthor:
		return nil, errors.New("unimplemented")
	case models.BurndownTypeGlobal:
		deltas, err := q.db.Repos().GetGlobalBurndown(ctx, id)
		if err != nil {
			return nil, err
		}
		return &models.Burndown{
			Type:    t,
			Entries: deltas,
		}, nil
	default:
		return nil, fmt.Errorf("invalid burndown type '%v'", t)
	}

	return nil, nil
}
