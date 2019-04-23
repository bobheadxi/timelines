package server

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/graphql/go/timelines"
	"github.com/bobheadxi/timelines/graphql/go/timelines/models"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/store"
)

type analyticsResolver struct {
	db *db.Database
	s  *store.Client

	l *zap.SugaredLogger
}

func newAnalyticsResolver(
	l *zap.SugaredLogger,
	database *db.Database,
	s *store.Client,
) timelines.RepositoryAnalyticsResolver {
	return &analyticsResolver{database, s, l.Named("analytics")}
}

func (a *analyticsResolver) Burndown(
	ctx context.Context,
	repo *models.RepositoryAnalytics,
	t *models.BurndownType,
) (models.Burndown, error) {
	if t == nil {
		*t = models.BurndownTypeGlobal
	}
	id := repo.Repository.ID
	l := a.l.With("repo", id, "type", *t, log.LogKeyRID, log.HTTPRequestID(ctx))

	switch *t {
	case models.BurndownTypeFile:
		return nil, errors.New("unimplemented")
	case models.BurndownTypeAuthor:
		return nil, errors.New("unimplemented")
	case models.BurndownTypeGlobal:
		deltas, err := a.db.Repos().GetGlobalBurndown(ctx, id)
		if err != nil {
			if !db.IsNotFound(err) {
				l.Errorw(err.Error())
			}
			return nil, fmt.Errorf("could not find '%s' burndowns for repo '%d'", t, id)
		}
		return &models.GlobalBurndown{
			Entries: deltas,
		}, nil
	default:
		return nil, fmt.Errorf("invalid burndown type '%v'", t)
	}
}
