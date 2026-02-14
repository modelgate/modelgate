package dao

import (
	"context"

	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type DataMigrationDao struct {
	*db.BaseDAO[model.DataMigration, model.DataMigrationFilter]
}

func NewDataMigrationDao(i do.Injector) (system.DataMigrationDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &DataMigrationDao{
		BaseDAO: db.NewBaseDAO[model.DataMigration, model.DataMigrationFilter](dbConn),
	}, nil
}

func (d *DataMigrationDao) Create(ctx context.Context, m *model.DataMigration) error {
	return d.GetDB().Transaction(func(tx *gorm.DB) error {
		// 使用原生 SQL 执行（支持多条语句和特殊语法）
		db, err := tx.DB()
		if err != nil {
			return err
		}
		if _, err := db.Exec(m.RawSql); err != nil {
			return err
		}
		return tx.Create(m).Error
	})
}

func (d *DataMigrationDao) AutoMigrate() error {
	return d.GetDB().AutoMigrate(&model.DataMigration{})
}
