package logger

import "go.uber.org/fx"

// Module exports Logger
var Module = fx.Options(
	fx.Provide(NewLogger),
)
