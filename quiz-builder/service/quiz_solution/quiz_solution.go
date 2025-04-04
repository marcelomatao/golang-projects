package quiz_solution

import (
	"quiz-builder/entity"
	"quiz-builder/logger"

	"quiz-builder/repository/quiz_builder"

	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IQuizSolutionService interface {
	WithTrx(trxHandle *gorm.DB) IQuizSolutionService
	GetOneQuizSolution(uuid uuid.UUID) (entity.QuizSolution, error)
	GetAllQuizSolution() ([]entity.QuizSolution, error)
	GetOneQuizSolutionByQuizIDAndUserID(quizID uint, userID uint) (quizSolution entity.QuizSolution, err error)
	GetAllQuizSolutionsByQuizID(quizID uint) (quizSolutions []entity.QuizSolution, err error)
	GetAllQuizzesSolutionsByUserID(userID uint) ([]entity.QuizSolution, error)
	CreateQuizSolution(entity.QuizSolution) error
	UpdateQuizSolution(entity.QuizSolution) error
	DeleteQuizSolution(uuid uuid.UUID) error
}

// QuizSolutionService service layer
type QuizSolutionService struct {
	logger     logger.Logger
	repository quiz_builder.QuizBuilderRepository
}

// NewQuizSolutionService creates a new quiz solution service
func NewQuizSolutionService(
	logger logger.Logger,
	repository quiz_builder.QuizBuilderRepository,
) QuizSolutionService {
	return QuizSolutionService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s QuizSolutionService) WithTrx(trxHandle *gorm.DB) QuizSolutionService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// GetOneQuizSolution gets one quiz
func (s QuizSolutionService) GetOneQuizSolution(uuid uuid.UUID) (quizSolution entity.QuizSolution, err error) {
	return quizSolution, s.repository.Find(&quizSolution, "UUID = ?", uuid.String()).Error
}

// GetOneQuizSolutionByQuizIDAndUserID gets one quiz by quiz id and user id
func (s QuizSolutionService) GetOneQuizSolutionByQuizIDAndUserID(quizID uint, userID uint) (quizSolution entity.QuizSolution, err error) {
	return quizSolution, s.repository.Where("user_id = ? AND quiz_id = ?", userID, quizID).Find(&quizSolution).Error
}

// GetAllQuizzesSolutionsByUserID gets all quizzes solutions submitions by user
func (s QuizSolutionService) GetAllQuizzesSolutionsByUserID(userID uint) (quizzesSolutions []entity.QuizSolution, err error) {
	return quizzesSolutions, s.repository.Where("user_id = ?", userID).Preload("QuizSolutionScore").Preload("QuestionsSolution").Preload("QuestionsSolution.SelectedAnswers").Preload("QuestionsSolution.QuestionSolutionScore").Preload(clause.Associations).Find(&quizzesSolutions).Error
}

// GetAllQuizSolutions get all the quiz solutions
func (s QuizSolutionService) GetAllQuizSolutions() (quizSolutions []entity.QuizSolution, err error) {
	return quizSolutions, s.repository.Find(&quizSolutions).Error
}

// GetAllQuizSolutionsByQuizID get all the quiz solutions by quiz ID
func (s QuizSolutionService) GetAllQuizSolutionsByQuizID(quizID uint) (quizSolutions []entity.QuizSolution, err error) {
	return quizSolutions, s.repository.Where("quiz_id = ?", quizID).Preload("QuizSolutionScore").Preload("QuestionsSolution").Preload("QuestionsSolution.SelectedAnswers").Preload("QuestionsSolution.QuestionSolutionScore").Preload(clause.Associations).Find(&quizSolutions).Error
}

// CreateQuizSolution call to create the quiz solution
func (s QuizSolutionService) CreateQuizSolution(quizSolution entity.QuizSolution) error {
	return s.repository.Create(&quizSolution).Error
}

// UpdateQuizSolution updates the quiz solution
func (s QuizSolutionService) UpdateQuizSolution(quizSolution entity.QuizSolution) error {
	return s.repository.Save(&quizSolution).Error
}

// DeleteQuiz deletes the quiz solution
func (s QuizSolutionService) DeleteQuizSolution(uuid uuid.UUID) error {
	return s.repository.Delete(&entity.QuizSolution{}, "UUID = ?", uuid.String()).Error
}
