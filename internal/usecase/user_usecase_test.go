package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mnaufalhilmym/bookshelf/internal/config"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
	"github.com/mnaufalhilmym/bookshelf/internal/repository"
	"github.com/mnaufalhilmym/bookshelf/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func newUserUsecase() *usecase.UserUsecase {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	repo := repository.NewUserRepository(db)
	return usecase.NewUserUsecase(db, repo, "jwtKey", 10*time.Second)
}

func newFailUserUsecase() *usecase.UserUsecase {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	repo := &repository.UserRepository{}
	return usecase.NewUserUsecase(db, repo, "jwtKey", 10*time.Second)
}

func TestUserUsecase_Register(t *testing.T) {
	uc := newUserUsecase()
	failUc := newFailUserUsecase()

	type params struct {
		ctx     context.Context
		request *model.RegisterUserRequest
	}
	type returned struct {
		resp *model.UserResponse
		err  error
	}

	t.Run("Positive Case - register", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.RegisterUserRequest{
				Username: "unique_username",
				Password: "randompassword",
			},
		}
		returned := returned{
			resp: &model.UserResponse{
				ID:       1,
				Username: params.request.Username,
			},
			err: nil,
		}

		resp, err := uc.Register(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp.ID, resp.ID)
		assert.EqualValues(t, returned.resp.Username, resp.Username)
	})

	t.Run("Negative Case 1 - duplicate username", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.RegisterUserRequest{
				Username: "unique_username",
				Password: "randompassword",
			},
		}
		returned := returned{
			resp: nil,
			err:  model.BadRequest(errors.New("duplicate username")),
		}

		uc.Register(params.ctx, params.request)
		resp, err := uc.Register(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp, resp)
	})

	t.Run("Negative Case 2 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.RegisterUserRequest{
				Username: "unique_username",
				Password: "randompassword",
			},
		}
		returned := returned{
			resp: nil,
			err:  model.InternalServerError(errors.New("failed to create new user")),
		}

		resp, err := failUc.Register(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp, resp)
	})
}

func TestUserUsecase_Login(t *testing.T) {
	uc := newUserUsecase()
	failUc := newFailUserUsecase()

	type params struct {
		ctx      context.Context
		request1 *model.RegisterUserRequest
		request2 *model.LoginRequest
	}
	type returned struct {
		resp *model.LoginResponse
		err  error
	}

	t.Run("Positive Case - login", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request1: &model.RegisterUserRequest{
				Username: "unique_username2",
				Password: "randompassword",
			},
			request2: &model.LoginRequest{
				Username: "unique_username2",
				Password: "randompassword",
			},
		}
		returned := returned{
			resp: &model.LoginResponse{
				ID:       1,
				Username: params.request2.Username,
			},
			err: nil,
		}

		resp1, err := uc.Register(params.ctx, params.request1)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp.ID, resp1.ID)
		assert.EqualValues(t, returned.resp.Username, resp1.Username)

		resp2, err := uc.Login(params.ctx, params.request2)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp.ID, resp2.ID)
		assert.EqualValues(t, returned.resp.Username, resp2.Username)
		assert.NotEmpty(t, resp2.Token)
	})

	t.Run("Negative Case 1 - username not found", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request2: &model.LoginRequest{
				Username: "unique_username3",
				Password: "randompassword",
			},
		}
		returned := returned{
			resp: nil,
			err:  model.NotFound(errors.New("username not found")),
		}

		resp, err := uc.Login(params.ctx, params.request2)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp, resp)
	})

	t.Run("Negative Case 2 - wrong password", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request2: &model.LoginRequest{
				Username: "unique_username2",
				Password: "randompassword1",
			},
		}
		returned := returned{
			resp: nil,
			err:  model.BadRequest(errors.New("wrong password")),
		}

		resp, err := uc.Login(params.ctx, params.request2)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp, resp)
	})

	t.Run("Negative Case 3 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request2: &model.LoginRequest{
				Username: "unique_username2",
				Password: "randompassword1",
			},
		}
		returned := returned{
			resp: nil,
			err:  model.InternalServerError(errors.New("failed to find user data by username")),
		}

		resp, err := failUc.Login(params.ctx, params.request2)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.resp, resp)
	})
}
