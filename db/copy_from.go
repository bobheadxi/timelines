package db

import (
	"github.com/jackc/pgx/pgtype"

	"github.com/bobheadxi/timelines/analysis"
	"github.com/bobheadxi/timelines/host"
)

type rowsItems struct {
	repo int
	rows []*host.Item
	idx  int
}

func copyFromItems(repo int, items []*host.Item) *rowsItems {
	return &rowsItems{repo: repo, rows: items, idx: -1}
}

func (r *rowsItems) Next() bool {
	r.idx++
	return r.idx < len(r.rows) && r.rows[r.idx] != nil
}

func (r *rowsItems) Values() ([]interface{}, error) {
	v := r.rows[r.idx]
	return []interface{}{
		r.repo, string(v.Type), v.GitHubID, v.Number,
		v.Author, v.Opened, v.Closed,
		v.Title, v.Body,
		v.Labels, v.Reactions, v.Details,
	}, nil
}

func (r *rowsItems) Err() error { return nil }

type rowsBurndowns struct {
	repo int
	name string
	meta *analysis.GitRepoMeta
	data [][]int64
	idx  int
}

func copyFromBurndowns(
	repo int, name string, m *analysis.GitRepoMeta, data [][]int64,
) *rowsBurndowns {
	return &rowsBurndowns{repo: repo, name: name, meta: m, data: data, idx: -1}
}

func (r *rowsBurndowns) Next() bool {
	r.idx++
	return r.idx < len(r.data) && r.data[r.idx] != nil
}

func (r *rowsBurndowns) Values() ([]interface{}, error) {
	start, end := r.meta.GetTickRange(r.idx)
	ts := &pgtype.Tsrange{
		// pg timestamps must be UTC
		Lower:     pgtype.Timestamp{Time: start.UTC(), Status: pgtype.Present},
		Upper:     pgtype.Timestamp{Time: end.UTC(), Status: pgtype.Present},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Status:    pgtype.Present,
	}
	if r.name != "" {
		return []interface{}{
			r.repo, ts, r.name,
			r.data[r.idx],
		}, nil
	}
	return []interface{}{
		r.repo, ts,
		r.data[r.idx],
	}, nil
}

func (r *rowsBurndowns) Err() error { return nil }
