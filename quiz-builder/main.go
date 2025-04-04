package main

import (
	"context"

	"fmt"

	"quiz-builder/app"
	"quiz-builder/config"
	"quiz-builder/filter"
	"quiz-builder/gin"
	"quiz-builder/handler"
	"quiz-builder/logger"

	"go.uber.org/fx"

	"github.com/joho/godotenv"
)

func opts() fx.Option {
	return fx.Options(
		app.Module,
		fx.Invoke(Run),
	)
}

func main() {
	_ = godotenv.Load()

	app := fx.New(opts())
	ctx := context.Background()
	err := app.Start(ctx)
	defer app.Stop(ctx)
	if err != nil {
		fmt.Errorf(err.Error())
	}

}

func Run(route handler.Routes, filter filter.Filters, config config.Config, logger logger.Logger, ginHandler gin.RequestHandler) {
	filter.Setup()
	route.Setup()

	logger.Info("Running server")
	if config.ServerPort == "" {
		_ = ginHandler.Gin.Run()
	} else {
		_ = ginHandler.Gin.Run(":" + config.ServerPort)
	}
	ginHandler.Gin.Run(":8080")
}
