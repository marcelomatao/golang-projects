package user

import (
	"quiz-builder/controller/user"
	ginHandler "quiz-builder/gin"
	"quiz-builder/logger"
)

// UserHandler struct
type UserHandler struct {
	logger         logger.Logger
	handler        ginHandler.RequestHandler
	userController user.UserController
}

// NewUser creates new user handler
func NewUser(
	logger logger.Logger,
	handler ginHandler.RequestHandler,
	userController user.UserController,
) UserHandler {
	return UserHandler{
		handler:        handler,
		logger:         logger,
		userController: userController,
	}
}

// Setup user routes
func (g UserHandler) Setup() {

	g.logger.Info("Setting up routes")

}
