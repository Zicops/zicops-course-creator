package global

import (
	"context"
	"sync"
)

// some global variables commonly used
var (
	CTX             context.Context
	Cancel          context.CancelFunc
	WaitGroupServer sync.WaitGroup
)

// initializes global package to read environment variables as needed
func init() {
}
