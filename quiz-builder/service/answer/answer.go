package answer

import (
	"quiz-builder/entity"

	"quiz-builder/repository/quiz_builder"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type IAnswerService interface {
	WithTrx(trxHandle *gorm.DB) IAnswerService
	GetOneAnswer(uuid uuid.UUID) (entity.Answer, error)
	GetAllAnswers() ([]entity.Answer, error)
	CreateAnswer(entity.Answer) error
	UpdateAnswer(entity.Answer) error
	DeleteAnswer(uuid uuid.UUID) error
}

// AnswerService service layer
type AnswerService struct {
	repository quiz_builder.QuizBuilderRepository
}

// NewAnswerService creates a new answer service
func NewAnswerService(
	repository quiz_builder.QuizBuilderRepository,
) AnswerService {
	return AnswerService{
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s AnswerService) WithTrx(trxHandle *gorm.DB) AnswerService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// GetOneAnswer gets one answer
func (s AnswerService) GetOneAnswer(uuid uuid.UUID) (answer entity.Answer, err error) {
	return answer, s.repository.Find(&answer, "UUID = ?", uuid.String()).Error
}

// GetAllAnswer get all the answer
func (s AnswerService) GetAllAnswers() (answers []entity.Answer, err error) {
	return answers, s.repository.Find(&answers).Error
}

// CreateAnswer call to create the answer
func (s AnswerService) CreateAnswer(answer entity.Answer) error {
	return s.repository.Create(&answer).Error
}

// UpdateAnswer updates the answer
func (s AnswerService) UpdateAnswer(answer entity.Answer) error {
	return s.repository.Save(&answer).Error
}

// DeleteAnswer deletes the answer
func (s AnswerService) DeleteAnswer(uuid uuid.UUID) error {
	return s.repository.Delete(&entity.Answer{}, "UUID = ?", uuid.String()).Error
}
