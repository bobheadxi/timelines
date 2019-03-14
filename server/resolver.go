package server

import "github.com/bobheadxi/timelines/graphql/go/timelines"

// resolver implements the timelines GraphQL API
type resolver struct{}

func newResolver() *resolver {
	return &resolver{}
}

func (r *resolver) Mutation() timelines.MutationResolver {
	// TODO
	return nil
}

func (r *resolver) Query() timelines.QueryResolver {
	// TODO
	return nil
}
