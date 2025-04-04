package entity

import (
	"time"

	"github.com/google/uuid"
)

const questionSolutionTable = "question_solution"

// QuestionSolution model
type QuestionSolution struct {
	UUID                  uuid.UUID             `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID                    uint                  `json:"-"`
	QuizSolutionID        uint                  `json:"-"`
	QuestionUUID          uuid.UUID             `json:"question_uuid"`
	SelectedAnswers       []SelectedAnswer      `json:"selected_answers";gorm:"foreignKey:QuestionSolutionID"`
	QuestionSolutionScore QuestionSolutionScore `json:"question_solution_score";gorm:"foreignKey:QuestionSolutionID"`
	CreatedAt             time.Time             `json:"-"`
	UpdatedAt             time.Time             `json:"-"`
}

// TableName gives table name of model
func (qs QuestionSolution) TableName() string {
	return questionSolutionTable
}
