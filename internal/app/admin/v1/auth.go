package v1

import (
	"context"
	"net/http"

	"connectrpc.com/connect/v2"
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

func (s *AuthService) Login(ctx context.Context, req *v1pb.LoginRequest) (resp *v1pb.LoginResponse, err error) {
	info, _ := connect.CallInfoForServerContext(ctx)
	result, err := s.systemService.Login(ctx, &model.LoginRequest{
		Username:   req.Username,
		Password:   req.Password,
		RememberMe: req.RememberMe,
	})
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.LoginResponse{
		AccessToken: result.AccessToken,
	}
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/api/admin.v1.AuthService",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		// Secure:   true, // 只有https才需要
	}
	info.ResponseHeader().Add("Set-Cookie", cookie.String())
	return resp, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *v1pb.RefreshTokenRequest) (resp *v1pb.RefreshTokenResponse, err error) {
	info, _ := connect.CallInfoForServerContext(ctx)
	cookieHeader := info.RequestHeader().Get("Cookie")
	header := http.Header{"Cookie": []string{cookieHeader}}
	r := &http.Request{Header: header}
	tokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, "missing refresh token")
		return
	}
	result, err := s.systemService.RefreshToken(ctx, &model.RefreshTokenRequest{
		RefreshToken: tokenCookie.Value,
	})
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.RefreshTokenResponse{
		AccessToken: result.AccessToken,
	}
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/api/admin.v1.AuthService",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		// Secure:   true, // 只有https才需要
	}
	info.ResponseHeader().Add("Set-Cookie", cookie.String())
	return resp, nil
}

func (s *AuthService) GetUserInfo(ctx context.Context, req *v1pb.GetUserInfoRequest) (resp *v1pb.GetUserInfoResponse, err error) {
	userId, ok := system.UserIDFromContext(ctx)
	if !ok {
		err = connect.NewError(connect.CodeUnauthenticated, "invalid access token")
		return
	}
	user, err := s.systemService.GetUser(ctx, &model.GetUserRequest{ID: userId})
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err.Error()).WithCause(err)
		return
	}
	if user.Status == model.EnableStatusDisabled {
		err = connect.NewError(connect.CodeUnauthenticated, "user has been archived")
		return
	}
	isSuperAdmin, perms, err := s.systemService.GetUserPermissions(ctx, userId)
	if err != nil {
		err = connect.NewError(connect.CodeUnauthenticated, err.Error()).WithCause(err)
		return
	}
	buttons := lo.FlatMap(perms, func(permission *model.RolePermission, index int) []string {
		return permission.Buttons
	})

	resp = &v1pb.GetUserInfoResponse{
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
	}
	return resp, nil
}
