package repository

import (
	"context"
	"database/sql"

	"github.com/6ill/go-article-rest-api/internal/pkg/model"
)

type ArticleRepo interface {
	CreateArticle(ctx context.Context, article model.CreateArticleRequest) (*model.Article, error)
}

type ArticleRepoImpl struct {
	db *sql.DB
}

func NewArticleRepo(db *sql.DB) ArticleRepo {
	return &ArticleRepoImpl{
		db: db,
	}
}

func (r *ArticleRepoImpl) CreateArticle(ctx context.Context, article model.CreateArticleRequest) (*model.Article, error) {
	var author model.Author
	err := r.db.QueryRowContext(ctx, "SELECT id, name FROM authors WHERE id = $1", article.AuthorID).Scan(&author.ID, &author.Name)

	if err != nil {
		return nil, err
	}
	// if err == sql.ErrNoRows {
	// 	return nil, fmt.Errorf("author with ID %s not found", article.AuthorID)
	// } else if err != nil {
	// 	return nil, err
	// }

	newArticle := model.Article{
		Title:  article.Title,
		Body:   article.Body,
		Author: author,
	}

	query := `INSERT INTO articles (title, body, author_id) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`

	err = r.db.QueryRowContext(ctx, query, newArticle.Title, newArticle.Body, newArticle.Author.ID).Scan(&newArticle.ID, &newArticle.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &newArticle, nil
}
