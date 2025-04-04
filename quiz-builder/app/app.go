package app

import (
	"quiz-builder/config"
	"quiz-builder/controller"
	"quiz-builder/database"
	"quiz-builder/filter"
	"quiz-builder/gin"
	"quiz-builder/handler"
	"quiz-builder/logger"
	"quiz-builder/repository"
	"quiz-builder/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	logger.Module,
	gin.Module,
	database.Module,
	handler.Module,
	filter.Module,
	service.Module,
	repository.Module,
	controller.Module,
)
