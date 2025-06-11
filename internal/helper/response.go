package helper

import "github.com/6ill/go-article-rest-api/internal/pkg/model"

type ResponseCreate struct {
	Message string        `json:"message"`
	Data    model.Article `json:"data"`
}

type ResponseGetAll struct {
	Data []model.Article `json:"data"`
}
