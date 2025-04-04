package quiz_solution

import (
	"errors"
	"math"
	"net/http"
	"quiz-builder/constants"
	"quiz-builder/entity"
	"quiz-builder/logger"
	"quiz-builder/mapper"
	"quiz-builder/utils"

	"quiz-builder/service/answer"
	"quiz-builder/service/quiz"
	"quiz-builder/service/quiz_solution"
	"quiz-builder/service/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// QuizSolutionController data type
type QuizSolutionController struct {
	userService         user.UserService
	quizService         quiz.QuizService
	quizSolutionService quiz_solution.QuizSolutionService
	answerService       answer.AnswerService
	logger              logger.Logger
}

// NewQuizSolutionController creates new quiz solution controller
func NewQuizSolutionController(
	userService user.UserService,
	quizService quiz.QuizService,
	quizSolutionService quiz_solution.QuizSolutionService,
	answerService answer.AnswerService,
	logger logger.Logger,
) QuizSolutionController {
	return QuizSolutionController{
		userService:         userService,
		quizService:         quizService,
		quizSolutionService: quizSolutionService,
		answerService:       answerService,
		logger:              logger,
	}
}

// GetAllPublishedQuizzes gets the published quizzes so the user can solve
func (u QuizSolutionController) GetAllPublishedQuizzes(c *gin.Context) {
	quizzes, err := u.quizService.GetAllQuizzesByPublished(true)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	quizzesToSolve := mapper.GetAllPublishedQuizzesToSolve(quizzes)
	c.JSON(200, gin.H{"data": quizzesToSolve})
}

// SubmitQuizSolution publishes quiz solution
func (u QuizSolutionController) SubmitQuizSolution(c *gin.Context) {
	quiz, httpStatus, err := utils.GetQuizFromParams(c, u.logger, u.quizService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, httpStatus, err := u.validateUserHasSubmitionToQuiz(c, quiz)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	quizSolution := entity.QuizSolution{}
	if err := c.ShouldBindJSON(&quizSolution); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	httpStatus, err = u.validateAndSaveQuizSolution(c, quiz, user, quizSolution)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "quiz solution submited"})
}

// GetSolutionsToOublishedQuiz gets the quizzes solutions to the published quiz
func (u QuizSolutionController) GetSolutionsToPublishedQuiz(c *gin.Context) {
	quiz, httpStatus, err := utils.GetQuizFromParams(c, u.logger, u.quizService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, httpStatus, err = utils.ValidateUserAuthorizedAndQuizCreator(c, quiz, u.userService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	quizSolutions, err := u.quizSolutionService.GetAllQuizSolutionsByQuizID(quiz.ID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{"data": quizSolutions})
}

// GetUserQuizzesSolutions gets the quizzes solutions for the authenticated user
func (u QuizSolutionController) GetUserQuizzesSolutions(c *gin.Context) {
	user, httpStatus, err := utils.GetAuthenticatedUser(c, u.userService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err,
		})
		return
	}

	quizSolutions, err := u.quizSolutionService.GetAllQuizzesSolutionsByUserID(user.ID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{"data": quizSolutions})
}

