package host

import (
	"fmt"
	"strconv"
)

// Repo represents a repository
type Repo interface {
	Hosted
	GetOwner() string
	GetName() string
	IsPrivate() bool
	fmt.Stringer
}

// BaseRepo is an implementation of Repo, should probably be only used for testing
type BaseRepo struct {
	Host    Host
	ID      int
	Owner   string
	Name    string
	Private bool
}

// GetHost returns the repo's host
func (b *BaseRepo) GetHost() Host { return b.Host }

// GetID returns the repo's ID
func (b *BaseRepo) GetID() string { return strconv.Itoa(b.ID) }

// GetOwner returns the repo's owner
func (b *BaseRepo) GetOwner() string { return b.Owner }

// GetName returns the repo's name
func (b *BaseRepo) GetName() string { return b.Name }

// IsPrivate indicates whether the repo is private of not
func (b *BaseRepo) IsPrivate() bool { return b.Private }

func (b *BaseRepo) String() string { return fmt.Sprintf("%s:%s/%s", b.Host, b.Owner, b.Name) }

// doIstillImplementRepo is a dumb way of checking that BaseRepo is up to date
// with the Repo interface
func (b *BaseRepo) doIstillImplementRepo() Repo { return b }
