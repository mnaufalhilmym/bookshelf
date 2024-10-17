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

func newAuthorAndBookUsecase() (*usecase.AuthorUsecase, *usecase.BookUsecase) {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	authorRepo := repository.NewAuthorRepository(db)
	bookRepo := repository.NewBookRepository(db)
	authorUc := usecase.NewAuthorUsecase(db, authorRepo)
	bookUc := usecase.NewBookUsecase(db, bookRepo, authorRepo)
	return authorUc, bookUc
}

func newFailAuthorAndBookUsecase() (*usecase.AuthorUsecase, *usecase.BookUsecase) {
	db := config.NewDatabase(":memory:", 1, 1, 100)
	authorRepo := &repository.AuthorRepository{}
	bookRepo := &repository.BookRepository{}
	authorUc := usecase.NewAuthorUsecase(db, authorRepo)
	bookUc := usecase.NewBookUsecase(db, bookRepo, authorRepo)
	return authorUc, bookUc
}

func TestBookUsecase_GetMany(t *testing.T) {
	authorUc, bookUc := newAuthorAndBookUsecase()
	_, bookFailUc := newFailAuthorAndBookUsecase()

	type params struct {
		ctx     context.Context
		request *model.GetManyBooksRequest
	}
	type returned struct {
		data  []model.BookResponse
		total int64
		err   error
	}

	author1, err := authorUc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
	})
	assert.NoError(t, err)
	author2, err := authorUc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 2",
		Birthdate: time.Date(2022, 2, 22, 2, 22, 22, 2222, time.UTC),
	})
	assert.NoError(t, err)

	book1, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
		Title:    "Book Title 1",
		ISBN:     "978-0451524935",
		AuthorID: author1.ID,
	})
	assert.NoError(t, err)
	book2, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
		Title:    "Book Title 2",
		ISBN:     "978-0743273565",
		AuthorID: author1.ID,
	})
	assert.NoError(t, err)
	book3, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
		Title:    "Book Title 3",
		ISBN:     "979-0061120084",
		AuthorID: author2.ID,
	})
	assert.NoError(t, err)

	t.Run("Positive Case 1 - search by title", func(t *testing.T) {
		request := &model.GetManyBooksRequest{
			Title: util.ToPointer("book title"),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.BookResponse{*book1, *book2, *book3},
			total: 3,
			err:   nil,
		}

		resp, total, err := bookUc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Positive Case 2 - search by isbn", func(t *testing.T) {
		request := &model.GetManyBooksRequest{
			ISBN: util.ToPointer("978"),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.BookResponse{*book1, *book2},
			total: 2,
			err:   nil,
		}

		resp, total, err := bookUc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Positive Case 3 - search by author id", func(t *testing.T) {
		request := &model.GetManyBooksRequest{
			AuthorID: &author2.ID,
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.BookResponse{*book3},
			total: 1,
			err:   nil,
		}

		resp, total, err := bookUc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Positive Case 4 - search by author name", func(t *testing.T) {
		request := &model.GetManyBooksRequest{
			AuthorName: util.ToPointer("name 1"),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.BookResponse{*book1, *book2},
			total: 2,
			err:   nil,
		}

		resp, total, err := bookUc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Positive Case 5 - empty search result", func(t *testing.T) {
		request := &model.GetManyBooksRequest{
			ISBN: util.ToPointer("zzz"),
		}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  []model.BookResponse{},
			total: 0,
			err:   nil,
		}

		resp, total, err := bookUc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})

	t.Run("Negative Case - db error", func(t *testing.T) {
		request := &model.GetManyBooksRequest{}
		request.Page = 1
		request.Size = 10

		params := params{
			ctx:     context.Background(),
			request: request,
		}
		returned := returned{
			data:  nil,
			total: 0,
			err:   model.InternalServerError(errors.New("failed to get many books")),
		}

		resp, total, err := bookFailUc.GetMany(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
		assert.EqualValues(t, returned.total, total)
	})
}

