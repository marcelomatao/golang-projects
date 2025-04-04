package handler

import (
	"quiz-builder/handler/auth"
	"quiz-builder/handler/quiz_builder"
	"quiz-builder/handler/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(quiz_builder.NewQuizBuilder),
	fx.Provide(user.NewUser),
	fx.Provide(auth.NewAuth),
	fx.Provide(NewRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	quizBuilderHandler quiz_builder.QuizBuilderHandler,
	userHandler user.UserHandler,
	authHandler auth.AuthHandler,
) Routes {
	return Routes{
		quizBuilderHandler,
		userHandler,
		authHandler,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
