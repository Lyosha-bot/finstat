package service

import (
	"finstat/internal/apperr"
	"finstat/internal/lib"
	"finstat/internal/models"
	"finstat/internal/repository"
	"finstat/internal/token"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const ACCESS_TOKEN_LIFE_TIME = 15 * 60
const REFRESH_TOKEN_LIFE_TIME = 7 * 24 * 60 * 60

type User models.User
type RefreshToken models.RefreshToken

type AuthService struct {
	repo             repository.Auth
	jwtAccessSecret  []byte
	jwtRefreshSecret []byte
}

func NewAuthService(repo repository.Auth, jwtAccessSecret, jwtRefreshSecret []byte) *AuthService {
	return &AuthService{
		repo:             repo,
		jwtAccessSecret:  jwtAccessSecret,
		jwtRefreshSecret: jwtRefreshSecret,
	}
}

func (s *AuthService) generateTokens(userID uint) (accessToken string, refreshToken string, err error) {
	uuid, err := s.repo.InsertRefreshToken(userID, time.Now().Add(time.Second*REFRESH_TOKEN_LIFE_TIME).UTC())
	if err != nil {
		return "", "", err
	}

	refreshToken, err = token.NewRefreshToken(uuid, s.jwtRefreshSecret, REFRESH_TOKEN_LIFE_TIME)
	if err != nil {
		return "", "", lib.Ewrap("Couldn't generate new refresh token", err)
	}

	accessToken, err = token.NewAccessToken(userID, s.jwtAccessSecret, ACCESS_TOKEN_LIFE_TIME)
	if err != nil {
		return "", "", lib.Ewrap("Couldn't generate new access token", err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return lib.Ewrap("Couldn't generate hashed password", err)
	}

	if err := s.repo.InsertUser(username, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(username, password string) (accessToken string, refreshToken string, err error) {
	user, err := s.repo.User(username)
	if err != nil {
		return "", "", lib.Ewrap("Couldn't get user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", "", apperr.PasswordMismatched
	}

	return s.generateTokens(user.ID)
}

func (s *AuthService) Refresh(refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	claims, err := token.Claims(refreshToken, s.jwtRefreshSecret)
	if err != nil {
		return "", "", err
	}

	token, err := s.repo.RefreshToken(claims.UUID)
	if err != nil {
		return "", "", err
	}

	_, err = s.repo.DeleteRefreshToken(token.UUID)
	if err != nil {
		return "", "", err
	}

	if token.ExpiresAt.Before(time.Now().UTC()) {
		return "", "", apperr.TokenExpired
	}

	return s.generateTokens(token.UserID)
}

func (s *AuthService) Logout(refreshToken string) error {
	claims, err := token.Claims(refreshToken, s.jwtRefreshSecret)
	if err != nil {
		return lib.Ewrap("Couldn't get claims to logout", err)
	}

	_, err = s.repo.DeleteRefreshToken(claims.UUID)

	return err
}

func (s *AuthService) ID(jwtAccessToken string) (uint, error) {
	claims, err := token.Claims(jwtAccessToken, s.jwtAccessSecret)
	if err != nil {
		return 0, lib.Ewrap("Invalid token", err)
	}

	return claims.UserID, nil
}
