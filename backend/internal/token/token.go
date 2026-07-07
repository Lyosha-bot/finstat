// TODO: Настроить context
// TODO: Добавить интерфейсы

package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TOKEN_ISSUER    = "finstat"
	TOKEN_LIFE_TIME = 20 // В минутах
)

type customClaims struct {
	ID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewToken(user_id uint, jwt_secret []byte, lifetime uint) (string, error) {
	claims := customClaims{
		ID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TOKEN_ISSUER,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(lifetime))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwt_secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ID(jwt_token string, jwt_secret []byte) (uint, error) {
	var result customClaims
	token, err := jwt.ParseWithClaims(jwt_token, &result, func(token *jwt.Token) (any, error) {
		return jwt_secret, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("Token is dead")
	}

	return result.ID, nil
}
