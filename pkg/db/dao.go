package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// BaseDAO provides common CRUD operations for a model M and filter F.
type BaseDAO[M schema.Tabler, F any] struct {
	db *gorm.DB
}

func NewBaseDAO[M schema.Tabler, F any](conn *gorm.DB) *BaseDAO[M, F] {
	return &BaseDAO[M, F]{db: conn}
}

func (d *BaseDAO[M, F]) GetDB() *gorm.DB {
	return d.db
}

func (d *BaseDAO[M, _]) Create(ctx context.Context, m *M) error {
	return d.db.Create(m).Error
}

func (d *BaseDAO[M, _]) Save(ctx context.Context, m *M) error {
	return d.db.Save(m).Error
}

func (d *BaseDAO[M, F]) Update(ctx context.Context, filter *F, update map[string]any) (int64, error) {
	var m M
	res := Apply(d.db, WithFilter(filter)).Model(m).Updates(update)
	return res.RowsAffected, res.Error
}

func (d *BaseDAO[M, _]) UpdateOne(ctx context.Context, m *M, update map[string]any) error {
	res := Apply(d.db, WithFilter(nil)).Model(m).Updates(update)
	return res.Error
}

func (d *BaseDAO[M, F]) Count(ctx context.Context, f *F) (total int64, err error) {
	var m M
	err = Apply(d.db, WithFilter(f)).Model(&m).Count(&total).Error
	return
}

func (d *BaseDAO[M, F]) Find(ctx context.Context, f *F, opts ...Option) (ms []*M, err error) {
	err = Apply(d.db, WithFilter(f), opts...).Find(&ms).Error
	return
}

func (d *BaseDAO[M, F]) FindOne(ctx context.Context, f *F, opts ...Option) (*M, error) {
	var m M
	if err := Apply(d.db, WithFilter(f), opts...).Take(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *BaseDAO[M, _]) FindOneByID(ctx context.Context, id int64) (*M, error) {
	var m M
	if err := d.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *BaseDAO[M, F]) Delete(ctx context.Context, filter *F) (int64, error) {
	var m M
	res := Apply(d.db, WithFilter(filter)).Delete(&m)
	return res.RowsAffected, res.Error
}
