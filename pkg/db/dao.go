package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type BaseDAO[T schema.Tabler, F any] struct {
	db *gorm.DB
}

func NewBaseDAO[T schema.Tabler, F any](db *gorm.DB) *BaseDAO[T, F] {
	return &BaseDAO[T, F]{db: db}
}

func (d *BaseDAO[T, _]) Create(ctx context.Context, m *T, opts ...Option) error {
	tx := d.db.WithContext(ctx)
	for _, opt := range opts {
		tx = opt(tx)
	}
	return tx.Create(m).Error
}

func (d *BaseDAO[T, _]) Save(ctx context.Context, m *T, opts ...Option) error {
	tx := d.db.WithContext(ctx)
	for _, opt := range opts {
		tx = opt(tx)
	}
	return tx.Save(m).Error
}

func (d *BaseDAO[T, F]) Update(ctx context.Context, filter *F, update map[string]any) (int64, error) {
	tx := d.db.WithContext(ctx).Model(new(T))
	tx = d.applyFilter(tx, filter)
	result := tx.Updates(update)
	return result.RowsAffected, result.Error
}

func (d *BaseDAO[T, F]) UpdateOne(ctx context.Context, m *T, update map[string]any) error {
	return d.db.WithContext(ctx).Model(m).Updates(update).Error
}

func (d *BaseDAO[T, F]) Count(ctx context.Context, filter *F) (int64, error) {
	var count int64
	tx := d.db.WithContext(ctx).Model(new(T))
	tx = d.applyFilter(tx, filter)
	err := tx.Count(&count).Error
	return count, err
}

func (d *BaseDAO[T, F]) Find(ctx context.Context, filter *F, opts ...Option) ([]*T, error) {
	var list []*T
	tx := d.db.WithContext(ctx).Model(new(T))
	tx = d.applyFilter(tx, filter)
	for _, opt := range opts {
		opt(tx)
	}
	err := tx.Find(&list).Error
	return list, err
}

func (d *BaseDAO[T, F]) FindOneByID(ctx context.Context, id int64) (*T, error) {
	var m T
	err := d.db.WithContext(ctx).First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *BaseDAO[T, F]) FindOne(ctx context.Context, filter *F, opts ...Option) (*T, error) {
	var m T
	tx := d.db.WithContext(ctx).Model(&m)
	tx = d.applyFilter(tx, filter)
	for _, opt := range opts {
		opt(tx)
	}
	err := tx.First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *BaseDAO[T, F]) Delete(ctx context.Context, filter *F) (int64, error) {
	tx := d.db.WithContext(ctx)
	tx = d.applyFilter(tx, filter)
	result := tx.Delete(new(T))
	return result.RowsAffected, result.Error
}

func (d *BaseDAO[_, _]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn)
}

func (d *BaseDAO[T, F]) applyFilter(tx *gorm.DB, filter *F) *gorm.DB {
	if filter == nil {
		return tx
	}
	where, args, err := BuildQuery(filter)
	if err != nil {
		tx.AddError(err)
		return tx
	}
	if where != "" {
		tx = tx.Where(where, args...)
	}
	return tx
}
