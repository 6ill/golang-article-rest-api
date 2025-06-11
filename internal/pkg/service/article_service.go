package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/6ill/go-article-rest-api/internal/helper"
	"github.com/6ill/go-article-rest-api/internal/pkg/model"
	"github.com/6ill/go-article-rest-api/internal/pkg/repository"
	"github.com/gofiber/fiber/v2"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, params model.CreateArticleRequest) (*model.Article, *helper.ErrorStruct)
	// GetArticles(ctx context.Context, authorName, searchQuery string, limit, offset int) ([]model.Article, *helper.ErrorStruct)
}

type ArticleServiceImpl struct {
	repo repository.ArticleRepo
}

func NewArticleService(articleRepo repository.ArticleRepo) ArticleService {
	return &ArticleServiceImpl{
		repo: articleRepo,
	}
}

func (s *ArticleServiceImpl) CreateArticle(ctx context.Context, params model.CreateArticleRequest) (*model.Article, *helper.ErrorStruct) {
	article, err := s.repo.CreateArticle(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, &helper.ErrorStruct{
				Err:  fmt.Errorf("author with ID %s not found", params.AuthorID),
				Code: fiber.StatusBadRequest,
			}
		default:
			return nil, &helper.ErrorStruct{
				Err:  err,
				Code: fiber.StatusInternalServerError,
			}
		}
	}

	return article, nil
}

// func (s *ArticleServiceImpl) GetArticles(ctx context.Context, authorName string, searchQuery string, limit int, offset int) ([]model.Article, *helper.ErrorStruct) {
// 	articles, err := s.repo.
// }
