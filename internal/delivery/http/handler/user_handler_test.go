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
	"github.com/stretchr/testify/assert"
)

func newUserHandler() *handler.UserHandler {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUsecase(db, repo, "jwtKey", 10*time.Second)
	return handler.NewUserHandler(uc)
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	handler := newUserHandler()

	router.POST("/auth/register", handler.Register)

	t.Run("Positive Case - register", func(t *testing.T) {
		payload := model.RegisterUserRequest{
			Username: "unique_username",
			Password: "random-password",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		expectedRes := model.UserResponse{
			ID:       1,
			Username: payload.Username,
		}

		httpReq, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusCreated, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes.ID, res.Data.ID)
		assert.EqualValues(t, expectedRes.Username, res.Data.Username)
	})

	t.Run("Negative Case 1 - duplicate username", func(t *testing.T) {
		payload := model.RegisterUserRequest{
			Username: "unique_username",
			Password: "random-password",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		httpReq, err = http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec = httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "duplicate username", res.Error)
	})

	t.Run("Negative Case 2 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte{}))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})

	t.Run("Negative Case 3 - validation error", func(t *testing.T) {
		payload := model.RegisterUserRequest{
			Username: "unique_username",
			Password: "",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field Password", res.Error)
	})
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	handler := newUserHandler()

	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)

	registerPayload := model.RegisterUserRequest{
		Username: "unique_username",
		Password: "random-password",
	}
	reqBody, err := json.Marshal(registerPayload)
	assert.NoError(t, err)

	httpReq, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	testRec := httptest.NewRecorder()
	router.ServeHTTP(testRec, httpReq)

	t.Run("Positive Case - login", func(t *testing.T) {
		payload := model.RegisterUserRequest{
			Username: "unique_username",
			Password: "random-password",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		expectedRes := model.LoginResponse{
			ID:       1,
			Username: payload.Username,
		}

		httpReq, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusOK, testRec.Code)

		res := new(model.Response[model.LoginResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, expectedRes.ID, res.Data.ID)
		assert.EqualValues(t, expectedRes.Username, res.Data.Username)
		assert.NotEmpty(t, res.Data.Token)
	})

	t.Run("Negative Case 1 - wrong username", func(t *testing.T) {
		payload := model.RegisterUserRequest{
			Username: "wrong_username",
			Password: "random-password",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusNotFound, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "username not found", res.Error)
	})

	t.Run("Negative Case 2 - wrong password", func(t *testing.T) {
		payload := model.RegisterUserRequest{
			Username: "unique_username",
			Password: "wrong-password",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "wrong password", res.Error)
	})

	t.Run("Negative Case 3 - invalid request body", func(t *testing.T) {
		httpReq, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte{}))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "failed to parse request", res.Error)
	})

	t.Run("Negative Case 4 - validation error", func(t *testing.T) {
		payload := model.RegisterUserRequest{
			Username: "unique_username",
			Password: "",
		}
		reqBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		testRec := httptest.NewRecorder()
		router.ServeHTTP(testRec, httpReq)

		assert.EqualValues(t, http.StatusBadRequest, testRec.Code)

		res := new(model.Response[model.UserResponse])
		assert.NoError(t, json.Unmarshal(testRec.Body.Bytes(), res))

		assert.EqualValues(t, "validation error in field Password", res.Error)
	})
}
