package user

import (
	"errors"
	"net/http"
	"net/mail"

	"gorm.io/gorm"

	"quiz-builder/constants"
	"quiz-builder/entity"
	"quiz-builder/logger"
	"quiz-builder/service/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"quiz-builder/utils"
)

// UserController data type
type UserController struct {
	service user.UserService
	logger  logger.Logger
}

// NewUserController creates new user controller
func NewUserController(
	userService user.UserService,
	logger logger.Logger,
) UserController {
	return UserController{
		service: userService,
		logger:  logger,
	}
}

// GetOneUser gets one user
func (u UserController) GetOneUser(c *gin.Context) {
	paramID := c.Param("id")

	id, err := uuid.Parse(paramID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	user, err := u.service.GetOneUser(id)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})

}

// SaveUser saves the user
func (u UserController) SaveUser(c *gin.Context) {

	user := entity.User{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userByEmail, err := u.service.GetOneUserByEmail(user.Email)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check errors
	if userByEmail.Email != "" {
		err = errors.New("user with given email already registered")
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	encoded, err := utils.HashPassword(user.Password)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Password = encoded

	user.UUID = uuid.New()
	if err := u.service.WithTrx(trxHandle).CreateUser(user); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "user created"})
}
