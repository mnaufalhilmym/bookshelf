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
	"github.com/mnaufalhilmym/bookshelf/internal/util"
	"github.com/stretchr/testify/assert"
)

func newAuthorUsecase() *usecase.AuthorUsecase {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	repo := repository.NewAuthorRepository(db)
	return usecase.NewAuthorUsecase(db, repo)
}

func newFailAuthorUsecase() *usecase.AuthorUsecase {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	repo := &repository.AuthorRepository{}
	return usecase.NewAuthorUsecase(db, repo)
}

func TestAuthorUsecase_GetMany(t *testing.T) {
	uc := newAuthorUsecase()
	failUc := newFailAuthorUsecase()

	type params struct {
		ctx     context.Context
		request *model.GetManyAuthorsRequest
	}
	type returned struct {
		data  []model.AuthorResponse
		total int64
		err   error
	}

	author1, err := uc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
	})
	assert.NoError(t, err)
	author2, err := uc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 2",
		Birthdate: time.Date(2022, 2, 22, 2, 22, 22, 2222, time.UTC),
	})
	assert.NoError(t, err)

	t.Run("Positive Case 1 - search by name", func(t *testing.T) {
		request := &model.GetManyAuthorsRequest{
			Name: util.ToPointer("Name"),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.AuthorResponse{*author1, *author2},
			total: 2,
			err:   nil,
		}

		resp, total, err := uc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Positive Case 2 - search by birthdate start", func(t *testing.T) {
		request := &model.GetManyAuthorsRequest{
			BirthdateStart: util.ToPointer(time.Date(2012, 1, 11, 1, 11, 11, 1111, time.UTC)),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.AuthorResponse{*author2},
			total: 1,
			err:   nil,
		}

		resp, total, err := uc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Positive Case 3 - search by birthdate end", func(t *testing.T) {
		request := &model.GetManyAuthorsRequest{
			BirthdateEnd: util.ToPointer(time.Date(2012, 1, 11, 1, 11, 11, 1111, time.UTC)),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.AuthorResponse{*author1},
			total: 1,
			err:   nil,
		}

		resp, total, err := uc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Positive Case 4 - empty search result", func(t *testing.T) {
		request := &model.GetManyAuthorsRequest{
			Name:           util.ToPointer("Author 3"),
			BirthdateStart: util.ToPointer(time.Date(2010, 1, 11, 1, 11, 11, 1111, time.UTC)),
			BirthdateEnd:   util.ToPointer(time.Date(2023, 1, 11, 1, 11, 11, 1111, time.UTC)),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.AuthorResponse{},
			total: 0,
			err:   nil,
		}

		resp, total, err := uc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Negative Case - db error", func(t *testing.T) {
		request := &model.GetManyAuthorsRequest{
			Name:           util.ToPointer("Author 3"),
			BirthdateStart: util.ToPointer(time.Date(2010, 1, 11, 1, 11, 11, 1111, time.UTC)),
			BirthdateEnd:   util.ToPointer(time.Date(2023, 1, 11, 1, 11, 11, 1111, time.UTC)),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  nil,
			total: 0,
			err:   model.InternalServerError(errors.New("failed to get many authors")),
		}

		resp, total, err := failUc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})
}

