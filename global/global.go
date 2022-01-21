package global

import (
	"context"

	"github.com/zicops/zicops-course-creator/lib/db/cassandra"
)

// some global variables commonly used
var (
	CTX         context.Context
	CassSession *cassandra.Cassandra
	Cancel      context.CancelFunc
)

// initializes global package to read environment variables as needed
func init() {
}
