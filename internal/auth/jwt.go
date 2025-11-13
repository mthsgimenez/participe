package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mthsgimenez/participe/internal/env"
)

var key = []byte(env.GetStringFallback("SECRET_KEY", "secret"))

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(email string) (string, error) {
	expirationDate := time.Now().Add(7 * (time.Hour * 24))

	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return key, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
