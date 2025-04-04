package controller

import (
	"quiz-builder/controller/auth"
	"quiz-builder/controller/quiz"
	"quiz-builder/controller/quiz_solution"
	"quiz-builder/controller/user"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(user.NewUserController),
	fx.Provide(quiz.NewQuizController),
	fx.Provide(quiz_solution.NewQuizSolutionController),
	fx.Provide(auth.NewAuthController),
)
