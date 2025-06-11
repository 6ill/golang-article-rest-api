package http

import (
	"github.com/6ill/go-article-rest-api/internal/infrastructure"
	"github.com/6ill/go-article-rest-api/internal/server/http/handler"
	"github.com/gofiber/fiber/v2"
)

func HttpRouteInit(r *fiber.App, container *infrastructure.Container) {
	api := r.Group("/api/v1")

	handler.ArticleHandler(api, container.ArticleService)
}
