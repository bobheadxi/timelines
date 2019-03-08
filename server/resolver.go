package server

import "github.com/bobheadxi/projector/graphql/go/projector"

// resolver implements the Projector GraphQL API
type resolver struct{}

func newResolver() *resolver {
	return &resolver{}
}

func (r *resolver) Mutation() projector.MutationResolver {
	// TODO
	return nil
}

func (r *resolver) Query() projector.QueryResolver {
	// TODO
	return nil
}
