package global

import (
	"context"
	"sync"

	gocqlx "github.com/scylladb/gocqlx/v2"
)

// some global variables commonly used
var (
	CTX             context.Context
	CassSession     *gocqlx.Session
	CassSessioQBank *gocqlx.Session
	Cancel          context.CancelFunc
	WaitGroupServer sync.WaitGroup
)

// initializes global package to read environment variables as needed
func init() {
}
