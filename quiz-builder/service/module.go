package service

import (
	"go.uber.org/fx"

	"quiz-builder/service/answer"
	"quiz-builder/service/auth"
	"quiz-builder/service/question"
	"quiz-builder/service/quiz"
	"quiz-builder/service/quiz_solution"
	"quiz-builder/service/user"
)

// Module exports services present
var Module = fx.Options(
	fx.Provide(user.NewUserService),
	fx.Provide(quiz.NewQuizService),
	fx.Provide(quiz_solution.NewQuizSolutionService),
	fx.Provide(question.NewQuestionService),
	fx.Provide(answer.NewAnswerService),
	fx.Provide(auth.NewAuthService),
)
