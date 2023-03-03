package global

import (
	"context"
	"sync"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/zicops/zicops-cass-pool/cassandra"
)

// some global variables commonly used
var (
	CTX             context.Context
	Cancel          context.CancelFunc
	WaitGroupServer sync.WaitGroup
	App             *firebase.App
	Client          *firestore.Client
	Ct              = context.Background()
	CassPool        *cassandra.CassandraPool
)

// initializes global package to read environment variables as needed
func init() {
}
