package config

import (
	"go.uber.org/fx"
)

// Module exports Config
var Module = fx.Options(
	fx.Provide(NewConfig),
)
