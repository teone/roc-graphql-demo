package gnmi_southbound

import (
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/grpc/retry"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"os"
	"time"
)

var log = logging.GetLogger("GnmiClient")

func NewGnmiClient(address string) (gnmi.GNMIClient, error) {

	// TODO allow secure connections
	opts, err := certs.HandleCertPaths("", "", "", true)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	optsWithRetry := []grpc.DialOption{
		grpc.WithStreamInterceptor(retry.RetryingStreamClientInterceptor(retry.WithInterval(100 * time.Millisecond))),
	}
	optsWithRetry = append(opts, optsWithRetry...)
	gnmiConn, err := grpc.Dial(address, optsWithRetry...)
	if err != nil {
		log.Errorw("cannot-connect-to-gnmi-server", "err", err, "address", address)
		return nil, err
	}

	return gnmi.NewGNMIClient(gnmiConn), nil
}
