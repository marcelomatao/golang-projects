package database

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewDatabase),
)
