package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

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
const authorID2 string = "1f66ba86-2798-42f6-a994-fa869b4702f9"

var testArticles []model.Article = []model.Article{
	{
		Title: "Sea Food", Body: "fish crab", Author: model.Author{ID: authorID, Name: "Test Author1"},
	},
	{
		Title: "Turtle", Body: "plastic sea", Author: model.Author{ID: authorID, Name: "Test Author1"},
	},
	{
		Title: "Micro Plastic", Body: "kitchen utensils", Author: model.Author{ID: authorID2, Name: "Test Author2"},
	},
}

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
					Name: "Test Author1",
				},
			},
		}

		resp, err := makeHttpRequest("POST", "/api/v1/article", reqBody)
		assert.Nil(t, err)

		assert.Equalf(t, fiber.StatusCreated, resp.StatusCode, "should return 201 code")

		var responseBody helper.ResponseCreate
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Nil(t, err)

		assert.Equalf(t, expectedResponse.Message, responseBody.Message, "should return expected response message")
		if !compareArticles(t, expectedResponse.Data, responseBody.Data) {
			t.Errorf("should return expected %+v but got %+v", expectedResponse.Data, responseBody.Data)
		}
	})
}

func TestGetArticles(t *testing.T) {
	setupTestData()
	err := setupArticles()
	assert.NoError(t, err)

	t.Run("should return all articles", func(t *testing.T) {
		expectedResponse := testArticles
		resp, err := makeHttpRequest("GET", "/api/v1/article"+"?", nil)
		assert.NoError(t, err)

		assert.Equalf(t, fiber.StatusOK, resp.StatusCode, "should return 200 code")

		var responseBody helper.ResponseGetAll
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)

		length := len(expectedResponse)

		for i := range responseBody.Data {
			// Remember GET article  returned in descending order
			if !compareArticles(t, expectedResponse[length-i-1], responseBody.Data[i]) {
				t.Errorf("should return expected articles but got %+v", responseBody.Data)
			}
		}
	})
	t.Run("should 2 articles from writer Test Author1", func(t *testing.T) {
		expectedResponse := testArticles[:2]

		query := "author=author1"
		resp, err := makeHttpRequest("GET", "/api/v1/article"+"?"+query, nil)
		assert.NoError(t, err)

		assert.Equalf(t, fiber.StatusOK, resp.StatusCode, "should return 200 code")

		var responseBody helper.ResponseGetAll
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)

		length := len(expectedResponse)

		for i := range responseBody.Data {
			// Remember GET article  returned in descending order
			if !compareArticles(t, expectedResponse[length-i-1], responseBody.Data[i]) {
				t.Errorf("should return expected articles but got %+v", responseBody.Data)
			}
		}
	})

	t.Run("should return articles contain 'plastic'", func(t *testing.T) {
		expectedResponse := testArticles[1:]

		query := "query=plastic"
		resp, err := makeHttpRequest("GET", "/api/v1/article"+"?"+query, nil)
		assert.NoError(t, err)

		assert.Equalf(t, fiber.StatusOK, resp.StatusCode, "should return 200 code")

		var responseBody helper.ResponseGetAll
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)

		length := len(expectedResponse)

		for i := range responseBody.Data {
			// Remember GET article  returned in descending order
			if !compareArticles(t, expectedResponse[length-i-1], responseBody.Data[i]) {
				t.Errorf("should return expected articles but got %+v", responseBody.Data)
			}
		}
	})

	t.Run("should return nothing", func(t *testing.T) {
		query := "query=crab&author=author2"
		resp, err := makeHttpRequest("GET", "/api/v1/article"+"?"+query, nil)
		assert.NoError(t, err)

		assert.Equalf(t, fiber.StatusOK, resp.StatusCode, "should return 200 code")

		var responseBody helper.ResponseGetAll
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)

		assert.Equalf(t, 0, len(responseBody.Data), "Should return empty list but got %+v", responseBody.Data)
	})
}

func setupArticles() error {
	// Ensure the articles are inserted in order.
	for _, article := range testArticles {
		args := []any{article.Title, article.Body, article.Author.ID}
		_, err := container.Db.Exec("INSERT INTO articles (title, body, author_id) VALUES ($1, $2, $3)", args...)
		if err != nil {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
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
	args := []any{authorID, "Test Author1", authorID2, "Test Author2"}
	_, err := container.Db.Exec("INSERT INTO authors (id, name) VALUES ($1, $2), ($3, $4)", args...)
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

func compareArticles(t *testing.T, expected model.Article, actual model.Article) bool {
	t.Helper()
	if expected.Title != actual.Title {
		t.Errorf("expected title: %s but got %s", expected.Title, actual.Title)
		return false
	}

	if expected.Body != actual.Body {
		t.Errorf("expected body: %s but got %s", expected.Body, actual.Body)
		return false
	}

	if !reflect.DeepEqual(expected.Author, actual.Author) {
		t.Errorf("expected author: %+v, but got %+v", expected.Author, expected.Body)
		return false
	}

	return true
}
