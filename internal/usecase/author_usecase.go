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

type AuthorUsecase struct {
	db         *gorm.DB
	repository *repository.AuthorRepository
}

func NewAuthorUsecase(
	db *gorm.DB,
	repository *repository.AuthorRepository,
) *AuthorUsecase {
	return &AuthorUsecase{
		db,
		repository,
	}
}

func (uc *AuthorUsecase) GetMany(ctx context.Context, request *model.GetManyAuthorsRequest) ([]model.AuthorResponse, int64, error) {
	tx := uc.db.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	defer tx.Rollback()

	authors, total, err := uc.repository.Search(
		ctx,
		tx,
		request.Name,
		request.BirthdateStart,
		request.BirthdateEnd,
		request.Page,
		request.Size,
	)
	if err != nil {
		return nil, 0, model.InternalServerError(errors.New("failed to get many authors"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, 0, model.InternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToAuthorsResponse(authors), total, nil
}

func (uc *AuthorUsecase) Get(ctx context.Context, request *model.GetAuthorRequest) (*model.AuthorResponse, error) {
	tx := uc.db.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	defer tx.Rollback()

	author, err := uc.repository.FindByID(tx, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.BadRequest(errors.New("author not found"))
		}
		return nil, model.InternalServerError(errors.New("failed to find author data by id"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.InternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToAuthorResponse(author), nil
}

func (uc *AuthorUsecase) Create(ctx context.Context, request *model.CreateAuthorRequest) (*model.AuthorResponse, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	author := &entity.Author{
		Name:      request.Name,
		Birthdate: request.Birthdate,
	}

	if err := uc.repository.Create(tx, author); err != nil {
		return nil, model.InternalServerError(errors.New("failed to create new author"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.InternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToAuthorResponse(author), nil
}

func (uc *AuthorUsecase) Update(ctx context.Context, request *model.UpdateAuthorRequest) (*model.AuthorResponse, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	author, err := uc.repository.FindByID(tx, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.BadRequest(errors.New("id not found"))
		}
		return nil, model.InternalServerError(errors.New("failed to find author data by id"))
	}

	if request.Name != nil && *request.Name != "" {
		author.Name = *request.Name
	}

	if request.Birthdate != nil {
		author.Birthdate = *request.Birthdate
	}

	if err := uc.repository.Update(tx, author); err != nil {
		return nil, model.InternalServerError(errors.New("failed to update author"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.InternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToAuthorResponse(author), nil
}

func (uc *AuthorUsecase) Delete(ctx context.Context, request *model.DeleteAuthorRequest) (*int, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	author, err := uc.repository.FindByID(tx, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.BadRequest(errors.New("id not found"))
		}
		return nil, model.InternalServerError(errors.New("failed to find author data by id"))
	}

	if err := uc.repository.Delete(tx, author); err != nil {
		return nil, model.InternalServerError(errors.New("failed to delete author"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.InternalServerError(errors.New("failed to commit transaction"))
	}

	return &author.ID, nil
}
