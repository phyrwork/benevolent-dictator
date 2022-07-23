package main

import (
	"github.com/phyrwork/benevolent-dictator/pkg/api/auth"
	"github.com/phyrwork/benevolent-dictator/pkg/api/database"
	"log"
	"net/http"
	"os"

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

	db, err := database.Open("")
	if err != nil {
		log.Fatalf("database open error: %v", err)
	}
	if err = database.Migrate(db); err != nil {
		log.Fatalf("database migrate error: %v", err)
	}

	resolver := &graph.Resolver{
		DB: db,
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	mux := http.NewServeMux()
	mux.Handle("/query", auth.Middleware(srv))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
