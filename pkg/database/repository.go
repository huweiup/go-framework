package database

import "gorm.io/gorm"

// Repository provides a generic data access layer
type Repository[T any] struct {
	DB *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{DB: db}
}

func (r *Repository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *Repository[T]) FindByID(id interface{}) (*T, error) {
	var entity T
	if err := r.DB.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *Repository[T]) Delete(id interface{}) error {
	var entity T
	return r.DB.Delete(&entity, id).Error
}

func (r *Repository[T]) List(offset, limit int) ([]T, int64, error) {
	var entities []T
	var count int64

	if err := r.DB.Model(new(T)).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, count, nil
}
