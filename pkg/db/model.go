package db

import (
	"time"
)

// Model 基础结构
type Model struct {
	ID        int64     `gorm:"type:bigint unsigned;comment:主键ID" json:"id,string"`                      // 主键ID
	CreatedAt time.Time `gorm:"type:datetime(3);not null;autoCreateTime;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"type:datetime(3);not null;autoUpdateTime;comment:更新时间" json:"updated_at"` // 更新时间
}