func (u QuizSolutionController) validateAndSaveQuizSolution(c *gin.Context, quiz entity.Quiz, user entity.User, quizSolution entity.QuizSolution) (int, error) {

	quizSolutionData := &entity.QuizSolution{
		UUID:   uuid.New(),
		UserID: user.ID,
	}

	quizSolutionScore := &entity.QuizSolutionScore{
		UUID: uuid.New(),
	}

	if len(quiz.Questions) < len(quizSolution.QuestionsSolution) {
		err := errors.New("the number of questions solutions are higher than the number of questions in the quiz")
		return http.StatusBadRequest, err
	}

	mapQuestionSolutionByQuestionUUID := make(map[string]entity.QuestionSolution)
	for _, questionSolution := range quizSolution.QuestionsSolution {

		if _, ok := mapQuestionSolutionByQuestionUUID[questionSolution.QuestionUUID.String()]; ok {
			err := errors.New("there is more than one question solution for one question in the quiz")
			return http.StatusBadRequest, err
		}

		mapQuestionSolutionByQuestionUUID[questionSolution.QuestionUUID.String()] = questionSolution
	}

	mapQuestionByQuestionUUID := make(map[string]entity.Question)
	for _, question := range quiz.Questions {
		mapQuestionByQuestionUUID[question.UUID.String()] = question
	}

	questionsScoreSum := 0
	for k, v := range mapQuestionSolutionByQuestionUUID {
		question, ok := mapQuestionByQuestionUUID[k]
		if !ok {
			err := errors.New("there is one question solution where the question for QuestionUUID does not exist")
			return http.StatusBadRequest, err
		}

		if len(v.SelectedAnswers) > 1 && question.AnswerType == entity.Single {
			err := errors.New("multiple answers select to a question that expects a single answer")
			return http.StatusBadRequest, err
		}

		questionSolutionData := &entity.QuestionSolution{
			UUID:         uuid.New(),
			QuestionUUID: v.QuestionUUID,
		}

		totalCorrect, totalWrong := u.computerTotalPossibleAnswers(question)
		questionSolutionScore := &entity.QuestionSolutionScore{
			UUID:                 uuid.New(),
			TotalCorrectPossible: totalCorrect,
			TotalWrongPossible:   totalWrong,
		}

		totalCorrectSelected := 0
		totalWrongSelected := 0
		for _, selectedAnswer := range v.SelectedAnswers {
			answer, err := u.answerService.GetOneAnswer(selectedAnswer.AnswerUUID)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			if answer.ID <= 0 && answer.QuestionID == question.ID {
				err := errors.New("every selected answer must exists as option for one question of this quiz")
				return http.StatusBadRequest, err
			}

			if answer.IsCorrect {
				totalCorrectSelected++
			} else {
				totalWrongSelected++
			}

			selectedAnswerData := &entity.SelectedAnswer{
				UUID:       uuid.New(),
				AnswerUUID: selectedAnswer.AnswerUUID,
			}

			questionSolutionData.SelectedAnswers = append(questionSolutionData.SelectedAnswers, *selectedAnswerData)
		}

		questionSolutionScore.TotalCorrectSelected = totalCorrectSelected
		questionSolutionScore.TotalWrongSelected = totalWrongSelected
		if question.AnswerType == entity.Single {
			questionSolutionScore.Score = (totalCorrectSelected - totalWrongSelected) * 100
		} else {
			scoreCorrect := math.Round((float64(totalCorrectSelected) / float64(totalCorrect)) * 100.0)
			if totalWrong == 0 {
				totalWrong++
			}
			scoreWrong := math.Round((float64(totalWrongSelected) / float64(totalWrong)) * 100.0)
			questionSolutionScore.Score = int(scoreCorrect - scoreWrong)
		}

		questionsScoreSum = questionsScoreSum + questionSolutionScore.Score
		questionSolutionData.QuestionSolutionScore = *questionSolutionScore
		quizSolutionData.QuestionsSolution = append(quizSolutionData.QuestionsSolution, *questionSolutionData)
	}

	quizSolutionScore.Score = int(math.Round(float64(questionsScoreSum) / float64((len(quiz.Questions) * 100.0)) * 100))
	quizSolutionData.QuizSolutionScore = *quizSolutionScore

	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	quiz.QuizSolutions = append(quiz.QuizSolutions, *quizSolutionData)
	if err := u.quizService.WithTrx(trxHandle).UpdateQuiz(quiz); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (u QuizSolutionController) computerTotalPossibleAnswers(question entity.Question) (int, int) {
	totalCorrect := 0
	totalWrong := 0
	for _, answer := range question.Answers {
		if answer.IsCorrect {
			totalCorrect++
		} else {
			totalWrong++
		}
	}

	return totalCorrect, totalWrong
}

func (u QuizSolutionController) validateUserHasSubmitionToQuiz(c *gin.Context, quiz entity.Quiz) (entity.User, int, error) {
	user, httpStatus, err := utils.GetAuthenticatedUser(c, u.userService)
	if err != nil {
		return entity.User{}, httpStatus, err
	}

	quizSolution, err := u.quizSolutionService.GetOneQuizSolutionByQuizIDAndUserID(quiz.ID, user.ID)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}

	if quizSolution.UserID == user.ID && quizSolution.QuizID == quiz.ID {
		err = errors.New("user already submitted a solution to this quiz")
		return entity.User{}, http.StatusBadRequest, err
	}

	return user, http.StatusOK, nil
}
