package repository

import (
	"context"
	"fmt"

	"github.com/mnaufalhilmym/bookshelf/internal/entity"
	"github.com/mnaufalhilmym/goasync"
	"github.com/mnaufalhilmym/gotracing"
	"gorm.io/gorm"
)

type BookRepository struct {
	repository[entity.Book]
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	if err := db.Migrator().CreateTable(&entity.Book{}); err != nil {
		panic(fmt.Errorf("failed to migrate entity: %w", err))
	}

	return &BookRepository{}
}

func (r *BookRepository) Search(
	ctx context.Context,
	db *gorm.DB,
	title *string,
	isbn *string,
	authorID *int,
	authorName *string,
	page int,
	size int,
) ([]entity.Book, int64, error) {
	offset := 0
	if page > 0 {
		offset = (page - 1) * size
	}

	filter := r.searchFilter(title, isbn, authorID, authorName)

	booksTask := goasync.Spawn(func(ctx context.Context) (books []entity.Book, err error) {
		err = db.Joins("Author").Scopes(filter).Offset(offset).Limit(size).Find(&books).Error
		return
	})

	totalTask := goasync.Spawn(func(ctx context.Context) (total int64, err error) {
		err = db.Model(&entity.Book{}).Scopes(filter).Count(&total).Error
		return
	})

	books, err := booksTask.Await(ctx)
	if err != nil {
		gotracing.Error("Failed to find entities from database", err)
		return nil, 0, err
	}

	total, err := totalTask.Await(ctx)
	if err != nil {
		gotracing.Error("Failed to count entities from database", err)
		return nil, 0, err
	}

	return books, total, nil
}

func (*BookRepository) searchFilter(
	title *string,
	isbn *string,
	authorID *int,
	authorName *string,
) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if title != nil && *title != "" {
			ftitle := "%" + *title + "%"
			tx = tx.Where("books.title ILIKE ?", ftitle)
		}

		if isbn != nil && *isbn != "" {
			fisbn := "%" + *isbn + "%"
			tx = tx.Where("books.isbn ILIKE ?", fisbn)
		}

		if authorID != nil {
			tx = tx.Where("authors.id = ?", *authorID)
		}

		if authorName != nil && *authorName != "" {
			fname := "%" + *authorName + "%"
			tx = tx.Where("authors.name ILIKE ?", fname)
		}

		return tx
	}
}