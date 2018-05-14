package app

import (
	"os"
	"testing"

	"github.com/goph/fxt/testing"
	"github.com/goph/fxt/testing/is"
)

var runnerFactoryRegistry fxtesting.RunnerFactoryRegistry

func TestMain(m *testing.M) {
	runner, err := runnerFactoryRegistry.CreateRunner()
	if err != nil {
		panic(err)
	}

	if is.Unit || is.Integration || !is.Acceptance {
		runner = fxtesting.AppendRunner(runner, m)
	}

	result := runner.Run()

	os.Exit(result)
}
