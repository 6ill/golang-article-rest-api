package model

import "time"

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Author    Author    `json:"author"`
}

type CreateArticleRequest struct {
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
	AuthorID string `json:"author_id" validate:"required,uuid4"`
}

type ArticleFilter struct {
	Query      string `query:"query" validate:"omitempty"`
	AuthorName string `query:"author" validate:"omitempty"`
	Page       int    `query:"page" validate:"omitempty,gte=1"`
	PageSize   int    `query:"page_size" validate:"omitempty,gte=1,lte=100"`
}
