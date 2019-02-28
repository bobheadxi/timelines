package store

import (
	"github.com/blevesearch/bleve"
	"go.uber.org/zap"
)

type Store struct {
	i bleve.Index

	l *zap.Logger
}
