package app

import "github.com/go-kit/kit/log"

// Service implements the RPC server
type Service struct {
	logger log.Logger
}

// NewService creates a new service object
func NewService(logger log.Logger) *Service {
	return &Service{
		logger: logger,
	}
}
