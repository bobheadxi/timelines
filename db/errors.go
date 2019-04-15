package db

import (
	"errors"
	"strings"
)

const (
	errMsgNotFound = "no results found"
)

// IsNotFound checks if the given errors is a database result-not-found error
func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), errMsgNotFound)
}

func errNotFound() error { return errors.New(errMsgNotFound) }

func isPgxNotFound(err error) bool {
	return strings.Contains(err.Error(), "no rows")
}
