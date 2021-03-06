// Package http provides an HTTP client for the TestService service.
package http

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

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
func New(instance string, options ...ClientOption) (handler.Service, error) {
	var cc clientConfig

	for _, f := range options {
		err := f(&cc)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	clientOptions := []httptransport.ClientOption{
		httptransport.ClientBefore(
			contextValuesToHttpHeaders(cc.headers)),
	}

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
			clientOptions...,
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

type clientConfig struct {
	headers []string
}

// ClientOption is a function that modifies the client config
type ClientOption func(*clientConfig) error

// CtxValuesToSend configures the http client to pull the specified keys out of
// the context and add them to the http request as headers.  Note that keys
// will have net/http.CanonicalHeaderKey called on them before being send over
// the wire and that is the form they will be available in the server context.
func CtxValuesToSend(keys []string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys
		return nil
	}
}

func contextValuesToHttpHeaders(keys []string) httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				r.Header.Set(k, v)
			}
		}

		return ctx
	}
}
