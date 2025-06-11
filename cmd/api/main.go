package main

import (
	"fmt"

	"github.com/6ill/go-article-rest-api/internal/helper"
	"github.com/6ill/go-article-rest-api/internal/infrastructure"
	api "github.com/6ill/go-article-rest-api/internal/server/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	container := infrastructure.InitContainer()

	app := fiber.New()
	app.Use(logger.New())

	api.HttpRouteInit(app, container)
	port := fmt.Sprintf("%s:%d", container.App.ServerHost, container.App.ServerPort)
	if err := app.Listen(port); err != nil {
		helper.Logger(helper.LoggerLevelFatal, "error", err)
	}
}
