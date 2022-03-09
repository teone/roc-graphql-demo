package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/onosproject/aether-roc-api/pkg/utils"
	externalRef1 "github.com/onosproject/config-models/modelplugin/aether-2.0.0/aether_2_0_0"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/teone/roc-graphql-demo/internal/gnmiConverter"
	"time"

	"github.com/teone/roc-graphql-demo/graph/generated"
	"github.com/teone/roc-graphql-demo/graph/model"

)

var log = logging.GetLogger("resolvers")

const enterprisePath = "/aether/v2.0.0/{target}/enterprises"

func (r *queryResolver) Enterprises(ctx context.Context) ([]*model.Enterprise, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)

	defer cancel()

	gnmiGet, err := utils.NewGnmiGetRequest(enterprisePath, string(r.target))
	if err != nil {
		return nil, err
	}
	log.Info("gnmi-get-request", "req", gnmiGet.String())
	gnmiVal, err := utils.GetResponseUpdate(r.client.Get(ctx, gnmiGet))

	if err != nil {
		log.Errorw("gnmi-get-error", "err", err)
		return nil, err
	}

	if gnmiVal == nil {
		return nil, nil
	}

	gnmiJsonVal, ok := gnmiVal.Value.(*gnmi.TypedValue_JsonVal)
	if !ok {
		return nil, fmt.Errorf("unexpected-type-of-reply-from-server: %v", gnmiVal.Value)
	}

	var gnmiResponse externalRef1.Device
	if err = externalRef1.Unmarshal(gnmiJsonVal.JsonVal, &gnmiResponse); err != nil {
		return nil, fmt.Errorf("error-unmarshalling-gnmiResponse: %v", err)
	}

	mpd := gnmiConverter.ModelPluginDevice{
		Device: gnmiResponse,
	}

	rocEnterprises, err := mpd.ToEnterprises()
	if err != nil {
		return nil, fmt.Errorf("error-casting-gnmiResponse-to-enterprise: %v", err)
	}

	log.Infow("received-gnmi-get-response", "rocEnterprises", rocEnterprises)

	enterprises := []*model.Enterprise{}
	for _, re := range *rocEnterprises.Enterprise {
		e := &model.Enterprise{
			ID: re.EnterpriseId,
			Name: re.DisplayName,
			Description: re.Description,
		}
		enterprises = append(enterprises, e)
	}

	return enterprises, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