func TestBookUsecase_Get(t *testing.T) {
	authorUc, bookUc := newAuthorAndBookUsecase()
	_, bookFailUc := newFailAuthorAndBookUsecase()

	type params struct {
		ctx     context.Context
		request *model.GetBookRequest
	}
	type returned struct {
		data *model.BookResponse
		err  error
	}

	author, err := authorUc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
	})
	assert.NoError(t, err)

	book, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
		Title:    "Book Title 1",
		ISBN:     "978-0451524935",
		AuthorID: author.ID,
	})
	assert.NoError(t, err)

	t.Run("Positive Case - get by id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.GetBookRequest{
				ID: book.ID,
			},
		}
		returned := returned{
			data: book,
			err:  nil,
		}

		resp, err := bookUc.Get(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.GetBookRequest{
				ID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("book not found")),
		}

		resp, err := bookUc.Get(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 2 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.GetBookRequest{
				ID: 1,
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to find book data by id")),
		}

		resp, err := bookFailUc.Get(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}

func TestBookUsecase_Create(t *testing.T) {
	authorUc, bookUc := newAuthorAndBookUsecase()
	_, bookFailUc := newFailAuthorAndBookUsecase()

	type params struct {
		ctx     context.Context
		request *model.CreateBookRequest
	}
	type returned struct {
		data *model.BookResponse
		err  error
	}

	author, err := authorUc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
	})
	assert.NoError(t, err)

	t.Run("Positive Case - create book", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.CreateBookRequest{
				Title:    "Book Title",
				ISBN:     "978-0451524935",
				AuthorID: author.ID,
			},
		}
		returned := returned{
			data: &model.BookResponse{
				ID:         1,
				Title:      params.request.Title,
				ISBN:       params.request.ISBN,
				AuthorID:   author.ID,
				AuthorName: author.Name,
			},
			err: nil,
		}

		resp, err := bookUc.Create(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 1 - author not found", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.CreateBookRequest{
				Title:    "Book Title",
				ISBN:     "978-0451524935",
				AuthorID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("author not found")),
		}

		resp, err := bookUc.Create(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 2 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.CreateBookRequest{
				Title:    "Book Title",
				ISBN:     "978-0451524935",
				AuthorID: author.ID,
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to find author data by id")),
		}

		resp, err := bookFailUc.Create(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}

func TestBookUsecase_Update(t *testing.T) {
	authorUc, bookUc := newAuthorAndBookUsecase()
	_, bookFailUc := newFailAuthorAndBookUsecase()

	type params struct {
		ctx     context.Context
		request *model.UpdateBookRequest
	}
	type returned struct {
		data *model.BookResponse
		err  error
	}

	author1, err := authorUc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
	})
	assert.NoError(t, err)
	author2, err := authorUc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 2",
		Birthdate: time.Date(2022, 2, 22, 2, 22, 22, 2222, time.UTC),
	})
	assert.NoError(t, err)

	t.Run("Positive Case 1 - update book title", func(t *testing.T) {
		book, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
			Title:    "Book Title",
			ISBN:     "978-0451524935",
			AuthorID: author1.ID,
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateBookRequest{
				ID:    book.ID,
				Title: util.ToPointer("Book Title Changed"),
			},
		}
		returned := returned{
			data: &model.BookResponse{
				ID:         book.ID,
				Title:      *params.request.Title,
				ISBN:       book.ISBN,
				AuthorID:   book.AuthorID,
				AuthorName: book.AuthorName,
			},
			err: nil,
		}

		resp, err := bookUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Positive Case 2 - update book isbn", func(t *testing.T) {
		book, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
			Title:    "Book Title",
			ISBN:     "978-0316769488",
			AuthorID: author1.ID,
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateBookRequest{
				ID:   book.ID,
				ISBN: util.ToPointer("978-0060850524"),
			},
		}
		returned := returned{
			data: &model.BookResponse{
				ID:         book.ID,
				Title:      book.Title,
				ISBN:       *params.request.ISBN,
				AuthorID:   book.AuthorID,
				AuthorName: book.AuthorName,
			},
			err: nil,
		}

		resp, err := bookUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Positive Case 3 - update book author", func(t *testing.T) {
		book, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
			Title:    "Book Title",
			ISBN:     "978-0547928227",
			AuthorID: author1.ID,
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateBookRequest{
				ID:       book.ID,
				AuthorID: &author2.ID,
			},
		}
		returned := returned{
			data: &model.BookResponse{
				ID:         book.ID,
				Title:      book.Title,
				ISBN:       book.ISBN,
				AuthorID:   author2.ID,
				AuthorName: author2.Name,
			},
			err: nil,
		}

		resp, err := bookUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Positive Case 4 - update all book properties", func(t *testing.T) {
		book, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
			Title:    "Book Title",
			ISBN:     "978-1451673319",
			AuthorID: author1.ID,
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateBookRequest{
				ID:       book.ID,
				Title:    util.ToPointer("Book Title Changed"),
				ISBN:     util.ToPointer("978-1503290563"),
				AuthorID: &author2.ID,
			},
		}
		returned := returned{
			data: &model.BookResponse{
				ID:         book.ID,
				Title:      *params.request.Title,
				ISBN:       *params.request.ISBN,
				AuthorID:   author2.ID,
				AuthorName: author2.Name,
			},
			err: nil,
		}

		resp, err := bookUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 1 - wrong book id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.UpdateBookRequest{
				ID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("id not found")),
		}

		resp, err := bookUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 2 - wrong author id", func(t *testing.T) {
		book, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
			Title:    "Book Title",
			ISBN:     "978-0062315007",
			AuthorID: author1.ID,
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.UpdateBookRequest{
				ID:       book.ID,
				AuthorID: util.ToPointer(0),
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("id not found")),
		}

		resp, err := bookUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 3 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.UpdateBookRequest{
				ID: 1,
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to find book data by id")),
		}

		resp, err := bookFailUc.Update(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}

