package entity

import (
	"time"

	"github.com/google/uuid"
)

const selectedAnswerTable = "selected_answer"

// SelectedAnswer model
type SelectedAnswer struct {
	UUID               uuid.UUID `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID                 uint      `json:"-"`
	QuestionSolutionID uint      `json:""-"`
	AnswerUUID         uuid.UUID `json:"answer_uuid"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

// TableName gives table name of model
func (qs SelectedAnswer) TableName() string {
	return selectedAnswerTable
}
