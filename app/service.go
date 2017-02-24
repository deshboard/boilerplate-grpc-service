package app

import (
	context "golang.org/x/net/context"

	"github.com/deshboard/boilerplate-grpc-service/model/boilerplate"
)

// Service implements the Protocol Buffer RPC server
type Service struct{}

// NewService creates a new service object
func NewService() *Service {
	return &Service{}
}

// Method is supposed to do something
func (s *Service) Method(ctx context.Context, r *boilerplate.Request) (*boilerplate.Response, error) {
	return &boilerplate.Response{}, nil
}

// StreamingMethod is supposed to do something else
func (s *Service) StreamingMethod(r *boilerplate.Request, stream boilerplate.Boilerplate_StreamingMethodServer) error {
	for {
		if err := stream.Send(&boilerplate.Response{}); err != nil {
			return err
		}
	}
}
