package server

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/graphql/go/timelines"
	"github.com/bobheadxi/timelines/graphql/go/timelines/models"
	"github.com/bobheadxi/timelines/store"
)

type combinedAnalyticsResolver interface {
	timelines.FileBurndownResolver
	timelines.RepositoryAnalyticsResolver
}

type analyticsResolver struct {
	db *db.Database
	s  *store.Client

	l *zap.SugaredLogger
}

func newAnalyticsResolver(
	l *zap.SugaredLogger,
	database *db.Database,
	s *store.Client,
) combinedAnalyticsResolver {
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
	l := a.l.With("repo", id, "type", *t, config.LogKeyRID, requestID(ctx))

	switch *t {
	case models.BurndownTypeFile:
		return &models.FileBurndown{
			RepoID: id,
			Type:   models.BurndownTypeFile,
			// Entries to be populated by FileResolver
		}, nil
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
			RepoID:  id,
			Type:    models.BurndownTypeGlobal,
			Entries: deltas,
		}, nil
	default:
		return nil, fmt.Errorf("invalid burndown type '%v'", t)
	}
}

func (a *analyticsResolver) File(
	ctx context.Context,
	bd *models.FileBurndown,
	filename *string,
) ([]*models.FileBurndownEntry, error) {
	var deltas []*models.FileBurndownEntry
	var err error
	if filename == nil {
		deltas, err = a.db.Repos().GetFilesBurndown(ctx, bd.RepoID, "")
		if err != nil {
			return nil, fmt.Errorf("could not find burndowns repo '%d': %v", bd.RepoID, err)
		}
	} else {
		deltas, err = a.db.Repos().GetFilesBurndown(ctx, bd.RepoID, *filename)
		if err != nil {
			return nil, fmt.Errorf("could not find burndowns for file '%s' in repo '%d': %v", *filename, bd.RepoID, err)
		}
	}
	return deltas, nil
}
