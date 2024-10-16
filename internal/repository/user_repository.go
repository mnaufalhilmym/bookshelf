package repository

import (
	"errors"
	"fmt"

	"github.com/mnaufalhilmym/bookshelf/internal/entity"
	"github.com/mnaufalhilmym/gotracing"
	"gorm.io/gorm"
)

type UserRepository struct {
	repository[entity.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	if err := db.Migrator().CreateTable(&entity.User{}); err != nil {
		panic(fmt.Errorf("failed to migrate entity: %w", err))
	}

	return &UserRepository{}
}

func (*UserRepository) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	var entity *entity.User
	if err := db.Where("username = ?", username).First(&entity).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			gotracing.Error("Failed to find entity from database", err)
		}
		return nil, err
	}
	return entity, nil
}
