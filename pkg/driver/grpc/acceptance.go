// +build acceptance

package grpc

import (
	"github.com/DATA-DOG/godog"
	"github.com/goph/fxt/testing"
)

func RegisterSuite(runner *fxtesting.GodogRunner) {
	runner.RegisterFeaturePath("../../features")
	runner.RegisterFeatureContext(FeatureContext)
}

func FeatureContext(s *godog.Suite) {
	// Add steps here
}
