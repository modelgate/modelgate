package model

import (
	"time"

	"github.com/modelgate/modelgate/pkg/db"
)

type DataMigration struct {
	ID          int64     `gorm:"type:bigint unsigned;comment:主键ID" json:"id,string"`
	Version     string    `gorm:"type:varchar(200);not null;default:''"`   // 版本号
	Description string    `gorm:"type:varchar(200);not null;default:''"`   // 描述
	ExecutedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 执行时间

	FilePath string `gorm:"-"` // 文件路径
	RawSql   string `gorm:"-"` // sql 语句
}

func (DataMigration) TableName() string {
	return TableDataMigrations
}

type DataMigrationFilter struct {
	ID      db.F[int64]
	Version db.F[string]
}
