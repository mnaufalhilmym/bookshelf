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

func newAuthorHandler() *handler.AuthorHandler {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	repo := repository.NewAuthorRepository(db)
	uc := usecase.NewAuthorUsecase(db, repo)
	return handler.NewAuthorHandler(uc)
}

func createAuthor(t *testing.T, router *gin.Engine, payload *model.CreateAuthorRequest) {
	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	httpReq, err := http.NewRequest(http.MethodPost, "/v1/authors", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	testRec := httptest.NewRecorder()
	router.ServeHTTP(testRec, httpReq)
}

func TestAuthorHandler_GetMany(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	handler := newAuthorHandler()

	router.GET("/v1/authors", handler.GetMany)
	router.POST("/v1/authors", handler.Create)

	req1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}
	req2 := &model.CreateAuthorRequest{
		Name:      "Author Name 2",
		Birthdate: time.Date(2022, 2, 22, 2, 22, 22, 222, time.UTC),
	}

	createAuthor(t, router, req1)
	createAuthor(t, router, req2)

	t.Run("Positive Case 1 - get many by name", func(t *testing.T) {
		expectedRes := []model.AuthorResponse{{
			ID:        1,
			Name:      req1.Name,
			Birthdate: req1.Birthdate,
		}}

		httpReq, err := http.NewRequest(http.MethodGet, "/v1/authors?name=name%201", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[[]model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Positive Case 2 - get many by birthdate", func(t *testing.T) {
		expectedRes := []model.AuthorResponse{{
			ID:        2,
			Name:      req2.Name,
			Birthdate: req2.Birthdate,
		}}

		httpReq, err := http.NewRequest(http.MethodGet, "/v1/authors?birthdate_start=2021-10-16T00:00:00Z", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[[]model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/authors?birthdate_start=2021-10-16", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[[]model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}

func TestAuthorHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	handler := newAuthorHandler()

	router.GET("/v1/authors/:id", handler.Get)
	router.POST("/v1/authors", handler.Create)

	req1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}

	createAuthor(t, router, req1)

	t.Run("Positive Case - get by id", func(t *testing.T) {
		expectedRes := model.AuthorResponse{
			ID:        1,
			Name:      req1.Name,
			Birthdate: req1.Birthdate,
		}

		httpReq, err := http.NewRequest(http.MethodGet, "/v1/authors/1", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - validation error", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/authors/0", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field ID", res.Error)
	})

	t.Run("Negative Case 2 - wrong id", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/authors/2", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusNotFound, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "author not found", res.Error)
	})

	t.Run("Negative Case 3 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodGet, "/v1/authors/xx", nil)
		assert.NoError(t, err)

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}

func TestAuthorHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	handler := newAuthorHandler()

	router.POST("/v1/authors", handler.Create)

	t.Run("Positive Case - create author", func(t *testing.T) {
		payload := model.CreateAuthorRequest{
			Name:      "Author Name 1",
			Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		expectedRes := model.AuthorResponse{
			ID:        1,
			Name:      payload.Name,
			Birthdate: payload.Birthdate,
		}

		httpReq, err := http.NewRequest(http.MethodPost, "/v1/authors", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusCreated, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - validation error", func(t *testing.T) {
		payload := model.CreateAuthorRequest{
			Name: "Author Name 1",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/v1/authors", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field Birthdate", res.Error)
	})

	t.Run("Negative Case 2 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodPost, "/v1/authors", bytes.NewReader([]byte{}))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}

func TestAuthorHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	handler := newAuthorHandler()

	router.POST("/v1/authors", handler.Create)
	router.PUT("/v1/authors/:id", handler.Update)

	req1 := &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
	}

	createAuthor(t, router, req1)

	t.Run("Positive Case - update author", func(t *testing.T) {
		payload := model.UpdateAuthorRequest{
			ID:        1,
			Name:      util.ToPointer("Author Name Changed"),
			Birthdate: util.ToPointer(time.Date(2022, 2, 22, 2, 22, 22, 222, time.UTC)),
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		expectedRes := model.AuthorResponse{
			ID:        1,
			Name:      *payload.Name,
			Birthdate: *payload.Birthdate,
		}

		httpReq, err := http.NewRequest(http.MethodPut, "/v1/authors/1", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		payload := model.UpdateAuthorRequest{
			ID:        2,
			Name:      util.ToPointer("Author Name Changed"),
			Birthdate: util.ToPointer(time.Date(2022, 2, 22, 2, 22, 22, 222, time.UTC)),
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPut, "/v1/authors/2", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusNotFound, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "id not found", res.Error)
	})

	t.Run("Negative Case 2 - validation error", func(t *testing.T) {
		payload := model.UpdateAuthorRequest{
			ID:        1,
			Name:      util.ToPointer("Author Name Changed"),
			Birthdate: util.ToPointer(time.Date(2022, 2, 22, 2, 22, 22, 222, time.UTC)),
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPut, "/v1/authors/0", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field ID", res.Error)
	})

	t.Run("Negative Case 3 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodPut, "/v1/authors/1", bytes.NewReader([]byte{}))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}

func TestAuthorHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	handler := newAuthorHandler()

	router.POST("/v1/authors", handler.Create)
	router.DELETE("/v1/authors/:id", handler.Delete)

	t.Run("Positive Case - delete author", func(t *testing.T) {
		req1 := &model.CreateAuthorRequest{
			Name:      "Author Name 1",
			Birthdate: time.Date(2011, 11, 11, 11, 11, 11, 111, time.UTC),
		}

		createAuthor(t, router, req1)

		expectedRes := model.AuthorResponse{
			ID: 1,
		}

		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/authors/1", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes, res.Data)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/authors/2", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusNotFound, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "id not found", res.Error)
	})

	t.Run("Negative Case 2 - validation error", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/authors/0", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field ID", res.Error)
	})

	t.Run("Negative Case 3 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodDelete, "/v1/authors/xx", nil)
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.AuthorResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})
}
