package controller

import (
	"github.com/6ill/go-article-rest-api/internal/helper"
	"github.com/6ill/go-article-rest-api/internal/pkg/model"
	"github.com/6ill/go-article-rest-api/internal/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type ArticleController interface {
	CreateArticle(c *fiber.Ctx) error
}

type ArticleControllerImpl struct {
	service service.ArticleService
}

func NewArticleController(articleService service.ArticleService) ArticleController {
	return &ArticleControllerImpl{
		service: articleService,
	}
}

func (co *ArticleControllerImpl) CreateArticle(c *fiber.Ctx) error {
	ctx := c.Context()

	data := new(model.CreateArticleRequest)
	if isValid, errValidation := helper.ExtractValidateRequestBody(data, c); !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(errValidation)
	}

	res, errMsg := co.service.CreateArticle(ctx, *data)
	if errMsg != nil {
		return c.Status(errMsg.Code).JSON(fiber.Map{
			"errors": []string{errMsg.Err.Error()},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(helper.ResponseCreate{
		Message: "Your article has been created!",
		Data:    *res,
	})
}
