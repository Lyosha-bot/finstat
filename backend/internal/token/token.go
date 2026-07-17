package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TOKEN_ISSUER = "finstat"
)

type CustomClaims struct {
	UserID uint   `json:"user_id,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	jwt.RegisteredClaims
}

func NewAccessToken(user_id uint, jwt_secret []byte, lifetime uint) (string, error) {
	claims := CustomClaims{
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TOKEN_ISSUER,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(lifetime))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwt_secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(uuid string, jwt_secret []byte, lifetime uint) (string, error) {
	claims := CustomClaims{
		UUID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TOKEN_ISSUER,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(lifetime))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwt_secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Claims(jwt_token string, jwt_secret []byte) (*CustomClaims, error) {
	var result CustomClaims
	token, err := jwt.ParseWithClaims(jwt_token, &result, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwt_secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Token is dead")
	}

	log.Println(result.UserID)
	log.Println(result.UUID)

	return &result, nil
}
