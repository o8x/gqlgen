package main

import (
	"net/http"

	"github.com/99designs/gqlgen/_examples/sse/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func NewServer() *handler.Server {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))
	srv.AddTransport(transport.SSE{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}

func main() {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", NewServer())

	if err := http.ListenAndServe("0.0.0.0:5001", nil); err != nil {
		panic(err)
	}
}
