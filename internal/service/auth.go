package service

import (
	"fmt"
	"github.com/PaulYakow/gophermart/config"
	"github.com/PaulYakow/gophermart/internal/repo"
	"github.com/PaulYakow/gophermart/internal/util"
	"github.com/PaulYakow/gophermart/internal/util/token"
)

type AuthService struct {
	cfg        config.AuthConfig
	repo       repo.IAuthorization
	tokenMaker token.IMaker
}

func NewAuthService(repo repo.IAuthorization) (*AuthService, error) {
	cfg, err := config.LoadAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot create config: %w", err)
	}

	tokenMaker, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	return &AuthService{
		cfg:        cfg,
		repo:       repo,
		tokenMaker: tokenMaker,
	}, nil
}

func (s *AuthService) CreateUser(login, password string) (int, error) {
	passwordHash, err := util.HashPassword(password)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateUser(login, passwordHash)
}

func (s *AuthService) GetUser(login, password string) (int, error) {
	user, err := s.repo.GetUser(login)
	if err != nil {
		return 0, ErrLoginNotExist
	}

	err = util.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return 0, ErrMismatchPassword
	}

	return user.ID, nil
}

func (s *AuthService) GenerateToken(userID int) (string, error) {
	authToken, err := s.tokenMaker.CreateToken(userID, s.cfg.AccessTokenDuration)
	if err != nil {
		return "", err
	}

	return authToken, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return 0, err
	}

	return payload.UserID, nil
}
