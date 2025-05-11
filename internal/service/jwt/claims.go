package jwt

import golang "github.com/golang-jwt/jwt/v5"

type ExtendClaims struct {
	UserID         string `json:"user_id"`
	IsRefreshToken bool   `json:"is_refresh_token"`
}

type CustomClaims struct {
	golang.RegisteredClaims
	ExtendClaims
}
