package entity

import (
	"time"

	"github.com/google/uuid"
)

const quizTable = "quiz"

// Quiz model
type Quiz struct {
	UUID          uuid.UUID      `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID            uint           `json:"-"`
	Title         string         `json:"title"`
	UserID        uint           `json:"-"`
	Published     bool           `json:"published"`
	Questions     []Question     `json:"questions";gorm:"foreignKey:QuizID;references:ID"`
	QuizSolutions []QuizSolution `json:"-";gorm:"foreignKey:QuizID;references:ID"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
}

// TableName gives table name of model
func (u Quiz) TableName() string {
	return quizTable
}
