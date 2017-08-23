package app

import (
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"github.com/goph/stdlib/errors"
)

// Service implements the RPC server.
type Service struct {
	Logger       log.Logger
	ErrorHandler errors.Handler
}

// NewService creates a new service object.
func NewService() *Service {
	return &Service{
		Logger:       log.NewNopLogger(),
		ErrorHandler: emperror.NewNullHandler(),
	}
}
