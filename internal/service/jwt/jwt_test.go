package jwt

import (
	"crypto/ed25519"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWT_CreateAndVerify(t *testing.T) {
	// Генеруємо ключі Ed25519
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	// Ініціалізуємо JWT-сервіс
	j, err := NewJWT(privKey, pubKey)
	require.NoError(t, err)

	// Готуємо дані
	claims := CustomClaims{
		ExtendClaims: ExtendClaims{
			UserID:         "123",
			IsRefreshToken: false,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "test",
		},
	}

	// Створюємо токен
	tokenStr, err := j.Create(claims)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	// Валідуємо токен
	parsedClaims, err := j.Verify(tokenStr)
	require.NoError(t, err)
	require.Equal(t, claims.UserID, parsedClaims.UserID)
	require.Equal(t, claims.Issuer, parsedClaims.Issuer)
	require.Equal(t, claims.IsRefreshToken, parsedClaims.IsRefreshToken)
}

func TestJWT_Verify_InvalidSignature(t *testing.T) {
	// Генеруємо дві пари ключів
	_, privKey1, _ := ed25519.GenerateKey(nil)
	pubKey2, _, _ := ed25519.GenerateKey(nil)

	j1, _ := NewJWT(privKey1, pubKey2) // неправильна перевірка

	// Створюємо валідний токен іншим підписом
	claims := CustomClaims{
		ExtendClaims: ExtendClaims{
			UserID:         "abc",
			IsRefreshToken: true,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "test",
		},
	}

	tokenStr, err := NewJWT(privKey1, privKey1.Public().(ed25519.PublicKey))
	require.NoError(t, err)
	validToken, err := tokenStr.Create(claims)
	require.NoError(t, err)

	// Тепер спробуємо перевірити іншим публічним ключем (повинна бути помилка)
	_, err = j1.Verify(validToken)
	require.Error(t, err)
}

func TestJWT_Verify_InvalidFormat(t *testing.T) {
	// Стандартна пара ключів
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	j, _ := NewJWT(privKey, pubKey)

	// Некоректний токен
	_, err := j.Verify("this_is_not_a_token")
	require.Error(t, err)
}
