package jwt

import (
	"crypto/ed25519"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func generateKeys(t *testing.T) (ed25519.PrivateKey, ed25519.PublicKey) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	return priv, pub
}

func TestJWT_CreateAndVerify_Success(t *testing.T) {
	privateKey, publicKey := generateKeys(t)
	j, err := NewJWT(privateKey, publicKey)
	require.NoError(t, err)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test_issuer",
			Subject:   "user_id",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		ExtendClaims: ExtendClaims{
			"custom_value",
		},
	}

	tokenString, err := j.Create(claims)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	parsedClaims, err := j.Verify(tokenString)
	require.NoError(t, err)

	c, ok := parsedClaims.(*Claims)
	require.True(t, ok)
	require.Equal(t, "test_issuer", c.Issuer)
	require.Equal(t, "user_id", c.Subject)
	require.Equal(t, "custom_value", c.UserID)
}

func TestJWT_Verify_InvalidSignature(t *testing.T) {
	privateKey1, publicKey1 := generateKeys(t)
	privateKey2, publicKey2 := generateKeys(t)

	j1, err := NewJWT(privateKey1, publicKey1)
	require.NoError(t, err)

	j2, err := NewJWT(privateKey2, publicKey2)
	require.NoError(t, err)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		ExtendClaims: ExtendClaims{
			"custom_value",
		},
	}

	tokenString, err := j1.Create(claims)
	require.NoError(t, err)

	// Спробуємо верифікувати токен не тим паблік ключем
	_, err = j2.Verify(tokenString)
	require.Error(t, err)
}

func TestJWT_Verify_ExpiredToken(t *testing.T) {
	privateKey, publicKey := generateKeys(t)
	j, err := NewJWT(privateKey, publicKey)
	require.NoError(t, err)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // вже прострочений
		},
		ExtendClaims: ExtendClaims{
			"value",
		},
	}

	tokenString, err := j.Create(claims)
	require.NoError(t, err)

	_, err = j.Verify(tokenString)
	require.Error(t, err)
}
