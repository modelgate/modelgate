package model

import (
	"time"

	"github.com/modelgate/modelgate/pkg/db"
)

// RefreshToken 刷新令牌
type RefreshToken struct {
	db.Model

	UserId      int64     `gorm:"type:bigint;not null;default:0"`
	Jti         string    `gorm:"type:varchar(100);not null;default:''"`
	ExpiresAt   time.Time `gorm:"type:datetime;not null;"`
	Description string    `gorm:"type:varchar(1000);not null;default:''"`
}

func (RefreshToken) TableName() string {
	return TableRefreshTokens
}

type RefreshTokenFilter struct {
	ID     db.F[int64]
	UserId db.F[int64]
	Jti    db.F[string]
}