func TestBookUsecase_Delete(t *testing.T) {
	authorUc, bookUc := newAuthorAndBookUsecase()
	_, bookFailUc := newFailAuthorAndBookUsecase()

	type params struct {
		ctx     context.Context
		request *model.DeleteBookRequest
	}
	type returned struct {
		data *int
		err  error
	}

	author, err := authorUc.Create(context.Background(), &model.CreateAuthorRequest{
		Name:      "Author Name 1",
		Birthdate: time.Date(2011, 1, 11, 1, 11, 11, 1111, time.UTC),
	})
	assert.NoError(t, err)

	t.Run("Positive Case - delete by id", func(t *testing.T) {
		book, err := bookUc.Create(context.Background(), &model.CreateBookRequest{
			Title:    "Book Title",
			ISBN:     "978-0439708180",
			AuthorID: author.ID,
		})
		assert.NoError(t, err)

		params := params{
			ctx: context.Background(),
			request: &model.DeleteBookRequest{
				ID: book.ID,
			},
		}
		returned := returned{
			data: &book.ID,
			err:  nil,
		}

		resp, err := bookUc.Delete(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 1 - wrong id", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.DeleteBookRequest{
				ID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.NotFound(errors.New("id not found")),
		}

		resp, err := bookUc.Delete(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})

	t.Run("Negative Case 2 - db error", func(t *testing.T) {
		params := params{
			ctx: context.Background(),
			request: &model.DeleteBookRequest{
				ID: 0,
			},
		}
		returned := returned{
			data: nil,
			err:  model.InternalServerError(errors.New("failed to find book data by id")),
		}

		resp, err := bookFailUc.Delete(params.ctx, params.request)
		assert.EqualValues(t, returned.err, err)
		assert.EqualValues(t, returned.data, resp)
	})
}
