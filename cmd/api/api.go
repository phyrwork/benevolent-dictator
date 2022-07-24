package main

import (
	"fmt"
	"github.com/phyrwork/benevolent-dictator/pkg/api/auth"
	"github.com/phyrwork/benevolent-dictator/pkg/api/database"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/phyrwork/benevolent-dictator/pkg/api/graph"
	"github.com/phyrwork/benevolent-dictator/pkg/api/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dsn := fmt.Sprintf(
		"host=%s user=dictator password=dictator dbname=dictator sslmode=disable",
		os.Getenv("DB_SERVICE_HOST"))
	log.Printf("database dsn=%s", dsn)

	var err error
	var db *gorm.DB
	for db == nil {
		db, err = database.Open(dsn)
		if db == nil {
			break
		}
		log.Printf("database open error: %v", err)
		time.Sleep(time.Second * 3)
	}
	log.Printf("database open ok")
	if err = database.Migrate(db); err != nil {
		log.Fatalf("database migrate error: %v", err)
	}
	log.Printf("database migration ok")

	resolver := &graph.Resolver{
		DB: db,
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	mux := http.NewServeMux()
	mux.Handle("/query", auth.Handle(srv))
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
