package auth

import (
	"quiz-builder/controller/auth"
	"quiz-builder/controller/user"
	ginHandler "quiz-builder/gin"
	"quiz-builder/logger"
)

// AuthHandler struct
type AuthHandler struct {
	logger         logger.Logger
	handler        ginHandler.RequestHandler
	authController auth.AuthController
	userController user.UserController
}

// NewAuth creates new auth handler
func NewAuth(
	logger logger.Logger,
	handler ginHandler.RequestHandler,
	authController auth.AuthController,
	userController user.UserController,
) AuthHandler {
	return AuthHandler{
		authController: authController,
		userController: userController,
		handler:        handler,
		logger:         logger,
	}
}

// Setup auth routes
func (g AuthHandler) Setup() {

	g.logger.Info("Setting up routes")

	api := g.handler.Gin.Group("/api/auth")
	{
		api.POST("/login", g.authController.SignIn)
		api.POST("/register", g.userController.SaveUser)
	}

}
