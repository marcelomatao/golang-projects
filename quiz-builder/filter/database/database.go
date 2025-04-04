package database

import (
	"net/http"

	"quiz-builder/constants"
	"quiz-builder/database"
	ginHandler "quiz-builder/gin"
	"quiz-builder/logger"

	"github.com/gin-gonic/gin"
)

// DatabaseTrx filter for transactions support for database
type DatabaseTrx struct {
	handler ginHandler.RequestHandler
	logger  logger.Logger
	db      database.Database
}

// statusInList function checks if context writer status is in provided list
func statusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

// NewDatabaseTrx creates new database transactions filter
func NewDatabaseTrx(
	handler ginHandler.RequestHandler,
	logger logger.Logger,
	db database.Database,
) DatabaseTrx {
	return DatabaseTrx{
		handler: handler,
		logger:  logger,
		db:      db,
	}
}

// Setup sets up database transaction filter
func (m DatabaseTrx) Setup() {
	m.logger.Info("setting up database transaction filter")

	m.handler.Gin.Use(func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		m.logger.Info("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
			}
		}()

		c.Set(constants.DBTransaction, txHandle)
		c.Next()

		// commit transaction on success status
		if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
			m.logger.Info("committing transactions")
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Error("trx commit error: ", err)
			}
		} else {
			m.logger.Info("rolling back transaction due to status code: 500")
			txHandle.Rollback()
		}
	})
}
