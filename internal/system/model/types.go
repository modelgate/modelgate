package model

import (
	"regexp"
	"time"
)

const (
	// issuer is the issuer of the jwt token.
	Issuer = "modelgate"
	// Signing key section. For now, this is only used for signing, not for verifying since we only
	// have 1 version. But it will be used to maintain backward compatibility if we change the signing mechanism.
	KeyID = "v1"
	// AccessTokenAudienceName is the audience name of the access token.
	AccessTokenAudienceName        = "user.access-token"
	AccessTokenDuration            = 15 * time.Minute
	DefaultRefreshTokenDuration    = 7 * 24 * time.Hour
	RememberMeRefreshTokenDuration = 30 * 24 * time.Hour
)

const (
	TableUsers          = "users"
	TableAccessTokens   = "access_tokens"
	TableRefreshTokens  = "refresh_tokens"
	TableRoles          = "roles"
	TableMenus          = "menus"
	TableApiPermissions = "api_permissions"
	TablePermissions    = "permissions"
)

var UsernameReg = regexp.MustCompile("^[a-zA-Z0-9]([a-zA-Z0-9-_]{1,30}[a-zA-Z0-9])$")

// EnableStatus 启用状态
type EnableStatus string

const (
	EnableStatusEnabled  EnableStatus = "enabled"  // 正常
	EnableStatusDisabled EnableStatus = "disabled" // 禁用
)
