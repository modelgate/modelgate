package service

import (
	"context"
	"testing"

	"github.com/samber/do/v2"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
	systempb "github.com/modelgate/modelgate/pkg/proto/model/system"
	"github.com/modelgate/modelgate/pkg/types"
)

func newTestService(t *testing.T) (do.Injector, *Service, *system.MockUserDAO) {
	i := do.New()
	ctl := gomock.NewController(t)
	userDaoMock := system.NewMockUserDAO(ctl)
	refreshTokenDaoMock := system.NewMockRefreshTokenDAO(ctl)
	roleDaoMock := system.NewMockRoleDAO(ctl)
	menuDaoMock := system.NewMockMenuDAO(ctl)

	do.Provide(i, func(i do.Injector) (system.UserDAO, error) { return userDaoMock, nil })
	do.Provide(i, func(i do.Injector) (system.RefreshTokenDAO, error) { return refreshTokenDaoMock, nil })
	do.Provide(i, func(i do.Injector) (system.RoleDAO, error) { return roleDaoMock, nil })
	do.Provide(i, func(i do.Injector) (system.MenuDAO, error) { return menuDaoMock, nil })

	s, _ := New(i)
	return i, s.(*Service), userDaoMock
}

func TestService_CreateUser(t *testing.T) {
	ctx := context.TODO()
	_, s, userDaoMock := newTestService(t)

	userDaoMock.
		EXPECT().
		Count(ctx, gomock.Any()).
		Return(int64(0), nil)

	userDaoMock.
		EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	user, err := s.CreateUser(ctx, &model.CreateUserRequest{
		User: &systempb.User{
			Username: "yearnfar",
			Password: "123456",
		},
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("user: %v", user)
}

func TestService_CreateUser_InvalidUsername(t *testing.T) {
	ctx := context.TODO()
	_, s, _ := newTestService(t)

	_, err := s.CreateUser(ctx, &model.CreateUserRequest{
		User: &systempb.User{
			Username: "!!",
		},
	})
	if err == nil {
		t.Fatal("expected error for invalid username")
	}
}

func TestService_CreateUser_UserExists(t *testing.T) {
	ctx := context.TODO()
	_, s, userDaoMock := newTestService(t)

	userDaoMock.EXPECT().
		Count(ctx, gomock.Any()).
		Return(int64(1), nil)

	_, err := s.CreateUser(ctx, &model.CreateUserRequest{
		User: &systempb.User{
			Username: "yearnfar",
		},
	})
	if err == nil {
		t.Fatal("expected error for existing username")
	}
}

func TestService_UpdateUser(t *testing.T) {
	ctx := context.TODO()
	_, s, userDaoMock := newTestService(t)

	userId := int64(1)
	existingUser := &model.User{
		Model: db.Model{
			ID: userId,
		},
		Username: "yearnfar",
		Nickname: "old_nick",
	}

	userDaoMock.EXPECT().
		FindOneByID(ctx, userId).
		Return(existingUser, nil)

	userDaoMock.EXPECT().
		UpdateOne(ctx, existingUser, gomock.Any()).
		Return(nil)

	_, err := s.UpdateUser(ctx, &model.UpdateUserRequest{
		User: &systempb.User{
			Id:       userId,
			Nickname: "new_nick",
		},
		UpdateMask: []string{"nickname"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_UpdateUser_NotFound(t *testing.T) {
	ctx := context.TODO()
	_, s, userDaoMock := newTestService(t)

	userId := int64(1)
	userDaoMock.EXPECT().
		FindOneByID(ctx, userId).
		Return(nil, gorm.ErrRecordNotFound)

	_, err := s.UpdateUser(ctx, &model.UpdateUserRequest{
		User: &systempb.User{
			Id: userId,
		},
	})
	if err == nil {
		t.Fatal("expected error for user not found")
	}
}

func TestService_GetUser(t *testing.T) {
	ctx := context.TODO()
	_, s, userDaoMock := newTestService(t)

	userId := int64(1)
	mockUser := &model.User{
		Model: db.Model{
			ID: userId,
		},
		Username: "yearnfar",
	}

	userDaoMock.EXPECT().
		FindOne(ctx, gomock.Any()).
		Return(mockUser, nil)

	user, err := s.GetUser(ctx, &model.GetUserRequest{ID: userId})
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != userId {
		t.Errorf("expected user id %d, got %d", userId, user.ID)
	}
}

func TestService_GetUserList(t *testing.T) {
	ctx := context.TODO()
	_, s, userDaoMock := newTestService(t)

	userDaoMock.EXPECT().
		Count(ctx, gomock.Any()).
		Return(int64(1), nil)

	userDaoMock.EXPECT().
		Find(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.User{{
			Model: db.Model{
				ID: 1,
			},
			Username: "yearnfar",
		}}, nil)

	total, list, err := s.GetUserList(ctx, &model.GetUserListRequest{
		PageParam: types.NewPageParam(1, 10, "id"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(list) != 1 {
		t.Errorf("expected list length 1, got %d", len(list))
	}
}

func TestService_DeleteUsers(t *testing.T) {
	ctx := context.TODO()
	_, s, userDaoMock := newTestService(t)

	ids := []int64{1, 2}

	userDaoMock.EXPECT().
		Delete(ctx, gomock.Any()).
		Return(int64(len(ids)), nil)

	err := s.DeleteUsers(ctx, &model.DeleteUsersRequest{IDs: ids})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_DeleteUsers_EmptyIDs(t *testing.T) {
	ctx := context.TODO()
	_, s, _ := newTestService(t)

	err := s.DeleteUsers(ctx, &model.DeleteUsersRequest{IDs: []int64{}})
	if err == nil {
		t.Fatal("expected error for empty ids")
	}
}
