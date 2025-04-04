package entity

import (
	"github.com/google/uuid"
)

type QuizRequestResponse struct {
	UUID      uuid.UUID                 `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	Title     string                    `json:"title"`
	Questions []QuestionRequestResponse `json:"questions";gorm:"foreignKey:QuizID;references:ID"`
}

type QuestionRequestResponse struct {
	UUID       uuid.UUID               `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	Text       string                  `json:"text"`
	AnswerType AnswerType              `json:"answerType"`
	Answers    []AnswerRequestResponse `json:"answers";gorm:"foreignKey:QuestionID;references:ID"`
}

type AnswerRequestResponse struct {
	UUID uuid.UUID `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	Text string    `json:"text"`
}
