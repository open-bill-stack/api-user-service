package jwt

import "github.com/golang-jwt/jwt/v5"

type ExtendClaims struct {
	UserID string `json:"user_id"`
}

type Claims struct {
	jwt.RegisteredClaims
	ExtendClaims
}
