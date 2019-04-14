package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bobheadxi/timelines/graphql/go/timelines/models"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/analysis"
	"github.com/bobheadxi/timelines/host"
)

const (
	tableGitBurndownGlobals      = "git_burndowns_globals"
	tableGitBurndownFiles        = "git_burndowns_files"
	tableGitBurndownContributors = "git_burndowns_contributors"
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

// GetRepositories fetches all repositories associated with the given owner
func (r *ReposDatabase) GetRepositories(ctx context.Context, owner string) ([]models.Repository, error) {
	rows, err := r.db.pg.QueryEx(ctx, `
		SELECT
			id, owner, name
		FROM
			repositories
		WHERE
			owner=$1`,
		&pgx.QueryExOptions{},
		owner)
	if err != nil {
		return nil, err
	}
	var repos = make([]models.Repository, 0)
	for rows.Next() {
		var repo models.Repository
		if err := rows.Scan(&repo.ID, &repo.Owner, &repo.Name); err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

// GetRepository fetches a specific repository
func (r *ReposDatabase) GetRepository(ctx context.Context, owner, name string) (models.Repository, error) {
	row := r.db.pg.QueryRowEx(ctx, `
		SELECT
			id, owner, name
		FROM
			repositories
		WHERE
			owner=$1 AND name=$2`,
		&pgx.QueryExOptions{},
		owner, name)
	var repo models.Repository
	return repo, row.Scan(&repo.ID, &repo.Owner, &repo.Name)
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
	} else {
		r.l.Errorw("error encountered when creating repo",
			"error", err)
		if strings.Contains(err.Error(), "violates unique constraint") {
			return fmt.Errorf("repo '%s:%s/%s' already exists",
				h, owner, name)
		}
	}
	return err
}

// DeleteRepository removes a repository and associated items
func (r *ReposDatabase) DeleteRepository(ctx context.Context, id int) error {
	res, err := r.db.pg.ExecEx(ctx, `
		DELETE FROM 
			repositories
		WHERE
			id = $1
		`, &pgx.QueryExOptions{},
		id)
	if err != nil {
		return err
	}
	if res.RowsAffected() < 1 {
		return errors.New("no repository was deleted")
	}
	return nil
}

// DropGitBurndownResults deletes all git burndown results associated with the
// given repository ID.
// TODO: should this be done? see https://github.com/bobheadxi/timelines/issues/44
func (r *ReposDatabase) DropGitBurndownResults(ctx context.Context, repoID int) error {
	var (
		start = time.Now()
		l     = r.l.Named("drop_burndowns").With("repo_id", repoID)
	)

	tx, err := r.db.pg.BeginEx(ctx, &pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	for _, t := range []string{
		tableGitBurndownGlobals,
		tableGitBurndownFiles,
		tableGitBurndownContributors,
	} {
		// need to custom format the table name into the query string, since the
		// pgx formatter doesn't like table names as arguments
		_, err := tx.ExecEx(ctx, fmt.Sprintf(`
			DELETE FROM
				%s
			WHERE
				fk_repo_id = $1
			`, t), &pgx.QueryExOptions{},
			repoID)
		if err != nil {
			l.Errorw("could not drop burndowns from table",
				"table", t,
				"error", err)
			if err := tx.RollbackEx(ctx); err != nil {
				l.Warnw("error when rolling back transaction",
					"error", err)
			}
			return fmt.Errorf("could not delete existing burndowns, rolling back: %v", err)
		}
	}

	if err := tx.CommitEx(ctx); err != nil {
		l.Errorw("could not commit burndowns drops",
			"error", err)
		return fmt.Errorf("could not delete existing burndowns: %v", err)
	}

	l.Infow("burndown results dropped",
		"duration", time.Since(start))
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
	if count, err := tx.CopyFrom(
		pgx.Identifier{tableGitBurndownGlobals},
		[]string{
			"fk_repo_id", "interval",
			"delta_bands",
		},
		copyFromBurndowns(repoID, "", m, bd.Global),
	); err != nil {
		l.Errorw("failed to insert global burndowns",
			"error", err)
		return err
	} else if count != width {
		l.Warnf("expected '%d' items, got '%d' items", width, count)
	}

	l.Debug("preparing per-file burndowns")
	for f, v := range bd.Files {
		if count, err := tx.CopyFrom(
			pgx.Identifier{tableGitBurndownFiles},
			[]string{
				"fk_repo_id", "interval", "filename",
				"delta_bands",
			},
			copyFromBurndowns(repoID, f, m, v),
		); err != nil {
			l.Errorw("failed to insert file burndowns",
				"file", f,
				"error", err)
			return err
		} else if count != width {
			l.Warnf("file '%s': expected '%d' items, got '%d' items", f, width, count)
		}
	}

	l.Debug("preparing per-contributor burndowns")
	for c, v := range bd.Files {
		if count, err := tx.CopyFrom(
			pgx.Identifier{"git_burndowns_contributors"},
			[]string{
				"fk_repo_id", "interval", "contributor",
				"delta_bands",
			},
			copyFromBurndowns(repoID, c, m, v),
		); err != nil {
			l.Errorw("failed to insert contributor burndowns",
				"contributor", c,
				"error", err)
			return err
		} else if count != width {
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

// GetGlobalBurndown retrieves global burndowns for the given repo
// TODO: one at a time? all at once?
func (r *ReposDatabase) GetGlobalBurndown(
	ctx context.Context,
	repoID int,
) ([]models.BurndownEntry, error) {
	rows, err := r.db.pg.Query(fmt.Sprintf(`
	SELECT
		interval, delta_bands
	FROM
		%s
	WHERE
		fk_repo_id = $1
	`, tableGitBurndownGlobals), repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]models.BurndownEntry, 0)
	for rows.Next() {
		var (
			interval   pgtype.Tsrange
			deltaBands []int64
		)
		if err = rows.Scan(&interval, &deltaBands); err != nil {
			return nil, err
		}
		// for now, we need to cast manually into a int[]
		intBands := make([]int, len(deltaBands))
		for i, v := range deltaBands {
			intBands[i] = int(v)
		}
		entries = append(entries, models.BurndownEntry{
			Start: interval.Lower.Time,
			Bands: intBands,
		})
	}

	return entries, nil
}
