package quiz

import (
	"quiz-builder/entity"
	"quiz-builder/logger"

	"quiz-builder/repository/quiz_builder"

	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IQuizService interface {
	WithTrx(trxHandle *gorm.DB) IQuizService
	GetOneQuiz(uuid uuid.UUID) (entity.Quiz, error)
	GetAllQuiz() ([]entity.Quiz, error)
	GetOneQuizById(quizID uint) (quiz entity.Quiz, err error)
	GetAllQuizzesByUserId(userID uint) (quizzes []entity.Quiz, err error)
	GetAllQuizzesByPublished(publish bool) (quizzes []entity.Quiz, err error)
	CreateQuiz(entity.Quiz) error
	UpdateQuiz(entity.Quiz) error
	DeleteQuiz(uuid uuid.UUID) error
}

// QuizService service layer
type QuizService struct {
	logger     logger.Logger
	repository quiz_builder.QuizBuilderRepository
}

// NewQuizService creates a new quiz service
func NewQuizService(
	logger logger.Logger,
	repository quiz_builder.QuizBuilderRepository,
) QuizService {
	return QuizService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s QuizService) WithTrx(trxHandle *gorm.DB) QuizService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// GetOneQuiz gets one quiz
func (s QuizService) GetOneQuiz(uuid uuid.UUID) (quiz entity.Quiz, err error) {
	return quiz, s.repository.Preload("Questions").Preload("Questions.Answers").Preload(clause.Associations).Find(&quiz, "UUID = ?", uuid.String()).Error
}

// GetOneQuizById gets one quiz by id
func (s QuizService) GetOneQuizById(quizID uint) (quiz entity.Quiz, err error) {
	return quiz, s.repository.Preload("Questions").Preload("Questions.Answers").Preload(clause.Associations).Find(&quiz, quizID).Error
}

// GetAllQuizzesByUserId gets one quiz by id
func (s QuizService) GetAllQuizzesByUserId(userID uint) (quizzes []entity.Quiz, err error) {
	return quizzes, s.repository.Where("user_id = ?", userID).Preload("Questions").Preload("Questions.Answers").Preload(clause.Associations).Find(&quizzes).Error
}

func (s QuizService) GetAllQuiz() (quizzes []entity.Quiz, err error) {
	return quizzes, s.repository.Find(&quizzes).Error
}

func (s QuizService) GetAllQuizzesByPublished(publish bool) (quizzes []entity.Quiz, err error) {
	return quizzes, s.repository.Preload("Questions").Preload("Questions.Answers").Preload(clause.Associations).Find(&quizzes, "Published = ?", publish).Error
}

// CreateQuiz call to create the quiz
func (s QuizService) CreateQuiz(quiz entity.Quiz) error {
	return s.repository.Create(&quiz).Error
}

// UpdateQuiz updates the quiz
func (s QuizService) UpdateQuiz(quiz entity.Quiz) error {
	return s.repository.Save(&quiz).Error
}

// DeleteQuiz deletes the quiz
func (s QuizService) DeleteQuiz(uuid uuid.UUID) error {
	return s.repository.Delete(&entity.Quiz{}, "UUID = ?", uuid.String()).Error
}
