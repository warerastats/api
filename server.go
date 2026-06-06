package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/warerastats/api/graph"
	"github.com/warerastats/api/graph/loaders"
	"github.com/warerastats/models/models"
)

const defaultPort = "8080"

// Public-playground guard limits.
const (
	playgroundMaxDepth     = 7
	playgroundMaxPageSize  = 50
	playgroundMaxQueryCost = 50000
	playgroundComplexity   = 2000
	playgroundRateInterval = 10 * time.Second
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx := context.Background()
	colls, err := models.Init(ctx)
	if err != nil {
		log.Fatalf("models init: %v", err)
	}
	defer colls.Close(ctx)

	es := graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Colls: colls}})

	// Authenticated endpoint: trusted clients (SvelteKit frontend, future API
	// key holders). No depth/complexity/time-window restrictions.
	srv := newGraphQLServer(es)

	// Public playground endpoint: open to guests but throttled and bounded so the
	// whole dataset cannot be scraped.
	playgroundSrv := newGraphQLServer(es)
	playgroundSrv.Use(extension.FixedComplexityLimit(playgroundComplexity))
	playgroundSrv.AroundOperations(graph.PlaygroundGuard(playgroundMaxDepth, playgroundMaxPageSize, playgroundMaxQueryCost))

	limiter := newIPRateLimiter(playgroundRateInterval)

	mux := http.NewServeMux()

	// Open playground UI; its queries are sent to the throttled public endpoint.
	mux.Handle("/", playground.Handler("GraphQL playground", "/playground/query"))

	// Public, rate-limited query endpoint backing the playground.
	mux.Handle("/playground/query", cors(limiter.middleware(
		playgroundContext(loaders.Middleware(colls, playgroundSrv)),
	)))

	// Authenticated query endpoint used by the frontend and API consumers.
	mux.Handle("/query", cors(apiKeyAuth(loaders.Middleware(colls, srv))))

	log.Printf("connect to http://localhost:%s/ for the public GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// newGraphQLServer builds a gqlgen server with the shared transports, query
// cache and introspection enabled.
func newGraphQLServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}
