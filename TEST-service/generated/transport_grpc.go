package svc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	"golang.org/x/net/context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/hasAdamr/test-service/TEST-service"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC TestServiceServer.
func MakeGRPCServer(ctx context.Context, endpoints Endpoints) pb.TestServiceServer {
	//options := []grpctransport.ServerOption{
	// grpctransport.ServiceBefore()
	//}
	return &grpcServer{
		// testservice

		readcontexttestvalue: grpctransport.NewServer(
			ctx,
			endpoints.ReadContextTestValueEndpoint,
			DecodeGRPCReadContextTestValueRequest,
			EncodeGRPCReadContextTestValueResponse,
			//options...,
		),
	}
}

// grpcServer implements the TestServiceServer interface
type grpcServer struct {
	readcontexttestvalue grpctransport.Handler
}

// Methods for grpcServer to implement TestServiceServer interface

func (s *grpcServer) ReadContextTestValue(ctx context.Context, req *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	_, rep, err := s.readcontexttestvalue.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.EmptyMessage), nil
}

// Server Decode

// DecodeGRPCReadContextTestValueRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC readcontexttestvalue request to a user-domain readcontexttestvalue request. Primarily useful in a server.
func DecodeGRPCReadContextTestValueRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.EmptyMessage)
	return req, nil
}

// Client Decode

// DecodeGRPCReadContextTestValueResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC readcontexttestvalue reply to a user-domain readcontexttestvalue response. Primarily useful in a client.
func DecodeGRPCReadContextTestValueResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.EmptyMessage)
	return reply, nil
}

// Server Encode

// EncodeGRPCReadContextTestValueResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain readcontexttestvalue response to a gRPC readcontexttestvalue reply. Primarily useful in a server.
func EncodeGRPCReadContextTestValueResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.EmptyMessage)
	return resp, nil
}

// Client Encode

// EncodeGRPCReadContextTestValueRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain readcontexttestvalue request to a gRPC readcontexttestvalue request. Primarily useful in a client.
func EncodeGRPCReadContextTestValueRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.EmptyMessage)
	return req, nil
}
