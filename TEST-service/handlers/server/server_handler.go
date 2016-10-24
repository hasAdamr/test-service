package handler

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	_ "errors"
	_ "time"

	"golang.org/x/net/context"

	_ "github.com/go-kit/kit/log"
	_ "github.com/go-kit/kit/metrics"

	pb "github.com/hasAdamr/test-service/TEST-service"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() Service {
	return TESTService{}
}

type TESTService struct{}

// ReadContextTestValue implements Service.
func (s TESTService) ReadContextTestValue(ctx context.Context, in *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	_ = ctx
	_ = in
	response := pb.EmptyMessage{}
	return &response, nil
}

type Service interface {
	ReadContextTestValue(ctx context.Context, in *pb.EmptyMessage) (*pb.EmptyMessage, error)
}
