package database

import (
	"fmt"

	"quiz-builder/config"
	"quiz-builder/entity"
	"quiz-builder/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database modal
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(config config.Config, logger logger.Logger) Database {

	username := config.DBUsername
	password := config.DBPassword
	host := config.DBHost
	port := config.DBPort
	dbname := config.DBName

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.GetGormLogger(),
	})

	if err != nil {
		logger.Info("Url: ", url)
		logger.Panic(err)
	}

	logger.Info("Database connection established")

	db.AutoMigrate(&entity.Answer{})
	db.AutoMigrate(&entity.QuestionSolutionScore{})
	db.AutoMigrate(&entity.QuizSolutionScore{})
	db.AutoMigrate(&entity.Question{})
	db.AutoMigrate(&entity.SelectedAnswer{})
	db.AutoMigrate(&entity.QuestionSolution{})
	db.AutoMigrate(&entity.QuizSolution{})
	db.AutoMigrate(&entity.Quiz{})
	db.AutoMigrate(&entity.User{})

	logger.Info("Database connection established")

	return Database{
		DB: db,
	}
}
