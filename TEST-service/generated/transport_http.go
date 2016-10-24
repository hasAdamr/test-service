package svc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"

	// This service
	pb "github.com/hasAdamr/test-service/TEST-service"
)

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.RegisterTestServiceServer
	_ = io.Copy
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandler(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
	}
	m := http.NewServeMux()

	m.Handle("/1", httptransport.NewServer(
		ctx,
		endpoints.ReadContextTestValueEndpoint,
		HttpDecodeLogger(DecodeHTTPReadContextTestValueZeroRequest, logger),
		EncodeHTTPGenericResponse,
		options...,
	))
	return m
}

func HttpDecodeLogger(next httptransport.DecodeRequestFunc, logger log.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		logger.Log("method", r.Method, "url", r.URL.String())
		rv, err := next(ctx, r)
		if err != nil {
			logger.Log("method", r.Method, "url", r.URL.String(), "Error", err)
		}
		return rv, err
	}
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// Server Decode

// DecodeHTTPReadContextTestValueZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded readcontexttestvalue request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPReadContextTestValueZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.EmptyMessage
	err := json.NewDecoder(r.Body).Decode(&req)
	// err = io.EOF if r.Body was empty
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "decoding body of http request")
	}

	pathParams, err := PathParams(r.URL.Path, "/1")
	_ = pathParams
	if err != nil {
		fmt.Printf("Error while reading path params: %v\n", err)
		return nil, errors.Wrap(err, "couldn't unmarshal path parameters")
	}
	queryParams, err := QueryParams(r.URL.Query())
	_ = queryParams
	if err != nil {
		fmt.Printf("Error while reading query params: %v\n", err)
		return nil, errors.Wrapf(err, "Error while reading query params: %v", r.URL.Query())
	}

	return &req, err
}

// Client Decode

// DecodeHTTPReadContextTestValue is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded EmptyMessage response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPReadContextTestValueResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp pb.EmptyMessage
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

// Client Encode

// EncodeHTTPReadContextTestValueZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a readcontexttestvalue request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPReadContextTestValueZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	fmt.Printf("Encoding request %v\n", request)
	req := request.(*pb.EmptyMessage)
	_ = req

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"1",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()

	// Set the body parameters
	var buf bytes.Buffer
	toRet := map[string]interface{}{}
	if err := json.NewEncoder(&buf).Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	fmt.Printf("URL: %v\n", r.URL)
	return nil
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// Helper functions

// PathParams takes a url and a gRPC-annotation style url template, and
// returns a map of the named parameters in the template and their values in
// the given url.
//
// PathParams does not support the entirety of the URL template syntax defined
// in third_party/googleapis/google/api/httprule.proto. Only a small subset of
// the functionality defined there is implemented here.
func PathParams(url string, urlTmpl string) (map[string]string, error) {
	rv := map[string]string{}
	pmp := BuildParamMap(urlTmpl)

	parts := strings.Split(url, "/")
	for k, v := range pmp {
		rv[k] = parts[v]
	}

	return rv, nil
}

// BuildParamMap takes a string representing a url template and returns a map
// indicating the location of each parameter within that url, where the
// location is the index as if in a slash-separated sequence of path
// components. For example, given the url template:
//
//     "/v1/{a}/{b}"
//
// The returned param map would look like:
//
//     map[string]int {
//         "a": 2,
//         "b": 3,
//     }
func BuildParamMap(urlTmpl string) map[string]int {
	rv := map[string]int{}

	parts := strings.Split(urlTmpl, "/")
	for idx, part := range parts {
		if strings.ContainsAny(part, "{}") {
			param := RemoveBraces(part)
			rv[param] = idx
		}
	}
	return rv
}

// RemoveBraces replace all curly braces in the provided string, opening and
// closing, with empty strings.
func RemoveBraces(val string) string {
	val = strings.Replace(val, "{", "", -1)
	val = strings.Replace(val, "}", "", -1)
	return val
}

// QueryParams takes query parameters in the form of url.Values, and returns a
// bare map of the string representation of each key to the string
// representation for each value. The representations of repeated query
// parameters is undefined.
func QueryParams(vals url.Values) (map[string]string, error) {

	rv := map[string]string{}
	for k, v := range vals {
		rv[k] = v[0]
	}
	return rv, nil
}

func headersToContext(ctx context.Context, r *http.Request) context.Context {
	for k, v := range r.Header {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx
}
