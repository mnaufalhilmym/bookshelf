package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mnaufalhilmym/bookshelf/internal/entity"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
	"github.com/mnaufalhilmym/bookshelf/internal/repository"
	"github.com/mnaufalhilmym/gotracing"
	"gorm.io/gorm"
)

type BookUsecase struct {
	db               *gorm.DB
	repository       *repository.BookRepository
	authorRepository *repository.AuthorRepository
}

func NewBookUsecase(
	db *gorm.DB,
	repository *repository.BookRepository,
	authorRepository *repository.AuthorRepository,
) *BookUsecase {
	return &BookUsecase{
		db,
		repository,
		authorRepository,
	}
}

func (uc *BookUsecase) GetMany(ctx context.Context, request *model.GetManyBooksRequest) ([]model.BookResponse, int64, error) {
	tx := uc.db.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	defer tx.Rollback()

	books, total, err := uc.repository.Search(
		ctx,
		tx,
		request.Title,
		request.ISBN,
		request.AuthorID,
		request.AuthorName,
		request.Page,
		request.Size,
	)
	if err != nil {
		return nil, 0, model.ErrorInternalServerError(errors.New("failed to get many books"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, 0, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToBooksResponse(books), total, nil
}

func (uc *BookUsecase) Get(ctx context.Context, request *model.GetBookRequest) (*model.BookResponse, error) {
	tx := uc.db.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	defer tx.Rollback()

	book, err := uc.repository.FindByID(tx, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrorNotFound(errors.New("book not found"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to find book data by id"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToBookResponse(book), nil
}

func (uc *BookUsecase) Create(ctx context.Context, request *model.CreateBookRequest) (*model.BookResponse, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	author, err := uc.authorRepository.FindByID(tx, request.AuthorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrorNotFound(errors.New("author not found"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to find author data by id"))
	}

	book := &entity.Book{
		Title:    request.Title,
		ISBN:     request.ISBN,
		AuthorID: request.AuthorID,
		Author:   *author,
	}

	if err := uc.repository.Create(tx, book); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, model.ErrorBadRequest(errors.New("duplicate isbn"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to create new book"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToBookResponse(book), nil
}

func (uc *BookUsecase) Update(ctx context.Context, request *model.UpdateBookRequest) (*model.BookResponse, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	book, err := uc.repository.FindByID(tx, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrorNotFound(errors.New("id not found"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to find book data by id"))
	}

	if request.Title != nil && *request.Title != "" {
		book.Title = *request.Title
	}

	if request.ISBN != nil && *request.ISBN != "" {
		book.ISBN = *request.ISBN
	}

	if request.AuthorID != nil {
		author, err := uc.authorRepository.FindByID(tx, *request.AuthorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, model.ErrorNotFound(errors.New("id not found"))
			}
			return nil, model.ErrorInternalServerError(errors.New("failed to find author data by id"))
		}

		book.AuthorID = author.ID
		book.Author = *author
	}

	if err := uc.repository.Update(tx, book); err != nil {
		return nil, model.ErrorInternalServerError(errors.New("failed to update book data"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToBookResponse(book), nil
}

func (uc *BookUsecase) Delete(ctx context.Context, request *model.DeleteBookRequest) (*int, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	book, err := uc.repository.FindByID(tx, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrorNotFound(errors.New("id not found"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to find book data by id"))
	}

	if err := uc.repository.Delete(tx, book); err != nil {
		return nil, model.ErrorInternalServerError(errors.New("failed to delete author"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	return &book.ID, nil
}
