package entity

import (
	"time"

	"github.com/google/uuid"
)

const questionSolutionScoreTable = "question_solution_score"

// QuestionSolutionScore model
type QuestionSolutionScore struct {
	UUID                 uuid.UUID `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID                   uint      `json:"-"`
	QuestionSolutionID   uint      `json:"-"`
	TotalCorrectPossible int       `json:"-"`
	TotalWrongPossible   int       `json:"-"`
	TotalCorrectSelected int       `json:"-"`
	TotalWrongSelected   int       `json:"-"`
	Score                int       `json:"score"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"-"`
}

// TableName gives table name of model
func (qs QuestionSolutionScore) TableName() string {
	return questionSolutionScoreTable
}
