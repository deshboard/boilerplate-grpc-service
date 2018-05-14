package fx

import (
	"go.uber.org/fx"
)

// Module is an fx compatible module.
var Module = fx.Options(
	fx.Provide(NewService),
)
