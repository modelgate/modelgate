package model

type LoginRequest struct {
	Username   string
	Password   string
	RememberMe bool
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenRequest struct {
	RefreshToken string
}
