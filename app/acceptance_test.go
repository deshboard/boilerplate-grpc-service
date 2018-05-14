// +build acceptance

package app

import (
	"github.com/goph/fxt/dev"
	"github.com/goph/fxt/testing"
	fxacceptance "github.com/goph/fxt/testing/acceptance"
	"go.uber.org/fx"
)

func init() {
	fxdev.LoadEnvFromFile("../.env.test")
	fxdev.LoadEnvFromFile("../.env.dist")

	runnerFactoryRegistry.Register(fxtesting.RunnerFactoryFunc(AcceptanceRunnerFactory))
}

func AcceptanceRunnerFactory() (fxtesting.Runner, error) {
	acceptanceRunner := fxtesting.NewGodogRunner()

	config, err := newConfig()
	if err != nil {
		panic(err)
	}

	appContext := fxacceptance.NewAppContext(
		fx.NopLogger,
		fx.Provide(func() Config { return config }, newApplicationInfo),
		Module,
	)

	acceptanceRunner.RegisterFeatureContext(appContext.BeforeFeatureContext)
	acceptanceRunner.RegisterFeatureContext(appContext.AfterFeatureContext)

	return acceptanceRunner, nil
}
