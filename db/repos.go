package db

import (
	"context"
	"errors"
	"time"

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

// init sets up all prepared statements associated with repositories
func (r *ReposDatabase) init() {
	// no-op for now
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
func (r *ReposDatabase) InsertGitBurndownResult(
	ctx context.Context,
	repoID int,
	m *analysis.GitRepoMeta,
	bd *analysis.BurndownResult,
) error {
	var (
		start = time.Now()
		width = len(bd.Global)
		l     = r.l.Named("insert_burndowns").With("repo_id", repoID, "items", width+
			len(bd.Files)*width+
			len(bd.People)*width)
	)

	tx, err := r.db.pg.BeginEx(ctx, &pgx.TxOptions{})
	if err != nil {
		return err
	}

	l.Debug("preparing global burndowns")
	count, err := tx.CopyFrom(
		pgx.Identifier{"git_burndowns_globals"},
		[]string{
			"fk_repo_id", "interval",
			"delta_bands",
		},
		copyFromBurndowns(repoID, "", m, bd.Global),
	)
	if err != nil {
		l.Errorw("failed to insert global burndowns",
			"error", err)
		return err
	}
	if count != width {
		l.Errorf("expected '%d' items, got '%d' items", width, count)
	}

	l.Debug("preparing per-file burndowns")
	for f, v := range bd.Files {
		count, err := tx.CopyFrom(
			pgx.Identifier{"git_burndowns_files"},
			[]string{
				"fk_repo_id", "interval", "filename",
				"delta_bands",
			},
			copyFromBurndowns(repoID, f, m, v),
		)
		if err != nil {
			l.Errorw("failed to insert file burndowns",
				"file", f,
				"error", err)
			return err
		}
		if count != width {
			l.Errorf("file '%s': expected '%d' items, got '%d' items", f, width, count)
		}
	}

	l.Debug("preparing per-contributor burndowns")
	for c, v := range bd.Files {
		count, err := tx.CopyFrom(
			pgx.Identifier{"git_burndowns_contributors"},
			[]string{
				"fk_repo_id", "interval", "contributor",
				"delta_bands",
			},
			copyFromBurndowns(repoID, c, m, v),
		)
		if err != nil {
			l.Errorw("failed to insert contributor burndowns",
				"contributor", c,
				"error", err)
			return err
		}
		if count != width {
			l.Errorf("contributor '%s': expected '%d' items, got '%d' items", c, width, count)
		}
	}

	if err := tx.CommitEx(ctx); err != nil {
		l.Errorw("failed to commit transaction for burndown insertion",
			"error", err)
		return err
	}

	l.Infow("burndown committed",
		"duration", time.Since(start))
	return nil
}

// InsertHostItems executes a batch insert on all given items
func (r *ReposDatabase) InsertHostItems(
	ctx context.Context,
	repoID int,
	items []*host.Item,
) error {
	var (
		l     = r.l.Named("insert_host_items").With("repo_id", repoID)
		cp    = copyFromItems(repoID, items)
		start = time.Now()
	)

	count, err := r.db.pg.CopyFrom(
		pgx.Identifier{"host_items"},
		[]string{
			"fk_repo_id", "type", "host_id", "number",
			"author", "open_date", "close_date",
			"title", "body",
			"labels", "reactions", "details",
		},
		cp,
	)
	if err != nil {
		l.Errorw("failed to insert host items",
			"error", err)
		return err
	}
	if count != cp.idx {
		l.Errorf("expected '%d' items, got '%d' items", cp.idx, count)
		return errors.New("unexpected mismatch actual items and inserted items")
	}

	l.Infow("items successfully inserted",
		"duration", time.Since(start))
	return nil
}
