package model

import (
	"strings"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	systempb "github.com/modelgate/modelgate/pkg/proto/model/system"
	"github.com/modelgate/modelgate/pkg/types"
)

// User 用户
type User struct {
	db.Model

	Username    string       `gorm:"type:varchar(100);not null;default:'';uniqueIndex:uk_username"`
	Nickname    string       `gorm:"type:varchar(100);not null;default:''"`
	Phone       string       `gorm:"type:varchar(20);not null;default:''"`
	Email       string       `gorm:"type:varchar(100);not null;default:''"`
	Gender      string       `gorm:"type:enum('male', 'female', 'unknown')	;not null;default:'unknown'"`
	Roles       string       `gorm:"type:varchar(2000);not null;default:''"`
	Password    string       `gorm:"type:varchar(200);not null;default:''"`
	Status      EnableStatus `gorm:"type:enum('enabled', 'disabled');not null;default:'enabled'"`
	AvatarUrl   string       `gorm:"type:varchar(1000);not null;default:''"`
	Description string       `gorm:"type:varchar(1000);not null;default:''"`
}

func (User) TableName() string {
	return TableUsers
}

func (u *User) ToProto() *systempb.User {
	return &systempb.User{
		Id:          u.ID,
		Username:    u.Username,
		Nickname:    u.Nickname,
		Email:       u.Email,
		Phone:       u.Phone,
		Gender:      u.Gender,
		Roles:       lo.Filter(strings.Split(u.Roles, ","), func(role string, _ int) bool { return role != "" }),
		Status:      string(u.Status),
		AvatarUrl:   u.AvatarUrl,
		Description: u.Description,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}
}

type UserFilter struct {
	ID       db.F[int64]
	IDs      db.F[[]int64] `gorm:"column:id"`
	Username db.F[string]
	Phone    db.F[string]
	Nickname db.F[string]
	Gender   db.F[string]
	Email    db.F[string]
	Status   db.F[EnableStatus]
}

type CreateUserRequest struct {
	User *systempb.User
}

type UpdateUserRequest struct {
	User       *systempb.User
	UpdateMask []string
}

type DeleteUsersRequest struct {
	IDs []int64
}

type GetUserRequest struct {
	ID int64
}

type GetUserListRequest struct {
	*types.PageParam

	Username string
	Phone    string
	Nickname string
	Gender   string
	Email    string
	Status   EnableStatus
}
