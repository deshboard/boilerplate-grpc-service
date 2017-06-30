package app

import (
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
)

// Service implements the RPC server
type Service struct {
	logger       log.Logger
	errorHandler emperror.Handler
}

// NewService creates a new service object
func NewService(logger log.Logger, errorHandler emperror.Handler) *Service {
	return &Service{logger, errorHandler}
}
