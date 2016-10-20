// Package grpc provides a gRPC client for the add service.
package grpc

import (
	//"time"

	//jujuratelimit "github.com/juju/ratelimit"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	//"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/ratelimit"
	//"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/hasAdamr/test-service/TEST-service"
	svc "github.com/hasAdamr/test-service/TEST-service/generated"
	handler "github.com/hasAdamr/test-service/TEST-service/handlers/server"
)

// New returns an AddService backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn /*, tracer stdopentracing.Tracer, logger log.Logger*/) handler.Service {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.

	//limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

	var readcontexttestvalueEndpoint endpoint.Endpoint
	{
		readcontexttestvalueEndpoint = grpctransport.NewClient(
			conn,
			"TEST.TestService",
			"ReadContextTestValue",
			svc.EncodeGRPCReadContextTestValueRequest,
			svc.DecodeGRPCReadContextTestValueResponse,
			pb.EmptyMessage{},
			//grpctransport.ClientBefore(opentracing.FromGRPCRequest(tracer, "ReadContextTestValue", logger)),
		).Endpoint()
		//readcontexttestvalueEndpoint = opentracing.TraceClient(tracer, "ReadContextTestValue")(readcontexttestvalueEndpoint)
		//readcontexttestvalueEndpoint = limiter(readcontexttestvalueEndpoint)
		//readcontexttestvalueEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		//Name:    "ReadContextTestValue",
		//Timeout: 30 * time.Second,
		//}))(readcontexttestvalueEndpoint)
	}

	var readcontextmetadataEndpoint endpoint.Endpoint
	{
		readcontextmetadataEndpoint = grpctransport.NewClient(
			conn,
			"TEST.TestService",
			"ReadContextMetadata",
			svc.EncodeGRPCReadContextMetadataRequest,
			svc.DecodeGRPCReadContextMetadataResponse,
			pb.EmptyMessage{},
			//grpctransport.ClientBefore(opentracing.FromGRPCRequest(tracer, "ReadContextMetadata", logger)),
		).Endpoint()
		//readcontextmetadataEndpoint = opentracing.TraceClient(tracer, "ReadContextMetadata")(readcontextmetadataEndpoint)
		//readcontextmetadataEndpoint = limiter(readcontextmetadataEndpoint)
		//readcontextmetadataEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		//Name:    "ReadContextMetadata",
		//Timeout: 30 * time.Second,
		//}))(readcontextmetadataEndpoint)
	}

	return svc.Endpoints{

		ReadContextTestValueEndpoint: readcontexttestvalueEndpoint,
		ReadContextMetadataEndpoint:  readcontextmetadataEndpoint,
	}
}
