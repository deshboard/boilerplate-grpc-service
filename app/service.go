package app

import (
	context "golang.org/x/net/context"

	"github.com/deshboard/boilerplate-grpc-service/protobuf"
)

// Service implements the Protocol Buffer RPC server
type Service struct{}

// NewService creates a new service object
func NewService() *Service {
	return &Service{}
}

// Method is supposed to do something
func (s *Service) Method(ctx context.Context, r *protobuf.Request) (*protobuf.Response, error) {
	return &protobuf.Response{}, nil
}

// StreamingMethod is supposed to do something else
func (s *Service) StreamingMethod(r *protobuf.Request, stream protobuf.Boilerplate_StreamingMethodServer) error {
	for {
		if err := stream.Send(&protobuf.Response{}); err != nil {
			return err
		}
	}
}
