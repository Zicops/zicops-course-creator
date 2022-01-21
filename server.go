package main

import (
	"context"

	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/config"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph"
	"github.com/zicops/zicops-course-creator/graph/generated"
	"github.com/zicops/zicops-course-creator/lib/db/cassandra"
)

const defaultPort = "8080"

func main() {
	log.Infof("Starting zicops course creator service")
	ctx, cancel := context.WithCancel(context.Background())
	cassConfig := config.NewCassandraConfig()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cassSession, err := cassandra.New(cassConfig)
	if err != nil {
		log.Errorf("Error connecting to cassandra: %s", err)
		log.Infof("zicops course creator intialization failed")
	}

	global.CTX = ctx
	global.CassSession = cassSession
	global.Cancel = cancel
	log.Infof("zicops course creator intialization complete")
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
