package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mnaufalhilmym/bookshelf/internal/entity"
	"github.com/mnaufalhilmym/goasync"
	"github.com/mnaufalhilmym/gotracing"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	repository[entity.Author]
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	if err := db.Migrator().CreateTable(&entity.Author{}); err != nil {
		panic(fmt.Errorf("failed to migrate entity: %w", err))
	}

	return &AuthorRepository{}
}

func (r *AuthorRepository) Search(
	ctx context.Context,
	db *gorm.DB,
	name *string,
	birthdateStart *time.Time,
	birthdateEnd *time.Time,
	page int,
	size int,
) ([]entity.Author, int64, error) {
	offset := 0
	if page > 0 {
		offset = (page - 1) * size
	}

	filter := r.searchFilter(name, birthdateStart, birthdateEnd)

	authorsTask := goasync.Spawn(func(ctx context.Context) (authors []entity.Author, err error) {
		err = db.Scopes(filter).Offset(offset).Limit(size).Find(&authors).Error
		return
	})

	totalTask := goasync.Spawn(func(ctx context.Context) (total int64, err error) {
		err = db.Model(&entity.Author{}).Scopes(filter).Count(&total).Error
		return
	})

	authors, err := authorsTask.Await(ctx)
	if err != nil {
		gotracing.Error("Failed to find entities from database", err)
		return nil, 0, err
	}

	total, err := totalTask.Await(ctx)
	if err != nil {
		gotracing.Error("Failed to count entities from database", err)
		return nil, 0, err
	}

	return authors, total, nil
}

func (*AuthorRepository) searchFilter(
	name *string,
	birthdateStart *time.Time,
	birthdateEnd *time.Time,
) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name != nil && *name != "" {
			fname := "%" + *name + "%"
			tx = tx.Where("LOWER(name) LIKE LOWER(?)", fname)
		}

		if birthdateStart != nil {
			tx = tx.Where("birthdate >= ?", *birthdateStart)
		}

		if birthdateEnd != nil {
			tx = tx.Where("birthdate <= ?", *birthdateEnd)
		}

		return tx
	}
}
