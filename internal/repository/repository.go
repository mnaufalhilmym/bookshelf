package repository

import (
	"errors"

	"github.com/mnaufalhilmym/gotracing"
	"gorm.io/gorm"
)

type repository[T any] struct{}

func (*repository[T]) Create(db *gorm.DB, entity *T) error {
	if err := db.Create(entity).Error; err != nil {
		gotracing.Error("Failed to create entity to database", err)
		return err
	}
	return nil
}

func (*repository[T]) Update(db *gorm.DB, entity *T) error {
	if err := db.Save(entity).Error; err != nil {
		gotracing.Error("Failed to update entity to database", err)
		return err
	}
	return nil
}

func (*repository[T]) Delete(db *gorm.DB, entity *T) error {
	if err := db.Delete(entity).Error; err != nil {
		gotracing.Error("Failed to delete entity from database", err)
		return err
	}
	return nil
}

func (*repository[T]) FindByID(db *gorm.DB, id int) (*T, error) {
	var entity *T
	if err := db.Where("id = ?", id).First(&entity).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			gotracing.Error("Failed to find entity from database", err)
		}
		return nil, err
	}
	return entity, nil
}

func (*repository[T]) FindAll(db *gorm.DB) ([]T, error) {
	var entities []T
	if err := db.Find(&entities).Error; err != nil {
		gotracing.Error("Failed to find entity from database", err)
		return nil, err
	}
	return entities, nil
}
