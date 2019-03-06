package server

import "github.com/bobheadxi/projector/graphql/go/projector"

// Resolver implements the Projector GraphQL API
type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Mutation() projector.MutationResolver {
	// TODO
	return nil
}

func (r *Resolver) Query() projector.QueryResolver {
	// TODO
	return nil
}
