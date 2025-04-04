package entity

import (
	"time"

	"github.com/google/uuid"
)

const answerTable = "answer"

// Answer model
type Answer struct {
	UUID       uuid.UUID `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID         uint      `json:"-"`
	QuestionID uint      `json:"-"`
	Text       string    `json:"text"`
	IsCorrect  bool      `json:"isCorrect"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

// Answer gives table name of model
func (a Answer) TableName() string {
	return answerTable
}
