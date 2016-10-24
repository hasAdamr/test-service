// Package http provides an HTTP client for the TestService service.

package http

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
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
	var cc *clientConfig

	for _, f := range options {
		_ = f(cc)
	}

	clientOptions := []httptransport.ClientOption{
		contextValuesToHttpHeaders(cc.headers),
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

// CtxValuesToSend will send specified keys as headers in http request with
// values that the context contains for those keys. Values must be strings.
func CtxValuesToSend(keys []string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys

		return nil
	}
}

func contextValuesToHttpHeaders(keys []string) httptransport.ClientOption {
	return httptransport.ClientBefore(
		func(ctx context.Context, r *http.Request) context.Context {
			for _, k := range keys {
				if v, ok := ctx.Value(k).(string); ok {
					r.Header[k] = append(r.Header[k], v)
				}
			}

			return ctx
		})
}
