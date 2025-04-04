package utils

import (
	"errors"
	"net/http"
	"quiz-builder/constants"
	"quiz-builder/entity"

	"quiz-builder/service/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthenticatedUser(c *gin.Context, u user.UserService) (entity.User, int, error) {
	uuidClaim := c.Request.Header.Get(constants.UUIDHeader)
	uuidObj, err := uuid.Parse(uuidClaim)
	if err != nil {
		return entity.User{}, http.StatusBadRequest, err
	}

	user, err := u.GetOneUser(uuidObj)
	if err != nil {
		return entity.User{}, http.StatusInternalServerError, err
	}

	if user.ID <= 0 {
		err = errors.New("user does not exist")
		return entity.User{}, http.StatusUnauthorized, err
	}

	return user, http.StatusOK, nil
}

func ValidateUserAuthorizedAndQuizCreator(c *gin.Context, quiz entity.Quiz, userService user.UserService) (entity.User, int, error) {
	user, httpStatus, err := GetAuthenticatedUser(c, userService)
	if err != nil {
		return entity.User{}, httpStatus, err
	}

	if quiz.UserID != user.ID {
		err = errors.New("user not authorized to handle this quiz")
		return entity.User{}, http.StatusUnauthorized, err
	}

	return user, http.StatusOK, nil
}
