package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"main/db"
	"main/graph"
	comment "main/packages/comment/repo/inMemory"
	commPostgres "main/packages/comment/repo/postgres"
	post "main/packages/post/repo/inMemory"
	postPostgres "main/packages/post/repo/postgres"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}
	var resolver *graph.Resolver
	storageType := os.Getenv("STORAGE")
	if storageType == "postgres" {
		postgresConfig := db.Config{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Dbname:   os.Getenv("POSTGRES_DBNAME"),
		}

		postgresConn, err := db.GetPostgresConnection(postgresConfig) // TODO
		if err != nil {
			log.Fatalf(err.Error())
		}

		postStorage := postPostgres.NewPostgresRepo(postgresConn)
		commStorage := commPostgres.NewPostgresRepo(postgresConn)
		resolver = graph.NewResolver(postStorage, commStorage)

	} else {
		postStorage := post.NewPostMemoryRepo()
		commStorage := comment.NewCommentMemoryRepo()
		resolver = graph.NewResolver(postStorage, commStorage)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
