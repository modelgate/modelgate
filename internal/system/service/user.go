package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"

	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) CreateUser(ctx context.Context, req *model.CreateUserRequest) (user *model.User, err error) {
	if !model.UsernameReg.MatchString(req.User.Username) {
		err = errors.Errorf("invalid username: %s", req.User.Username)
		return
	}
	total, err := s.userDao.Count(ctx, &model.UserFilter{Username: db.Eq(req.User.Username)})
	if err != nil {
		return
	} else if total > 0 {
		err = errors.Errorf("username %s already exist", req.User.Username)
		return
	}
	var passwordHash []byte
	if req.User.Password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
		if err != nil {
			err = errors.Errorf("failed to generate password hash: %v", err)
			return
		}
	}
	user = &model.User{
		Username: req.User.Username,
		Nickname: req.User.Nickname,
		Email:    req.User.Email,
		Phone:    req.User.Phone,
		Gender:   req.User.Gender,
		Roles:    strings.Join(req.User.Roles, ","),
		Status:   model.EnableStatusEnabled,
		Password: string(passwordHash),
	}
	if err = s.userDao.Create(ctx, user); err != nil {
		err = errors.Errorf("failed to create user: %v", err)
		return
	}
	return
}

func (s *Service) UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (user *model.User, err error) {
	user, err = s.userDao.FindOneByID(ctx, req.User.Id)
	if err != nil {
		return
	}

	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "nickname") {
		update["nickname"] = req.User.Nickname
	}
	if lo.Contains(req.UpdateMask, "email") {
		update["email"] = req.User.Email
	}
	if lo.Contains(req.UpdateMask, "phone") {
		update["phone"] = req.User.Phone
	}
	if lo.Contains(req.UpdateMask, "gender") {
		update["gender"] = req.User.Gender
	}
	if lo.Contains(req.UpdateMask, "roles") {
		update["roles"] = strings.Join(req.User.Roles, ",")
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.User.Status
	}
	if lo.Contains(req.UpdateMask, "password") && req.User.Password != "" {
		var passwordHash []byte
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
		if err != nil {
			err = errors.Errorf("failed to generate password hash: %v", err)
			return
		}
		update["password"] = string(passwordHash)
	}
	if len(update) == 0 {
		return
	}
	if err = s.userDao.UpdateOne(ctx, user, update); err != nil {
		err = errors.Errorf("failed to update user: %v", err)
		return
	}
	return
}

func (s *Service) GetUser(ctx context.Context, req *model.GetUserRequest) (user *model.User, err error) {
	f := &model.UserFilter{
		ID: db.Eq(req.ID),
	}
	return s.userDao.FindOne(ctx, f)
}

func (s *Service) GetUserList(ctx context.Context, req *model.GetUserListRequest) (total int64, list []*model.User, err error) {
	filter := &model.UserFilter{
		Username: db.Eq(req.Username, db.OmitIfZero[string]()),
		Phone:    db.Eq(req.Phone, db.OmitIfZero[string]()),
		Nickname: db.Eq(req.Nickname, db.OmitIfZero[string]()),
		Gender:   db.Eq(req.Gender, db.OmitIfZero[string]()),
		Email:    db.Eq(req.Email, db.OmitIfZero[string]()),
		Status:   db.Eq(req.Status, db.OmitIfZero[model.EnableStatus]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.userDao.Count(ctx, filter)
		if err != nil {
			return
		}
		if !db.HasRecrods(total, req.PageParam.Page, req.PageParam.PageSize) {
			return
		}
		options = append(options,
			db.WithPaging(req.PageParam.Page, req.PageParam.PageSize),
			db.WithOrder(req.PageParam.OrderBy, nil))
	}
	list, err = s.userDao.Find(ctx, filter, options...)
	return
}

func (s *Service) DeleteUsers(ctx context.Context, req *model.DeleteUsersRequest) (err error) {
	if len(req.IDs) == 0 {
		err = errors.New("ids is empty")
		return
	}
	_, err = s.userDao.Delete(ctx, &model.UserFilter{IDs: db.In(req.IDs)})
	return err
}

func (s *Service) GetUserPermissions(ctx context.Context, userId int64) (isSuperAdmin bool, perms []*model.RolePermission, err error) {
	user, err := s.userDao.FindOneByID(ctx, userId)
	if err != nil {
		return
	}
	roleIds := lo.FlatMap(strings.Split(user.Roles, ","), func(s string, _ int) []int64 {
		id, err := strconv.ParseInt(s, 10, 64)
		return lo.Ternary(err == nil, []int64{id}, nil)
	})
	roleList, err := s.roleDao.Find(ctx, &model.RoleFilter{IDs: db.In(roleIds), Status: db.Eq(model.EnableStatusEnabled)})
	if err != nil {
		return
	}
	for _, role := range roleList {
		if role.IsSuperAdmin {
			isSuperAdmin = true
			break
		}
	}
	perms = lo.FlatMap(roleList, func(role *model.Role, _ int) []*model.RolePermission {
		perm, err := role.GetPermission()
		return lo.Ternary(err == nil, []*model.RolePermission{perm}, nil)
	})
	return
}

func (s *Service) GetUserApiPermissions(ctx context.Context, userId int64) (isSuperAdmin bool, apiPermList []*model.ApiPerm, err error) {
	isSuperAdmin, rolePerms, err := s.GetUserPermissions(ctx, userId)
	if err != nil {
		return
	}
	if isSuperAdmin {
		return
	}
	permCodes := lo.Uniq(lo.FlatMap(rolePerms, func(perm *model.RolePermission, _ int) []string {
		return perm.Buttons
	}))
	permissionList, err := s.permissionDao.Find(ctx, &model.PermissionFilter{Codes: db.In(permCodes)})
	if err != nil {
		return
	}
	apiPermList = lo.FlatMap(permissionList, func(p *model.Permission, _ int) []*model.ApiPerm {
		var list []*model.ApiPerm
		if err := json.Unmarshal(p.Data, &list); err != nil {
			return nil
		}
		return list
	})
	return
}
