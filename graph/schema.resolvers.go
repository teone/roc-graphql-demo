package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/teone/roc-graphql-demo/graph/generated"
	"github.com/teone/roc-graphql-demo/graph/model"
	"github.com/teone/roc-graphql-demo/internal/stores"
)

func (r *queryResolver) Enterprises(ctx context.Context) ([]*model.Enterprise, error) {
	rocEnterprises, err := stores.ListEnterprises(ctx, r.client, r.target)
	if err != nil {
		return nil, err
	}

	enterprises := []*model.Enterprise{}
	for _, re := range *rocEnterprises.Enterprise {
		apps := []*model.Application{}
		for _, ra := range *re.Application {
			app := &model.Application{
				ID:   ra.ApplicationId,
				Name: ra.DisplayName,
			}
			apps = append(apps, app)
		}

		e := &model.Enterprise{
			ID:           re.EnterpriseId,
			Name:         re.DisplayName,
			Description:  re.Description,
			Applications: apps,
			Sites:        stores.ListSites(re),
		}
		enterprises = append(enterprises, e)
	}

	return enterprises, nil
}

func (r *queryResolver) Site(ctx context.Context) ([]*model.Site, error) {
	rocEnterprises, err := stores.ListEnterprises(ctx, r.client, r.target)
	if err != nil {
		return nil, err
	}
	sites := []*model.Site{}
	for _, re := range *rocEnterprises.Enterprise {
		sites = append(sites, stores.ListSites(re)...)
	}
	return sites, nil
}

func (r *queryResolver) SimCards(ctx context.Context) ([]*model.SimCard, error) {
	rocEnterprises, err := stores.ListEnterprises(ctx, r.client, r.target)
	if err != nil {
		return nil, err
	}
	simCards := []*model.SimCard{}
	for _, re := range *rocEnterprises.Enterprise {
		for _, rs := range *re.Site {
			simCards = append(simCards, stores.ListSimCard(rs)...)
		}
	}
	return simCards, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
