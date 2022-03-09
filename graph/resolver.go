package graph

import (
	externalRef0 "github.com/onosproject/aether-roc-api/pkg/aether_2_0_0/types"
	"github.com/openconfig/gnmi/proto/gnmi"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	client     gnmi.GNMIClient
	target     externalRef0.Target
}

func NewResolver(client gnmi.GNMIClient, gnmiTarget string) (*Resolver, error) {

	log.Infow("creating-resolver", "client", client, "gnmiTarget", gnmiTarget)
	return &Resolver{
		client: client,
		target: externalRef0.Target(gnmiTarget),
	}, nil
}
