package question

import (
	"quiz-builder/entity"
	"quiz-builder/logger"

	"quiz-builder/repository/quiz_builder"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type IQuestionService interface {
	WithTrx(trxHandle *gorm.DB) IQuestionService
	GetOneQuestion(uuid uuid.UUID) (entity.Question, error)
	GetAllQuestion() ([]entity.Question, error)
	CreateQuestion(entity.Question) error
	UpdateQuestion(entity.Question) error
	DeleteQuestion(uuid uuid.UUID) error
}

// QuestionService service layer
type QuestionService struct {
	logger     logger.Logger
	repository quiz_builder.QuizBuilderRepository
}

// NewQuestionService creates a new question service
func NewQuestionService(
	logger logger.Logger,
	repository quiz_builder.QuizBuilderRepository,
) QuestionService {
	return QuestionService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s QuestionService) WithTrx(trxHandle *gorm.DB) QuestionService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// GetOneQuestion gets one question
func (s QuestionService) GetOneQuestion(uuid uuid.UUID) (question entity.Question, err error) {
	return question, s.repository.Preload("Answers").Find(&question, "UUID = ?", uuid.String()).Error
}

// GetOneQuestionById gets one question by id
func (s QuestionService) GetOneQuestionById(questionID uint) (question entity.Question, err error) {
	return question, s.repository.Preload("Answers").Find(&question, questionID).Error
}

// GetAllQuestion get all the question
func (s QuestionService) GetAllQuestion() (questions []entity.Question, err error) {
	return questions, s.repository.Find(&questions).Error
}

// CreateQuestion call to create the question
func (s QuestionService) CreateQuestion(question entity.Question) error {
	return s.repository.Create(&question).Error
}

// UpdateQuestion updates the question
func (s QuestionService) UpdateQuestion(question entity.Question) error {
	return s.repository.Save(&question).Error
}

// DeleteQuestion deletes the question
func (s QuestionService) DeleteQuestion(uuid uuid.UUID) error {
	return s.repository.Delete(&entity.Question{}, "UUID = ?", uuid.String()).Error
}
