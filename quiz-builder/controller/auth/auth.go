package auth

import (
	"fmt"
	"net/http"
	"quiz-builder/logger"
	"quiz-builder/service/auth"
	"quiz-builder/service/user"

	"github.com/gin-gonic/gin"

	"quiz-builder/entity"
	"quiz-builder/utils"
)

// AuthController struct
type AuthController struct {
	logger      logger.Logger
	service     auth.AuthService
	userService user.UserService
}

// NewAuthController creates new controller
func NewAuthController(
	logger logger.Logger,
	service auth.AuthService,
	userService user.UserService,
) AuthController {
	return AuthController{
		logger:      logger,
		service:     service,
		userService: userService,
	}
}

// SignIn signs in user
func (auth AuthController) SignIn(c *gin.Context) {
	auth.logger.Info("SignIn route called")

	login := entity.Login{}
	if err := c.ShouldBindJSON(&login); err != nil {
		auth.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := auth.userService.GetOneUserByEmail(login.Email)
	if err != nil {
		auth.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.CheckPasswordHash(user.Password, login.Password); err != nil {
		auth.logger.Error("Password is not valid")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Password is not valid"),
		})
		return
	}

	token := auth.service.CreateToken(user)
	c.JSON(200, gin.H{
		"message": "logged in successfully",
		"token":   token,
	})
}
