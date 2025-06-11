package handler

import (
	"github.com/6ill/go-article-rest-api/internal/pkg/controller"
	"github.com/6ill/go-article-rest-api/internal/pkg/service"
	"github.com/gofiber/fiber/v2"
)

func ArticleHandler(r fiber.Router, articleService service.ArticleService) {
	controller := controller.NewArticleController(articleService)

	api := r.Group("/article")
	api.Post("/", controller.CreateArticle)
	api.Get("/", controller.GetArticles)
}
