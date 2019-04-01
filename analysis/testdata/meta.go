package testdata

import (
	"time"

	"github.com/bobheadxi/timelines/analysis"
)

// Meta is an example metadata set. It can't be included in the codegen because
// time.Time's '%#v' representation uses unexported fields, boo hoo.
var Meta = &analysis.GitRepoMeta{
	Commits:  96,
	First:    time.Unix(1509843192, 0),
	Last:     time.Unix(1522909723, 0),
	TickSize: 6,
}
