package quiz_builder

import (
	"quiz-builder/database"
	"quiz-builder/logger"

	"gorm.io/gorm"
)

// QuizBuilderRepository database structure
type QuizBuilderRepository struct {
	database.Database
	logger logger.Logger
}

// NewQuizBuilderRepository creates a new quiz builder repository
func NewQuizBuilderRepository(
	db database.Database,
	logger logger.Logger) QuizBuilderRepository {
	return QuizBuilderRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository with transaction
func (r QuizBuilderRepository) WithTrx(trxHandle *gorm.DB) QuizBuilderRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}
