package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mnaufalhilmym/bookshelf/internal/config"
	"github.com/mnaufalhilmym/bookshelf/internal/delivery/http/handler"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
	"github.com/mnaufalhilmym/bookshelf/internal/repository"
	"github.com/mnaufalhilmym/bookshelf/internal/usecase"
	"github.com/mnaufalhilmym/bookshelf/internal/util"
	"github.com/stretchr/testify/assert"
)

func newAuthorAndBookHandler() (*handler.AuthorHandler, *handler.BookHandler) {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	authorRepo := repository.NewAuthorRepository(db)
	bookRepo := repository.NewBookRepository(db)
	authorUc := usecase.NewAuthorUsecase(db, authorRepo)
	bookUc := usecase.NewBookUsecase(db, bookRepo, authorRepo)
	authorHandler := handler.NewAuthorHandler(authorUc)
	bookHandler := handler.NewBookHandler(bookUc)
	return authorHandler, bookHandler
}

func createBook(t *testing.T, router *gin.Engine, payload *model.CreateBookRequest) {
	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	httpReq, err := http.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	testRec := httptest.NewRecorder()
	router.ServeHTTP(testRec, httpReq)
}

func TestBookHandler_GetMany(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	authorHandler, userHandler := newAuthorAndBookHandler()

	router.POST("/v1/authors", authorHandler.Create)
	router.GET("/v1/books", userHandler.GetMany)
	router.POST("/v1/books", userHandler.Create)

	reqAuthor1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}

	createAuthor(t, router, reqAuthor1)

	reqBook1 := &model.CreateBookRequest{
		Title:    "Book Title 1",
		ISBN:     "978-1451673319",
		AuthorID: 1,
	}
	reqBook2 := &model.CreateBookRequest{
		Title:    "Book Title 2",
		ISBN:     "978-1503290563",
		AuthorID: 1,
	}

	createBook(t, router, reqBook1)
	createBook(t, router, reqBook2)

	t.Run("Positive Case 1 - get many by title", func(t *testing.T) {
		expectedRes := []model.BookResponse{{
			ID:         1,
			Title:      reqBook1.Title,
			ISBN:       reqBook1.ISBN,
			AuthorID:   reqBook1.AuthorID,
			AuthorName: reqAuthor1.Name,
		}}

		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books?title=1", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[[]model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Positive Case 2 - get many by isbn", func(t *testing.T) {
		expectedRes := []model.BookResponse{{
			ID:         2,
			Title:      reqBook2.Title,
			ISBN:       reqBook2.ISBN,
			AuthorID:   reqBook2.AuthorID,
			AuthorName: reqAuthor1.Name,
		}}

		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books?isbn=0563", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[[]model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - invalid request", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books?author_id=xx", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[[]model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})

	t.Run("Negative Case 2 - validation error", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books?author_id=0", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[[]model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field AuthorID", res.Error)
	})
}

func TestBookHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	authorHandler, userHandler := newAuthorAndBookHandler()

	router.POST("/v1/authors", authorHandler.Create)
	router.GET("/v1/books/:id", userHandler.Get)
	router.POST("/v1/books", userHandler.Create)

	reqAuthor1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}
	createAuthor(t, router, reqAuthor1)

	reqBook1 := &model.CreateBookRequest{
		Title:    "Book Title 1",
		ISBN:     "978-1451673319",
		AuthorID: 1,
	}
	createBook(t, router, reqBook1)

	t.Run("Positive Case - get by id", func(t *testing.T) {
		expectedRes := model.BookResponse{
			ID:         1,
			Title:      reqBook1.Title,
			ISBN:       reqBook1.ISBN,
			AuthorID:   reqBook1.AuthorID,
			AuthorName: reqAuthor1.Name,
		}

		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books/1", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - validation error", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books/0", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field ID", res.Error)
	})

	t.Run("Negative Case 2 - wrong id", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books/2", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusNotFound, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "book not found", res.Error)
	})

	t.Run("Negative Case 3 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/books/xx", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}

func TestBookHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	authorHandler, userHandler := newAuthorAndBookHandler()

	router.POST("/v1/authors", authorHandler.Create)
	router.POST("/v1/books", userHandler.Create)

	reqAuthor1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}
	createAuthor(t, router, reqAuthor1)

	t.Run("Positive Case - create book", func(t *testing.T) {
		payload := model.CreateBookRequest{
			Title:    "Book Title 1",
			ISBN:     "978-1503290563",
			AuthorID: 1,
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		expectedRes := model.BookResponse{
			ID:         1,
			Title:      payload.Title,
			ISBN:       payload.ISBN,
			AuthorID:   payload.AuthorID,
			AuthorName: reqAuthor1.Name,
		}

		httpReq, err := http.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusCreated, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - validation error", func(t *testing.T) {
		payload := model.CreateBookRequest{
			Title: "Book Title 1",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field ISBN", res.Error)
	})

	t.Run("Negative Case 2 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader([]byte{}))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}

func TestBookHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	authorHandler, userHandler := newAuthorAndBookHandler()

	router.POST("/v1/authors", authorHandler.Create)
	router.POST("/v1/books", userHandler.Create)
	router.PUT("/v1/books/:id", userHandler.Update)

	reqAuthor1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}
	createAuthor(t, router, reqAuthor1)

	reqBook1 := &model.CreateBookRequest{
		Title:    "Book Title 1",
		ISBN:     "978-1451673319",
		AuthorID: 1,
	}
	createBook(t, router, reqBook1)

	t.Run("Positive Case - update book", func(t *testing.T) {
		payload := model.UpdateBookRequest{
			ID:    1,
			Title: util.ToPointer("Book Title Changed"),
			ISBN:  util.ToPointer("978-0062315007"),
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		expectedRes := model.BookResponse{
			ID:         1,
			Title:      *payload.Title,
			ISBN:       *payload.ISBN,
			AuthorID:   reqBook1.AuthorID,
			AuthorName: reqAuthor1.Name,
		}

		httpReq, err := http.NewRequest(http.MethodPut, "/v1/books/1", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		payload := model.UpdateBookRequest{
			ID: 2,
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPut, "/v1/books/2", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusNotFound, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "id not found", res.Error)
	})

	t.Run("Negative Case 2 - validation error", func(t *testing.T) {
		payload := model.UpdateBookRequest{
			ID: 1,
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPut, "/v1/books/0", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field ID", res.Error)
	})

	t.Run("Negative Case 3 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodPut, "/v1/books/1", bytes.NewReader([]byte{}))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}

func TestBookHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	authorHandler, userHandler := newAuthorAndBookHandler()

	router.POST("/v1/authors", authorHandler.Create)
	router.POST("/v1/books", userHandler.Create)
	router.DELETE("/v1/books/:id", userHandler.Delete)

	reqAuthor1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}
	createAuthor(t, router, reqAuthor1)

	t.Run("Positive Case - delete book", func(t *testing.T) {
		reqBook1 := &model.CreateBookRequest{
			Title:    "Book Title 1",
			ISBN:     "978-1451673319",
			AuthorID: 1,
		}
		createBook(t, router, reqBook1)

		expectedRes := model.BookResponse{
			ID: 1,
		}

		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/books/1", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/books/2", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusNotFound, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "id not found", res.Error)
	})

	t.Run("Negative Case 2 - validation error", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/books/0", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field ID", res.Error)
	})

	t.Run("Negative Case 3 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/books/xx", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.BookResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}
