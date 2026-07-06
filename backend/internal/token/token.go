// TODO: Настроить context
// TODO: Добавить интерфейсы

package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TOKEN_ISSUER    = "auth.my-financials"
	TOKEN_LIFE_TIME = 20 // В минутах
)

type customClaims struct {
	ID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewToken(user_id uint, jwt_secret []byte) (string, error) {
	claims := customClaims{
		ID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TOKEN_ISSUER,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * TOKEN_LIFE_TIME)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwt_secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
