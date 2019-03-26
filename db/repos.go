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
	r.db.pg.Prepare(preparedStmtInsertHostItem, `
INSERT INTO
	host_items
VALUES
	(
		$1::INTEGER, $2::host_item_type, $3::INTEGER, $4::INTEGER,
		$5::TEXT, $6::DATE, $7::DATE, 
		$8::TEXT, $9::TEXT,
		$10::TEXT[], $11::JSONB, $12::JSONB
	)
`)
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
		batch     = r.db.pg.BeginBatch()
		itemCount int64
	)

	// queue all items for insertion
	for _, i := range items {
		if i == nil {
			break
		}
		itemCount++
		batch.Queue(preparedStmtInsertHostItem,
			[]interface{}{
				repoID, string(i.Type), i.GitHubID, i.Number,
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
		// TODO: this check does terrible things sometimes
		// https://github.com/bobheadxi/timelines/issues/22

		// r.l.Infow("provided items", "items", items[29].Number)
		// return fmt.Errorf("expected %d rows to change, only changed %d rows",
		//	itemCount, res.RowsAffected())
	}
	return nil
}
