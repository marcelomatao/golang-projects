package repository

import (
	"quiz-builder/repository/quiz_builder"

	"go.uber.org/fx"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(quiz_builder.NewQuizBuilderRepository),
)
