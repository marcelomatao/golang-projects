package quiz

import (
	"errors"
	"fmt"
	"net/http"
	"quiz-builder/constants"
	"quiz-builder/entity"
	"quiz-builder/logger"
	"quiz-builder/utils"

	"quiz-builder/service/answer"
	"quiz-builder/service/question"
	"quiz-builder/service/quiz"
	"quiz-builder/service/user"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// QuizController data type
type QuizController struct {
	quizService     quiz.QuizService
	questionService question.QuestionService
	answerService   answer.AnswerService
	userService     user.UserService
	logger          logger.Logger
}

// NewQuizController creates new quiz controller
func NewQuizController(
	quizService quiz.QuizService,
	questionService question.QuestionService,
	userService user.UserService,
	answerService answer.AnswerService,
	logger logger.Logger,
) QuizController {
	return QuizController{
		quizService:     quizService,
		questionService: questionService,
		answerService:   answerService,
		userService:     userService,
		logger:          logger,
	}
}

// GetOneQuiz gets one quiz
func (u QuizController) GetOneOwnedQuiz(c *gin.Context) {
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

	c.JSON(200, gin.H{
		"data": quiz,
	})

}

// GetAllOwnedQuiz gets all quizzes created by authenticated user
func (u QuizController) GetAllOwnedQuiz(c *gin.Context) {
	user, httpStatus, err := utils.GetAuthenticatedUser(c, u.userService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	quizzes, err := u.quizService.GetAllQuizzesByUserId(user.ID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": quizzes,
	})
}

// SaveQuiz saves the quiz
func (u QuizController) SaveQuiz(c *gin.Context) {
	quiz := entity.Quiz{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&quiz); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, httpStatus, err := utils.GetAuthenticatedUser(c, u.userService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	quiz.UUID = uuid.New()

	user.QuizzesCreated = append(user.QuizzesCreated, quiz)
	if err := u.userService.WithTrx(trxHandle).UpdateUser(user); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "quiz created"})
}

// PublishQuiz publishes quiz
func (u QuizController) PublishQuiz(c *gin.Context) {
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

	err = u.validateQuestionsAndAnswers(quiz.Questions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	quiz.Published = true

	if err := u.quizService.WithTrx(trxHandle).UpdateQuiz(quiz); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "quiz published"})
}

// AddQuestion add question
func (u QuizController) AddQuestion(c *gin.Context) {
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

	if quiz.Published {
		err := errors.New("quiz can't be edited after published")
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	question := entity.Question{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&question); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	question.UUID = uuid.New()

	quiz.Questions = append(quiz.Questions, question)

	err = u.validateNumberOfQuestions(quiz.Questions)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := u.quizService.WithTrx(trxHandle).UpdateQuiz(quiz); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "question added to the quiz"})
}

// AddAnswer add answer
func (u QuizController) AddAnswer(c *gin.Context) {
	paramQuestionID := c.Param("question_uuid")

	questionid, err := uuid.Parse(paramQuestionID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	question, err := u.questionService.GetOneQuestion(questionid)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	quiz, err := u.quizService.GetOneQuizById(question.QuizID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if quiz.Published {
		err := errors.New("quiz can't be edited after published")
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, httpStatus, err := utils.ValidateUserAuthorizedAndQuizCreator(c, quiz, u.userService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	answer := entity.Answer{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&answer); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	answer.UUID = uuid.New()

	question.Answers = append(question.Answers, answer)

	err = u.validateQuestionAnswers(question.Answers, question.AnswerType)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := u.questionService.WithTrx(trxHandle).UpdateQuestion(question); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "answer added to the question"})
}

// DeleteQuiz deletes quiz
func (u QuizController) DeleteQuiz(c *gin.Context) {
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

	for _, question := range quiz.Questions {
		for _, answer := range question.Answers {
			if err := u.answerService.DeleteAnswer(answer.UUID); err != nil {
				u.logger.Error(err)
			}
		}

		if err := u.questionService.DeleteQuestion(question.UUID); err != nil {
			u.logger.Error(err)
		}
	}

	if err := u.quizService.DeleteQuiz(quiz.UUID); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "quiz deleted"})
}

// DeleteQuestion delete question
func (u QuizController) DeleteQuestion(c *gin.Context) {
	paramID := c.Param("question_uuid")

	id, err := uuid.Parse(paramID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	question, err := u.questionService.GetOneQuestion(id)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	quiz, err := u.quizService.GetOneQuizById(question.QuizID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, httpStatus, err := utils.ValidateUserAuthorizedAndQuizCreator(c, quiz, u.userService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	if quiz.Published {
		err := errors.New("quiz can't be edited after published")
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, answer := range question.Answers {
		if err := u.answerService.DeleteAnswer(answer.UUID); err != nil {
			u.logger.Error(err)
		}
	}

	if err := u.questionService.DeleteQuestion(id); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "question deleted"})
}

// DeleteAnswer delete answer
func (u QuizController) DeleteAnswer(c *gin.Context) {
	paramID := c.Param("answer_uuid")

	id, err := uuid.Parse(paramID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	answer, err := u.answerService.GetOneAnswer(id)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	question, err := u.questionService.GetOneQuestionById(answer.QuestionID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	quiz, err := u.quizService.GetOneQuizById(question.QuizID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, httpStatus, err := utils.ValidateUserAuthorizedAndQuizCreator(c, quiz, u.userService)
	if err != nil {
		u.logger.Error(err)
		c.JSON(httpStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	if quiz.Published {
		err := errors.New("quiz can't be edited after published")
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := u.answerService.DeleteAnswer(id); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "answer deleted"})
}

func (u QuizController) validateQuestionsAndAnswers(questons []entity.Question) error {
	if len(questons) < 1 {
		err := errors.New("quiz must have at least one question")
		u.logger.Error(err)
		return err
	}

	for _, question := range questons {
		if len(question.Answers) < 1 {
			err := errors.New("each question must have at least one answer and it has " + fmt.Sprintf("%d", len(question.Answers)))
			u.logger.Error(err)
			return err
		}

		hasOneCorrectAnswer := false
		for _, answer := range question.Answers {
			if answer.IsCorrect {
				hasOneCorrectAnswer = true
				break
			}
		}
		if !hasOneCorrectAnswer {
			err := errors.New("each question must have at least one correct answer")
			u.logger.Error(err)
			return err
		}
	}

	return nil
}

func (u QuizController) validateNumberOfQuestions(questons []entity.Question) error {
	if len(questons) > 10 {
		err := errors.New("each quiz must have 10 questions at maximum")
		u.logger.Error(err)
		return err
	}

	return nil
}

func (u QuizController) validateQuestionAnswers(answers []entity.Answer, answerType entity.AnswerType) error {
	if len(answers) > 5 {
		err := errors.New("question must have 5 answers at maximum")
		u.logger.Error(err)
		return err
	}

	correctCount := 0
	for _, answer := range answers {
		if answer.IsCorrect {
			correctCount++
		}
	}

	if correctCount > 1 && answerType == entity.Single {
		err := errors.New("questions with single type expect only one correct answer")
		u.logger.Error(err)
		return err
	}

	return nil
}
