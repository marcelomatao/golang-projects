package entity

import (
	"time"

	"github.com/google/uuid"
)

const userTable = "user"

// User model
type User struct {
	UUID             uuid.UUID      `json:"uuid";gorm:"type:uuid;default:uuid_generate_v4()"`
	ID               uint           `json:"-"`
	Name             string         `json:"name"`
	Email            string         `json:"email"`
	QuizzesCreated   []Quiz         `json:"quizzes_created";gorm:"foreignKey:UserID;references:ID"`
	QuizzesSolutions []QuizSolution `json:"quizzes_solutions";gorm:"foreignKey:UserID;references:ID"`
	Password         string         `json:"password"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
}

// TableName gives table name of model
func (u User) TableName() string {
	return userTable
}
