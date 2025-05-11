package jwt

import (
	"crypto/ed25519"
	"fmt"
	golang "github.com/golang-jwt/jwt/v5"
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

func (s *JWT) Create(data CustomClaims) (string, error) {
	token := golang.NewWithClaims(golang.SigningMethodEdDSA, data)
	return token.SignedString(s.privateKey)
}

func (s *JWT) Verify(tokenString string) (*CustomClaims, error) {
	token, err := golang.ParseWithClaims(tokenString, &CustomClaims{}, func(token *golang.Token) (any, error) {
		if _, ok := token.Method.(*golang.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
