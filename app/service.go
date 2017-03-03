package app

import (
	deshboard "github.com/deshboard/boilerplate-grpc-service/model"
	context "golang.org/x/net/context"
)

// Service implements the Protocol Buffer RPC server
type Service struct{}

// NewService creates a new service object
func NewService() *Service {
	return &Service{}
}

// Method is supposed to do something
func (s *Service) Method(ctx context.Context, r *deshboard.BoilerplateRequest) (*deshboard.BoilerplateResponse, error) {
	return &deshboard.BoilerplateResponse{}, nil
}

// StreamingMethod is supposed to do something else
func (s *Service) StreamingMethod(r *deshboard.BoilerplateRequest, stream deshboard.Boilerplate_StreamingMethodServer) error {
	for {
		if err := stream.Send(&deshboard.BoilerplateResponse{}); err != nil {
			return err
		}
	}
}
