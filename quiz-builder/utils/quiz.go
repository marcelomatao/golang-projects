package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"quiz-builder/entity"
	"quiz-builder/service/quiz"

	"github.com/google/uuid"

	"quiz-builder/logger"
)

func GetQuizFromParams(c *gin.Context, logger logger.Logger, quizService quiz.QuizService) (entity.Quiz, int, error) {
	paramID := c.Param("quiz_uuid")
	quizUUID, err := uuid.Parse(paramID)
	if err != nil {
		return entity.Quiz{}, http.StatusBadRequest, err
	}

	quiz, err := quizService.GetOneQuiz(quizUUID)
	if err != nil {
		return entity.Quiz{}, http.StatusInternalServerError, err
	}

	if quiz.ID <= 0 {
		err := errors.New("there is no quiz for the given quiz uuid")
		return entity.Quiz{}, http.StatusNotFound, err
	}

	return quiz, http.StatusOK, nil
}
