package entity

import (
	"time"

	"github.com/google/uuid"
)

const questionTable = "question"

type AnswerType string

const (
	Single   AnswerType = "Single"
	Multiple AnswerType = "Multiple"
)

// Question model
type Question struct {
	UUID       uuid.UUID  `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID         uint       `json:"-"`
	QuizID     uint       `json:"-"`
	Text       string     `json:"text"`
	AnswerType AnswerType `json:"answerType"`
	Answers    []Answer   `json:"answers";gorm:"foreignKey:QuestionID;references:ID"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
}

// TableName gives table name of model
func (u Question) TableName() string {
	return questionTable
}
