package db

import "github.com/jackc/pgx"

// copyFromRows returns a CopyFromSource interface over the provided rows slice
// making it usable by *Conn.CopyFrom.
func copyFromRows(vals [][]interface{}) pgx.CopyFromSource {
	return &rows{rows: vals, idx: -1}
}

type rows struct {
	rows [][]interface{}
	idx  int
}

func (r *rows) Next() bool {
	r.idx++
	return r.idx < len(r.rows) && r.rows[r.idx] != nil
}

func (r *rows) Values() ([]interface{}, error) {
	return r.rows[r.idx], nil
}

func (r *rows) Err() error {
	return nil
}
