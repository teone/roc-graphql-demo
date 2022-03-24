package stores

import (
	"context"
	"fmt"
	"github.com/onosproject/aether-roc-api/pkg/aether_2_0_0/types"
	externalRef0 "github.com/onosproject/aether-roc-api/pkg/aether_2_0_0/types"
	"github.com/onosproject/aether-roc-api/pkg/utils"
	externalRef1 "github.com/onosproject/config-models/modelplugin/aether-2.0.0/aether_2_0_0"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/teone/roc-graphql-demo/graph/model"
	"github.com/teone/roc-graphql-demo/internal/gnmiConverter"
	"time"
)

var log = logging.GetLogger("stores")

const enterprisePath = "/aether/v2.0.0/{target}/enterprises"

func ListEnterprises(ctx context.Context, client gnmi.GNMIClient, target externalRef0.Target) (*types.Enterprises, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	gnmiGet, err := utils.NewGnmiGetRequest(enterprisePath, string(target))
	if err != nil {
		return nil, err
	}
	log.Info("gnmi-get-request", "req", gnmiGet.String())
	gnmiVal, err := utils.GetResponseUpdate(client.Get(ctx, gnmiGet))

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
	return rocEnterprises, nil
}

func ListSites(enterprise types.EnterprisesEnterprise) ([]*model.Site) {
	sites := []*model.Site{}
	for _, rs := range *enterprise.Site {

		alerts := 0
		switch rs.SiteId {
		case "acme-chicago":
			alerts = 2
		case "starbucks-seattle":
			alerts = 3
			
		}

		image := "https://chronos-dev.onlab.us/chronos-exporter/images/berlin-deutschland.png"
		
		simCards := ListSimCard(rs)
		s := &model.Site{
			ID:       rs.SiteId,
			Name:     rs.DisplayName,
			Devices:  ListDevice(rs),
			SimCards: simCards,
			SimCardsCount: len(simCards),
			Alerts: &alerts,
			Image: &image,
			Slices: ListSlices(rs),
		}
		sites = append(sites, s)
	}

	return sites
}

func ListDevice(site types.EnterprisesEnterpriseSite) ([]*model.Device) {
	devices := []*model.Device{}
	for _, rd := range *site.Device {
		d := &model.Device{
			ID:      rd.DeviceId,
			Name:    rd.DisplayName,
			SimCard: rd.SimCard,
		}
		devices = append(devices, d)
	}
	return devices
}

func ListSimCard(site types.EnterprisesEnterpriseSite) ([]*model.SimCard) {
	simCards := []*model.SimCard{}
	for _, sim := range *site.SimCard {
		d := &model.SimCard{
			ID:      sim.SimId,
			Name:    sim.DisplayName,
		}
		simCards = append(simCards, d)
	}
	return simCards
}

func ListSlices(site types.EnterprisesEnterpriseSite) ([]*model.Slices) {
	slices := []*model.Slices{}
	for _, slice := range *site.Slice {
		d := &model.Slices{
			ID:   slice.SliceId,
			Name: slice.DisplayName,
		}
		slices = append(slices, d)
	}
	return slices
}