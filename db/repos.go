package db

import (
	"context"
	"errors"
	"fmt"

	// need https://github.com/jackc/pgx/issues/335
	"github.com/bobheadxi/pgx"

	"go.uber.org/zap"

	"github.com/bobheadxi/projector/analysis"
	"github.com/bobheadxi/projector/github"
)

// Repository represents a stored repository. TODO
type Repository struct {
	ID           int
	Installation string

	GitBurndowns []*GitBurndown
}

// GitBurndown represents one sample per entry
type GitBurndown struct {
	ID     int
	RepoID int
	Month  string

	Global int
	Files  map[string]int
	People map[string]int
}

// ReposDatabase is a client for accessing repository-related databases
type ReposDatabase struct {
	db *Database
	l  *zap.SugaredLogger
}

const (
	preparedStmtInsertGitHubItem = "insert_github_item"
)

// init sets up all prepared statements associated with repositories
func (r *ReposDatabase) init() {
	r.db.pg.Prepare(preparedStmtInsertGitHubItem, `
INSERT INTO
	github_items
VALUES
	(
		$1::INTEGER, $2::INTEGER, $3::INTEGER, $4::github_item_type,
		$5::TEXT, $6::DATE, $7::DATE, 
		$8::TEXT, $9::TEXT,
		$10::TEXT[], $11::JSONB, $12::JSONB
	)
`)
}

// GetRepositoryID retrieves the ID associated with the given repository
func (r *ReposDatabase) GetRepositoryID(ctx context.Context, owner, name string) (int, error) {
	row := r.db.pg.QueryRowEx(ctx,
		"SELECT id FROM repositories WHERE owner=$1 AND name=$2",
		&pgx.QueryExOptions{},
		owner, name)
	var id int
	return id, row.Scan(&id)
}

// NewRepository creates a new repository entry
func (r *ReposDatabase) NewRepository(ctx context.Context, installation, owner, name string) error {
	if installation == "" {
		return errors.New("installation required")
	}
	if owner == "" || name == "" {
		return errors.New("repository identifiers (owner and name) required")
	}
	_, err := r.db.pg.ExecEx(ctx, `
	INSERT INTO 
		repositories (installation_id, owner, name)
	VALUES
		($1, $2, $3)
	`, &pgx.QueryExOptions{},
		installation, owner, name)
	if err == nil {
		r.l.Infow("created new entry for repo", "repo", owner+"/"+name)
	}
	return err
}

// DeleteRepository removes a repository and associated items
func (r *ReposDatabase) DeleteRepository(ctx context.Context, id int) error {
	res, err := r.db.pg.ExecEx(ctx, `
	DELETE FROM 
		repositories
	WHERE id = $1`, &pgx.QueryExOptions{}, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() < 1 {
		return errors.New("no repository was deleted")
	}
	return nil
}

// InsertGitBurndownResult processes a burndown analysis for insertion
func (r *ReposDatabase) InsertGitBurndownResult(ctx context.Context, burndown *analysis.BurndownResult) {
	// TODO
}

// InsertGitHubItems executes a batch insert on all given items
func (r *ReposDatabase) InsertGitHubItems(ctx context.Context, repoID int, items []*github.Item) error {
	var (
		batch     = r.db.pg.BeginBatch()
		itemCount int64
	)

	// queue all items for insertion
	for _, i := range items {
		if i == nil {
			break
		}
		itemCount++
		batch.Queue(preparedStmtInsertGitHubItem,
			[]interface{}{
				repoID, i.GitHubID, i.Number, i.Type,
				i.Author, i.Opened, i.Closed,
				i.Title, i.Body,
				i.Labels, i.Reactions, i.Details,
			}, nil, nil)
	}

	// send and fetch execution results
	if err := batch.Send(ctx, &pgx.TxOptions{}); err != nil {
		return err
	}
	res, err := batch.ExecResults()
	if err != nil {
		return err
	}

	// if an incorrect number of rows is modified, throw an error
	if res.RowsAffected() != itemCount {
		return fmt.Errorf("expected %d rows to change, only changed %d rows",
			itemCount, res.RowsAffected())
	}
	return nil
}
