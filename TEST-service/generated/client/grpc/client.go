// Package grpc provides a gRPC client for the TestService service.
package grpc

import (
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/hasAdamr/test-service/TEST-service"
	svc "github.com/hasAdamr/test-service/TEST-service/generated"
	handler "github.com/hasAdamr/test-service/TEST-service/handlers/server"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn) handler.Service {
	//options := []grpctransport.ServerOption{
	//grpctransport.ServerBefore(),
	//}
	var readcontexttestvalueEndpoint endpoint.Endpoint
	{
		readcontexttestvalueEndpoint = grpctransport.NewClient(
			conn,
			"TEST.TestService",
			"ReadContextTestValue",
			svc.EncodeGRPCReadContextTestValueRequest,
			svc.DecodeGRPCReadContextTestValueResponse,
			pb.EmptyMessage{},
			//options...,
		).Endpoint()
	}

	return svc.Endpoints{
		ReadContextTestValueEndpoint: readcontexttestvalueEndpoint,
	}
}
