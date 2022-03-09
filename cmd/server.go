package main

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/teone/roc-graphql-demo/internal/gnmi_southbound"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/teone/roc-graphql-demo/graph"
	"github.com/teone/roc-graphql-demo/graph/generated"
)

const defaultPort = "8080"
const onosConfigAddress = "localhost:5150"
const target = "connectivity-service-v2"

var log = logging.GetLogger("server")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	gnmiClient, err := gnmi_southbound.NewGnmiClient(onosConfigAddress)
	if err != nil {
		log.Fatalw("cannot-start-gnmi-client", "err", err)
	}

	resolver, err := graph.NewResolver(gnmiClient, target)
	if err != nil {
		log.Fatalw("cannot-create-resolver", "err", err)
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
