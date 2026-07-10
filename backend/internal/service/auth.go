package service

import (
	"errors"
	ewrap "finstat/internal/lib"
	"finstat/internal/repository"
	"finstat/internal/token"

	"golang.org/x/crypto/bcrypt"
)

const TOKEN_LIFE_TIME = 15

type AuthRepo interface {
	InsertUser(username, password string) (uint, error)
	User(username string) (*repository.User, error)
}

type AuthService struct {
	repo      AuthRepo
	jwtSecret []byte
}

func NewAuthService(repo AuthRepo, jwtSecret []byte) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ewrap.Wrap("Couldn't generate hashed password", err)
	}

	_, err = s.repo.InsertUser(username, string(hashedPassword))
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return err
		}
		return ewrap.Wrap("Couldn't insert new user", err)
	}

	return nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.User(username)
	if err != nil {
		return "", ewrap.Wrap("Couldn't get user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", ewrap.Wrap("Passwords mismatched", err)
	}

	newToken, err := token.AddToken(user.ID, s.jwtSecret, TOKEN_LIFE_TIME)
	if err != nil {
		return "", ewrap.Wrap("Couldn't generate new token", err)
	}

	return newToken, nil
}

func (s *AuthService) ID(jwtToken string) (uint, error) {
	id, err := token.ID(jwtToken, s.jwtSecret)
	if err != nil {
		return 0, ewrap.Wrap("Invalid token", err)
	}

	return id, nil
}
