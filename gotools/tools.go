// +build tools

package tools

import (
	// golint is a linter for Go source code.
	_ "golang.org/x/lint/golint"

	// gqlgen is a tool for generating Go bindings for the GraphQL API.
	// gqlgen configuration is tracked in `graphql/gqlgen.yml`.
	_ "github.com/99designs/gqlgen"
)
