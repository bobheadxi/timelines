package server

import "github.com/bobheadxi/projector/graphql/golang/api"

// Resolver implements the Projector GraphQL API
type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Mutation() api.MutationResolver {
	// TODO
	return nil
}

func (r *Resolver) Query() api.QueryResolver {
	// TODO
	return nil
}
