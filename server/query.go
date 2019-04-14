package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/graphql/go/timelines/models"
)

type queryResolver struct {
	l *zap.SugaredLogger
}

func (q *queryResolver) Repo(ctx context.Context, owner string, repo string) (*models.Repository, error) {
	return nil, nil
}

func (q *queryResolver) Burndown(ctx context.Context, id int, typeArg *models.BurndownType) (*models.Burndown, error) {
	return nil, nil
}
