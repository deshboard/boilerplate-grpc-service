package app

import "github.com/Sirupsen/logrus"

// Service implements the RPC server
type Service struct {
	logger logrus.FieldLogger
}

// NewService creates a new service object
func NewService(logger logrus.FieldLogger) *Service {
	return &Service{
		logger: logger,
	}
}
