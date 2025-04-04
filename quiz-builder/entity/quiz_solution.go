package entity

import (
	"time"

	"github.com/google/uuid"
)

const quizSolutionTable = "quiz_solution"

// QuizSolution model
type QuizSolution struct {
	UUID              uuid.UUID          `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID                uint               `json:"-"`
	UserID            uint               `json:"-"`
	QuizID            uint               `json:"-"`
	QuestionsSolution []QuestionSolution `json:"questions_solution";gorm:"foreignKey:QuizSolutionID;references:ID"`
	QuizSolutionScore QuizSolutionScore  `json:"quiz_solution_score";gorm:"foreignKey:QuestionSolutionID"`
	CreatedAt         time.Time          `json:"-"`
	UpdatedAt         time.Time          `json:"-"`
}

// TableName gives table name of model
func (u QuizSolution) TableName() string {
	return quizSolutionTable
}
