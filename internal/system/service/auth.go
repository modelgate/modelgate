package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) Login(ctx context.Context, req *model.LoginRequest) (resp *model.LoginResponse, err error) {
	user, err := s.userDao.FindOne(ctx, &model.UserFilter{Username: db.Eq(req.Username)})
	if err != nil {
		err = errors.Errorf("failed to find user by username %s", req.Username)
		return
	}
	if user.Status == model.EnableStatusDisabled {
		err = errors.Errorf("user has been archived with username %s", req.Username)
		return
	}
	// Compare the stored hashed password, with the hashed version of the password that was received.
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err = errors.New("unmatched username and password")
		return
	}
	issuedAt := time.Now()
	expireTime := time.Now().Add(model.AccessTokenDuration)
	_, accessTokenStr, err := s.generateJwtToken(ctx, user.ID, issuedAt, expireTime)
	if err != nil {
		return
	}
	tokenExpireTime := time.Now().Add(model.DefaultRefreshTokenDuration)
	if req.RememberMe {
		tokenExpireTime = time.Now().Add(model.RememberMeRefreshTokenDuration)
	}
	jti, refreshTokenStr, err := s.generateJwtToken(ctx, user.ID, issuedAt, tokenExpireTime)
	if err != nil {
		return
	}
	refreshToken := &model.RefreshToken{
		UserId:      user.ID,
		Jti:         jti,
		Description: "login",
		ExpiresAt:   tokenExpireTime,
	}
	err = s.refreshTokenDao.Create(ctx, refreshToken)
	if err != nil {
		return
	}
	resp = &model.LoginResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}
	return
}

func (s *Service) generateJwtToken(_ context.Context, userID int64, issuedAt, expirationTime time.Time) (jti string, tokenStr string, err error) {
	cfg := config.GetConfig().JWT

	jti = uuid.New().String()
	registeredClaims := jwt.RegisteredClaims{
		ID:       jti,
		Issuer:   model.Issuer,
		Audience: jwt.ClaimStrings{model.AccessTokenAudienceName},
		IssuedAt: jwt.NewNumericDate(issuedAt),
		Subject:  fmt.Sprint(userID),
	}
	if !expirationTime.IsZero() {
		registeredClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	token.Header["kid"] = model.KeyID

	// Create the JWT string.
	tokenStr, err = token.SignedString([]byte(cfg.Key))
	return
}

func (s *Service) parseJwtToken(_ context.Context, tokenStr string) (jti string, userId int64, err error) {
	cfg := config.GetConfig().JWT

	claims := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.Errorf("unexpected access token signing method=%v, expect %v", t.Header["alg"], jwt.SigningMethodHS256)
		}
		if kid, ok := t.Header["kid"].(string); ok {
			if kid == model.KeyID {
				return []byte(cfg.Key), nil
			}
		}
		return nil, errors.Errorf("unexpected access token kid=%v", t.Header["kid"])
	})
	if err != nil {
		err = errors.New("Invalid or expired access token")
		return
	}
	if claims.Issuer != model.Issuer {
		err = errors.New("invalid access token issuer")
		return
	}
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		err = errors.New("access token expired")
		return
	}
	jti = claims.ID
	userId, err = strconv.ParseInt(claims.Subject, 10, 32)
	if err != nil {
		return
	}
	return
}

// RefreshToken refresh the access token and refresh token
func (s *Service) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (resp *model.RefreshTokenResponse, err error) {
	jti, userId, err := s.parseJwtToken(ctx, req.RefreshToken)
	if err != nil {
		return
	}
	refreshToken, err := s.refreshTokenDao.FindOne(ctx, &model.RefreshTokenFilter{
		Jti:    db.Eq(jti),
		UserId: db.Eq(userId),
	})
	if err != nil {
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		err = errors.New("refresh token expired")
		return
	}
	user, err := s.userDao.FindOneByID(ctx, refreshToken.UserId)
	if err != nil {
		return
	}
	if user.Status == model.EnableStatusDisabled {
		err = errors.New("user has been archived")
		return
	}

	issuedAt := time.Now()
	accessTokenExpireTime := time.Now().Add(model.AccessTokenDuration)
	_, accessTokenStr, err := s.generateJwtToken(ctx, user.ID, issuedAt, accessTokenExpireTime)
	if err != nil {
		return
	}
	jti, refreshTokenStr, err := s.generateJwtToken(ctx, user.ID, issuedAt, refreshToken.ExpiresAt)
	if err != nil {
		return
	}
	err = s.refreshTokenDao.UpdateOne(ctx, refreshToken, map[string]any{"jti": jti})
	if err != nil {
		return
	}
	resp = &model.RefreshTokenResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}
	return
}

func (s *Service) Authenticate(ctx context.Context, tokenStr string) (userId int64, err error) {
	if tokenStr == "" {
		err = errors.New("access token not found")
		return
	}
	_, userId, err = s.parseJwtToken(ctx, tokenStr)
	if err != nil {
		return
	}
	user, err := s.userDao.FindOneByID(ctx, userId)
	if err != nil {
		return
	} else if user.Status == model.EnableStatusDisabled {
		err = errors.New("user has been archived")
		return
	}
	return
}
