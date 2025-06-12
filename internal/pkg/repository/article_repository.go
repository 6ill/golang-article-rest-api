package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/6ill/go-article-rest-api/internal/pkg/model"
)

type ArticleRepo interface {
	CreateArticle(ctx context.Context, article model.CreateArticleRequest) (*model.Article, error)
	GetArticles(ctx context.Context, filters model.ArticleFilter) ([]model.Article, error)
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

func (r *ArticleRepoImpl) GetArticles(ctx context.Context, filters model.ArticleFilter) (articles []model.Article, err error) {
	query := `
        SELECT a.id, a.title, a.body, a.created_at, au.id, au.name
        FROM articles a
        JOIN authors au ON a.author_id = au.id
        WHERE
            ($1 = '' OR au.name ILIKE '%' || $1 || '%')
            AND
            ($2 = '' OR to_tsvector('simple', a.title || ' ' || a.body) @@ plainto_tsquery('simple', $2))
        ORDER BY a.created_at DESC
		LIMIT $3 OFFSET $4	
    `

	offset := (filters.Page - 1) * filters.PageSize

	args := []any{filters.AuthorName, filters.Query, filters.PageSize, offset}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Body, &article.CreatedAt, &article.Author.ID, &article.Author.Name); err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return articles, nil
}
