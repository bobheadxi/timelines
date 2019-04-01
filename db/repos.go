package db

import (
	"context"
	"errors"

	// need https://github.com/jackc/pgx/issues/335
	"github.com/jackc/pgx"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/analysis"
	"github.com/bobheadxi/timelines/host"
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
	preparedStmtInsertHostItem          = "insert_host_item"
	preparedStmtInsertGitBurndownGlobal = "insert_git_burndowns_globals"
	preparedStmtInsertGitBurndownFile   = "insert_git_burndowns_files"
	preparedStmtInsertGitBurndownPeople = "insert_git_burndowns_contributors"
)

// init sets up all prepared statements associated with repositories
func (r *ReposDatabase) init() {
	r.db.pg.Prepare(preparedStmtInsertGitBurndownGlobal, `
INSERT INTO
	git_burndowns_globals
VALUES
	(
		$1::INTEGER, $2::TSRANGE, $3::INTEGER[]
	)
`)
	r.db.pg.Prepare(preparedStmtInsertGitBurndownFile, `
INSERT INTO
	git_burndowns_files
VALUES
	(
		$1::INTEGER, $2::TEXT, $3::TSRANGE, $4::INTEGER[]
	)
`)
	r.db.pg.Prepare(preparedStmtInsertGitBurndownPeople, `
INSERT INTO
	git_burndowns_contributors
VALUES
	(
		$1::INTEGER, $2::TEXT, $3::TSRANGE, $4::INTEGER[]
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
func (r *ReposDatabase) NewRepository(
	ctx context.Context,
	h host.Host,
	installation, owner, name string,
) error {
	if installation == "" {
		return errors.New("installation required")
	}
	if owner == "" || name == "" {
		return errors.New("repository identifiers (owner and name) required")
	}
	_, err := r.db.pg.ExecEx(ctx, `
	INSERT INTO 
		repositories (installation_id, type, owner, name)
	VALUES
		($1, $2, $3, $4)
	`, &pgx.QueryExOptions{},
		installation, string(h), owner, name)
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
func (r *ReposDatabase) InsertGitBurndownResult(ctx context.Context, burndown *analysis.BurndownResult) error {
	// TODO
	return nil
}

// InsertHostItems executes a batch insert on all given items
func (r *ReposDatabase) InsertHostItems(ctx context.Context, repoID int, items []*host.Item) error {
	var (
		itemCount int
		rows      = make([][]interface{}, len(items))
		l         = r.l.With("repo_id", repoID)
	)

	for i, v := range items {
		if v == nil {
			break
		}
		itemCount++
		rows[i] = []interface{}{
			repoID, string(v.Type), v.GitHubID, v.Number,
			v.Author, v.Opened, v.Closed,
			v.Title, v.Body,
			v.Labels, v.Reactions, v.Details,
		}
	}
	l = l.With("items", itemCount)
	l.Info("preparing to insert items")

	count, err := r.db.pg.CopyFrom(
		pgx.Identifier{"host_items"},
		[]string{
			"fk_repo_id", "type", "host_id", "number",
			"author", "open_date", "close_date",
			"title", "body",
			"labels", "reactions", "details",
		},
		copyFromRows(rows),
	)
	if err != nil {
		l.Errorw("failed to insert host items",
			"error", err)
		return err
	}
	if count != itemCount {
		l.Errorf("expected '%d' items, got '%d' items", itemCount, count)
		return errors.New("unexpected mismatch actual items and inserted items")
	}
	l.Infow("items successfully inserted")
	return nil
}
