package quiz_builder

import (
	"quiz-builder/controller/quiz"
	"quiz-builder/controller/quiz_solution"

	"quiz-builder/filter/auth"
	ginHandler "quiz-builder/gin"
	"quiz-builder/logger"
)

// QuizBuilderHandler struct
type QuizBuilderHandler struct {
	logger                 logger.Logger
	handler                ginHandler.RequestHandler
	quizController         quiz.QuizController
	quizSolutionController quiz_solution.QuizSolutionController
	authFilter             auth.AuthFilter
}

// NewQuizBuilder creates new quiz-builder handler
func NewQuizBuilder(
	logger logger.Logger,
	handler ginHandler.RequestHandler,
	quizController quiz.QuizController,
	quizSolutionController quiz_solution.QuizSolutionController,
	authFilter auth.AuthFilter,

) QuizBuilderHandler {
	return QuizBuilderHandler{
		handler:                handler,
		logger:                 logger,
		quizController:         quizController,
		quizSolutionController: quizSolutionController,
		authFilter:             authFilter,
	}
}

// Setup quiz-builder routes
func (g QuizBuilderHandler) Setup() {

	g.logger.Info("Setting up routes")

	api := g.handler.Gin.Group("/api/quiz").Use(g.authFilter.Handler())
	{
		api.POST("", g.quizController.SaveQuiz)                                 // OK
		api.GET("/all", g.quizController.GetAllOwnedQuiz)                       // OK
		api.POST("/:quiz_uuid/publish", g.quizController.PublishQuiz)           // OK
		api.POST("/:quiz_uuid/question", g.quizController.AddQuestion)          // OK
		api.POST("/question/:question_uuid/answer", g.quizController.AddAnswer) // OK
		api.DELETE("/:quiz_uuid", g.quizController.DeleteQuiz)                  // OK
		api.DELETE("/question/:question_uuid", g.quizController.DeleteQuestion) // OK
		api.DELETE("/answer/:answer_uuid", g.quizController.DeleteAnswer)       // OK
		api.GET("/:quiz_uuid", g.quizController.GetOneOwnedQuiz)                // OK

		api.GET("/published/all", g.quizSolutionController.GetAllPublishedQuizzes)                 // OK
		api.POST("/:quiz_uuid/solution/submit", g.quizSolutionController.SubmitQuizSolution)       // OK
		api.GET("/:quiz_uuid/solutions/all", g.quizSolutionController.GetSolutionsToPublishedQuiz) // OK
		api.GET("/solutions", g.quizSolutionController.GetUserQuizzesSolutions)                    // OK
	}

}
