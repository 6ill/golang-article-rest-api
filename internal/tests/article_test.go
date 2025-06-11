package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/6ill/go-article-rest-api/internal/helper"
	"github.com/6ill/go-article-rest-api/internal/infrastructure"
	"github.com/6ill/go-article-rest-api/internal/pkg/model"
	api "github.com/6ill/go-article-rest-api/internal/server/http"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var (
	app       *fiber.App
	container *infrastructure.Container
)

const authorID string = "1f66ba86-2798-42f6-a994-fa869b4702f8"

func TestMain(m *testing.M) {
	v := infrastructure.InitMockViper()
	container = infrastructure.InitMockContainer(v)

	app = fiber.New()

	api.HttpRouteInit(app, container)
	setupTestData()

	code := m.Run()

	cleanupTestData()
	os.Exit(code)
}

func TestCreateArticle(t *testing.T) {
	t.Run("should return validation field error", func(t *testing.T) {
		reqBody := model.CreateArticleRequest{
			Title: "Test title",
			Body:  "Test Body",
		}

		resp, err := makeHttpRequest("POST", "/api/v1/article", reqBody)
		assert.Nil(t, err)

		assert.Equalf(t, fiber.StatusBadRequest, resp.StatusCode, "should return 400 code")

		var resBodyError helper.ValidationResponse
		err = json.NewDecoder(resp.Body).Decode(&resBodyError)

		assert.Nil(t, err)
		assert.Equalf(t, "AuthorID", resBodyError.Errors[0].Field, "should have error field on AuthorID")
	})

	t.Run("should return no author found error", func(t *testing.T) {
		fakeID := "48d8e18a-b2b5-444d-9496-2ebafabcf175"
		reqBody := model.CreateArticleRequest{
			Title:    "Test title",
			Body:     "Test Body",
			AuthorID: fakeID,
		}
		expectedError := fmt.Sprintf("author with ID %s not found", fakeID)

		resp, err := makeHttpRequest("POST", "/api/v1/article", reqBody)
		assert.Nil(t, err)

		assert.Equalf(t, fiber.StatusBadRequest, resp.StatusCode, "should return 400 code")

		var resBodyError struct {
			Errors []string `json:"errors"`
		}

		err = json.NewDecoder(resp.Body).Decode(&resBodyError)
		fmt.Printf("\nerror body: %+v \n", resBodyError)

		assert.Nil(t, err)
		assert.Equalf(t, expectedError, resBodyError.Errors[0], "should return message error author not found")
	})

	t.Run("should return success response", func(t *testing.T) {
		reqBody := model.CreateArticleRequest{
			Title:    "Test title",
			Body:     "Test Body",
			AuthorID: authorID,
		}

		expectedResponse := helper.ResponseCreate{
			Message: "Your article has been created!",
			Data: model.Article{
				ID:    authorID,
				Title: "Test title",
				Body:  "Test Body",
				Author: model.Author{
					ID:   authorID,
					Name: "Test Author",
				},
			},
		}

		resp, err := makeHttpRequest("POST", "/api/v1/article", reqBody)
		assert.Nil(t, err)

		assert.Equalf(t, fiber.StatusCreated, resp.StatusCode, "should return 201 code")

		var responseBody helper.ResponseCreate
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		fmt.Printf("\nresp body: %+v\n", responseBody)
		assert.Nil(t, err)

		assert.Equalf(t, expectedResponse.Message, responseBody.Message, "should return expected response message")
		assert.Equalf(t, expectedResponse.Data.Body, responseBody.Data.Body, "should return expected article body but got %+v", responseBody)
	})
}

func cleanupTestData() {
	container.Db.Exec("DELETE FROM articles")
	container.Db.Exec("DELETE FROM authors")
	container.Db.Close()
}

func setupTestData() {
	// Clean tables
	container.Db.Exec("DELETE FROM articles")
	container.Db.Exec("DELETE FROM authors")

	// Seed author
	args := []any{authorID, "Test Author"}
	_, err := container.Db.Exec("INSERT INTO authors (id, name) VALUES ($1, $2)", args...)
	if err != nil {
		log.Fatalf("Failed to seed author: %+v", err)
	}
}

func makeHttpRequest(method string, endpoint string, reqBody any) (*http.Response, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, endpoint, bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	return resp, err
}
