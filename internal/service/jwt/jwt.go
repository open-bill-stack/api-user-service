package jwt

import (
	"crypto/ed25519"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewJWT(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (*JWT, error) {
	return &JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (s *JWT) Create(data Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, data)
	return token.SignedString(s.privateKey)
}

func (s *JWT) Verify(tokenString string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