func TestAuthorUsecase_Get(t *testing.T) {
	uc := newAuthorUsecase()
	failUc := newFailAuthorUsecase()

	type params struct {
		ctx     context.Context
		request *model.GetAuthorRequest
	}
	type returned struct {
		data *model.AuthorResponse
		err  error
	}

	author, err := uc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name",
		Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
	})
	assert.NoError(t, err)

	t.Run("Positive Case - get by id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.GetAuthorRequest{
				ID: author.ID,
			},
		}
		returned := returned{
			data: author,
			err:  nil,
		}

		resp, err := uc.Get(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.GetAuthorRequest{
				ID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("author not found")),
		}

		resp, err := uc.Get(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 2 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.GetAuthorRequest{
				ID: 1,
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to find author data by id")),
		}

		resp, err := failUc.Get(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}

func TestAuthorUsecase_Create(t *testing.T) {
	uc := newAuthorUsecase()
	failUc := newFailAuthorUsecase()

	type params struct {
		ctx     context.Context
		request *model.CreateAuthorRequest
	}
	type returned struct {
		data *model.AuthorResponse
		err  error
	}

	t.Run("Positive Case - create author", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.CreateAuthorRequest{
				Name:      "Author Name 1",
				Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
			},
		}
		returned := returned{
			data: &model.AuthorResponse{
				ID:        1,
				Name:      params.request.Name,
				Birthdate: params.request.Birthdate,
			},
			err: nil,
		}

		resp, err := uc.Create(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.CreateAuthorRequest{
				Name:      "Author Name 1",
				Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to create new author")),
		}

		resp, err := failUc.Create(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}

func TestAuthorUsecase_Update(t *testing.T) {
	uc := newAuthorUsecase()
	failUc := newFailAuthorUsecase()

	type params struct {
		ctx     context.Context
		request *model.UpdateAuthorRequest
	}
	type returned struct {
		data *model.AuthorResponse
		err  error
	}

	t.Run("Positive Case 1 - update author name", func(t *testing.T) {
		author, err := uc.Create(context.Background(), &model.CreateAuthorRequest{
			Name:      "Author Name",
			Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateAuthorRequest{
				ID:   author.ID,
				Name: util.ToPointer("Author Name Changed"),
			},
		}
		returned := returned{
			data: &model.AuthorResponse{
				ID:        author.ID,
				Name:      *params.request.Name,
				Birthdate: author.Birthdate,
			},
			err: nil,
		}

		resp, err := uc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Positive Case 2 - update author birthdate", func(t *testing.T) {
		author, err := uc.Create(context.Background(), &model.CreateAuthorRequest{
			Name:      "Author Name",
			Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateAuthorRequest{
				ID:        author.ID,
				Birthdate: util.ToPointer(time.Date(2012, 1, 11, 1, 11, 11, 1111, time.UTC)),
			},
		}
		returned := returned{
			data: &model.AuthorResponse{
				ID:        author.ID,
				Name:      author.Name,
				Birthdate: *params.request.Birthdate,
			},
			err: nil,
		}

		resp, err := uc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Positive Case 3 - update author name and birthdate", func(t *testing.T) {
		author, err := uc.Create(context.Background(), &model.CreateAuthorRequest{
			Name:      "Author Name",
			Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateAuthorRequest{
				ID:        author.ID,
				Name:      util.ToPointer("Author Name Changed"),
				Birthdate: util.ToPointer(time.Date(2012, 1, 11, 1, 11, 11, 1111, time.UTC)),
			},
		}
		returned := returned{
			data: &model.AuthorResponse{
				ID:        author.ID,
				Name:      *params.request.Name,
				Birthdate: *params.request.Birthdate,
			},
			err: nil,
		}

		resp, err := uc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.UpdateAuthorRequest{
				ID:        0,
				Name:      util.ToPointer("Author Name Changed"),
				Birthdate: util.ToPointer(time.Date(2012, 1, 11, 1, 11, 11, 1111, time.UTC)),
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("id not found")),
		}

		resp, err := uc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 2 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.UpdateAuthorRequest{
				ID:        0,
				Name:      util.ToPointer("Author Name Changed"),
				Birthdate: util.ToPointer(time.Date(2012, 1, 11, 1, 11, 11, 1111, time.UTC)),
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to find author data by id")),
		}

		resp, err := failUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}

func TestAuthorUsecase_Delete(t *testing.T) {
	uc := newAuthorUsecase()
	failUc := newFailAuthorUsecase()

	type params struct {
		ctx     context.Context
		request *model.DeleteAuthorRequest
	}
	type returned struct {
		data *int
		err  error
	}

	t.Run("Positive Case - delete by id", func(t *testing.T) {
		author, err := uc.Create(context.Background(), &model.CreateAuthorRequest{
			Name:      "Author Name",
			Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.DeleteAuthorRequest{
				ID: author.ID,
			},
		}
		returned := returned{
			data: &author.ID,
			err:  nil,
		}

		resp, err := uc.Delete(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.DeleteAuthorRequest{
				ID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("id not found")),
		}

		resp, err := uc.Delete(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 2 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.DeleteAuthorRequest{
				ID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to find author data by id")),
		}

		resp, err := failUc.Delete(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}
