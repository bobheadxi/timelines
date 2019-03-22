package server

import (
	"github.com/bobheadxi/timelines/graphql/go/timelines"
	"go.uber.org/zap"
)

// resolver implements the timelines GraphQL API
type resolver struct {
	l *zap.SugaredLogger
}

func newResolver(l *zap.SugaredLogger) *resolver {
	return &resolver{l}
}

func (r *resolver) Mutation() timelines.MutationResolver {
	// TODO
	return nil
}

func (r *resolver) Query() timelines.QueryResolver {
	// TODO
	return nil
}
