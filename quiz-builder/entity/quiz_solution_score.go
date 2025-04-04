package entity

import (
	"time"

	"github.com/google/uuid"
)

const quizSolutionScoreTable = "quiz_solution_score"

// QuizSolutionScore model
type QuizSolutionScore struct {
	UUID           uuid.UUID `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID             uint      `json:"-"`
	QuizSolutionID uint      `json:"-"`
	Score          int       `json:"score"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// TableName gives table name of model
func (qs QuizSolutionScore) TableName() string {
	return quizSolutionScoreTable
}
