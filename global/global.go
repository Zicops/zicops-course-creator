package global

import (
	"context"
	"sync"

	"github.com/zicops/zicops-course-creator/lib/db/cassandra"
)

// some global variables commonly used
var (
	CTX             context.Context
	CassSession     *cassandra.Cassandra
	CassSessioQBank *cassandra.Cassandra
	Cancel          context.CancelFunc
	WaitGroupServer sync.WaitGroup
)

// initializes global package to read environment variables as needed
func init() {
}
