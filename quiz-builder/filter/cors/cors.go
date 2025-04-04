package cors

import (
	"quiz-builder/config"
	ginHandler "quiz-builder/gin"
	"quiz-builder/logger"

	cors "github.com/rs/cors/wrapper/gin"
)

// CorsFilter filter for cors
type CorsFilter struct {
	handler ginHandler.RequestHandler
	logger  logger.Logger
	config  config.Config
}

// NewCorsFilter creates new cors filter
func NewCorsFilter(handler ginHandler.RequestHandler, logger logger.Logger, config config.Config) CorsFilter {
	return CorsFilter{
		handler: handler,
		logger:  logger,
		config:  config,
	}
}

// Setup sets up cors filter
func (m CorsFilter) Setup() {
	m.logger.Info("Setting up cors filter")

	debug := m.config.Environment == "development"
	m.handler.Gin.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		Debug:            debug,
	}))
}
