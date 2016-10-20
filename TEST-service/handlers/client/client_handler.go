package clienthandler

import (
	pb "github.com/hasAdamr/test-service/TEST-service"
)

// ReadContextTestValue implements Service.
func ReadContextTestValue() (*pb.EmptyMessage, error) {

	request := pb.EmptyMessage{}
	return &request, nil
}

// ReadContextMetadata implements Service.
func ReadContextMetadata() (*pb.EmptyMessage, error) {

	request := pb.EmptyMessage{}
	return &request, nil
}
