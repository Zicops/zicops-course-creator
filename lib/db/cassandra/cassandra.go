package cassandra

import (
	gocql "github.com/gocql/gocql"
	gocqlx "github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/config"
)

// Cassandra struct
type Cassandra struct {
	Session *gocqlx.Session
}

// New cassandra session and return Cassandra struct
func New(conf *config.Cassandra) (*Cassandra, error) {
	cluster := gocql.NewCluster(conf.Host)
	cluster.Keyspace = conf.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.Username,
		Password: conf.Password,
	}
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		// log via logrus
		log.Errorf("Error creating cassandra session: %v", err)
		return nil, err
	}
	return &Cassandra{
		Session: &session,
	}, nil
}
