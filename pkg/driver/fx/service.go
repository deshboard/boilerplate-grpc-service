package fx

import (
	"github.com/deshboard/boilerplate-grpc-service/pkg/driver/grpc"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"go.uber.org/dig"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService returns a new service instance.
func NewService(params ServiceParams) *grpc.Service {
	return grpc.NewService(
		grpc.Logger(params.Logger),
		grpc.ErrorHandler(params.ErrorHandler),
	)
}
