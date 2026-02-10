package v1

import (
	"context"
	"errors"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"github.com/samber/do/v2"
	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	v1pb "github.com/modelgate/modelgate/pkg/proto/admin/v1"
	modelpb "github.com/modelgate/modelgate/pkg/proto/model/system"
)

type AuthService struct {
	v1pb.UnimplementedAuthServiceHandler
	systemService system.Service
}

func NewAuthService(i do.Injector) (*AuthService, error) {
	return &AuthService{
		systemService: do.MustInvoke[system.Service](i),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *connect.Request[v1pb.LoginRequest]) (resp *connect.Response[v1pb.LoginResponse], err error) {
	result, err := s.systemService.Login(ctx, &model.LoginRequest{
		Username:   req.Msg.Username,
		Password:   req.Msg.Password,
		RememberMe: req.Msg.RememberMe,
	})
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err)
		return
	}
	resp = connect.NewResponse(
		&v1pb.LoginResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
		},
	)
	return resp, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *connect.Request[v1pb.RefreshTokenRequest]) (resp *connect.Response[v1pb.RefreshTokenResponse], err error) {
	result, err := s.systemService.RefreshToken(ctx, &model.RefreshTokenRequest{
		RefreshToken: req.Msg.RefreshToken,
	})
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err)
		return
	}
	resp = connect.NewResponse(
		&v1pb.RefreshTokenResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
		},
	)
	return resp, nil
}

func (s *AuthService) GetUserInfo(ctx context.Context, req *connect.Request[v1pb.GetUserInfoRequest]) (resp *connect.Response[v1pb.GetUserInfoResponse], err error) {
	userId, ok := authn.GetInfo(ctx).(int64)
	if !ok {
		err = connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
		return
	}
	user, err := s.systemService.GetUser(ctx, &model.GetUserRequest{ID: userId})
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err)
		return
	}
	if user.Status == model.EnableStatusDisabled {
		err = connect.NewError(connect.CodeUnauthenticated, errors.New("user has been archived"))
		return
	}
	isSuperAdmin, perms, err := s.systemService.GetUserPermissions(ctx, userId)
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err)
		return
	}
	buttons := lo.FlatMap(perms, func(permission *model.RolePermission, index int) []string {
		return permission.Buttons
	})
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err)
		return
	}

	resp = connect.NewResponse(
		&v1pb.GetUserInfoResponse{
			User: &modelpb.UserInfo{
				Id:           user.ID,
				Username:     user.Username,
				Nickname:     user.Nickname,
				IsSuperAdmin: isSuperAdmin,
				Buttons:      buttons,
				Gender:       user.Gender,
				AvatarUrl:    user.AvatarUrl,
				Description:  user.Description,
			},
		},
	)
	return resp, nil
}
