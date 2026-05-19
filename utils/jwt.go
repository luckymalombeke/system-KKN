package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minJWTSecretLen = 32

var weakJWTSecrets = map[string]struct{}{
	"secret": {}, "jwt_secret": {}, "changeme": {}, "password": {},
	"123456": {}, "your-secret-key": {}, "kkn-secret": {},
}

var jwtSecret []byte

// ValidateJWTSecret menolak nilai default/lemah agar token tidak mudah dipalsukan.
func ValidateJWTSecret(secret string) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return errors.New("JWT_SECRET wajib diisi di file .env (min. 32 karakter acak)")
	}
	if len(secret) < minJWTSecretLen {
		return errors.New("JWT_SECRET terlalu pendek (minimal 32 karakter)")
	}
	if _, weak := weakJWTSecrets[strings.ToLower(secret)]; weak {
		return errors.New("JWT_SECRET tidak boleh memakai nilai default; gunakan string acak yang kuat")
	}
	return nil
}

// InitJWTSecret mengatur secret dari environment (dipanggil saat startup).
func InitJWTSecret(secret string) error {
	if err := ValidateJWTSecret(secret); err != nil {
		return err
	}
	jwtSecret = []byte(strings.TrimSpace(secret))
	return nil
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func getJWTSecret() ([]byte, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT belum diinisialisasi, panggil InitJWTSecret terlebih dahulu")
	}
	return jwtSecret, nil
}
