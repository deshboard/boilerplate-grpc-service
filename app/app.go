package app

import (
	"github.com/goph/fxt/app/default"
	"go.uber.org/fx"
)

// Module is the collection of all modules of the application.
var Module = fx.Options(
	fxdefaultapp.Module,

	// Configuration
	fx.Provide(
		NewLoggerConfig,
		NewDebugConfig,
	),
)

type Runner = fxdefaultapp.Runner
