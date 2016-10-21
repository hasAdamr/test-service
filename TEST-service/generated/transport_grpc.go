package svc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	//stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	//"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/hasAdamr/test-service/TEST-service"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(ctx context.Context, endpoints Endpoints /*, tracer stdopentracing.Tracer, logger log.Logger*/) pb.TestServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerBefore(PrintContextValues),
	}
	return &grpcServer{
		// testservice

		readcontexttestvalue: grpctransport.NewServer(
			ctx,
			endpoints.ReadContextTestValueEndpoint, DecodeGRPCReadContextTestValueRequest,
			EncodeGRPCReadContextTestValueResponse,
			options...,
		),
	}
}

type grpcServer struct {
	readcontexttestvalue grpctransport.Handler
	readcontextmetadata  grpctransport.Handler
}

// Methods

func (s *grpcServer) ReadContextTestValue(ctx context.Context, req *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	_, rep, err := s.readcontexttestvalue.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.EmptyMessage), nil
}

func (s *grpcServer) ReadContextMetadata(ctx context.Context, req *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	_, rep, err := s.readcontextmetadata.ServeGRPC(ctx, req)
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

// DecodeGRPCReadContextMetadataRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC readcontextmetadata request to a user-domain readcontextmetadata request. Primarily useful in a server.
func DecodeGRPCReadContextMetadataRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
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

// DecodeGRPCReadContextMetadataResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC readcontextmetadata reply to a user-domain readcontextmetadata response. Primarily useful in a client.
func DecodeGRPCReadContextMetadataResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
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

// EncodeGRPCReadContextMetadataResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain readcontextmetadata response to a gRPC readcontextmetadata reply. Primarily useful in a server.
func EncodeGRPCReadContextMetadataResponse(_ context.Context, response interface{}) (interface{}, error) {
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

// EncodeGRPCReadContextMetadataRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain readcontextmetadata request to a gRPC readcontextmetadata request. Primarily useful in a client.
func EncodeGRPCReadContextMetadataRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.EmptyMessage)
	return req, nil
}

func PrintContextValues(ctx context.Context, inMD *metadata.MD) context.Context {
	md := *inMD
	for k, v := range md {
		ctx = context.WithValue(ctx, k, v[0])
	}

	ctx = metadata.NewContext(ctx, *inMD)

	return ctx
}
