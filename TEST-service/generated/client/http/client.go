// Package http provides an HTTP client for the TestService service.

package http

import (
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	// This Service
	svc "github.com/hasAdamr/test-service/TEST-service/generated"
	handler "github.com/hasAdamr/test-service/TEST-service/handlers/server"
)

var (
	_ = endpoint.Chain
	_ = httptransport.NewClient
)

// New returns a service backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port".
func New(instance string) (handler.Service, error) {
	//options := []httptransport.ServerOption{
	//httptransport.ServerBefore(),
	//}

	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	_ = u

	var ReadContextTestValueZeroEndpoint endpoint.Endpoint
	{
		ReadContextTestValueZeroEndpoint = httptransport.NewClient(
			"get",
			copyURL(u, "/1"),
			svc.EncodeHTTPReadContextTestValueZeroRequest,
			svc.DecodeHTTPReadContextTestValueResponse,
			//options...,
		).Endpoint()
	}

	return svc.Endpoints{
		ReadContextTestValueEndpoint: ReadContextTestValueZeroEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
