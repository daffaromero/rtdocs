package service

import (
	"context"
	"errors"
	"rtdocs/model/domain"
	"rtdocs/model/web"
	"rtdocs/repository"
	"rtdocs/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *web.RegisterRequest) (*web.RegisterResponse, error)
	Login(ctx context.Context, req *web.LoginRequest) (*web.LoginResponse, error)
	Logout(ctx context.Context, accessToken string) error
}

type authService struct {
	userRepo repository.UserRepository
	tokenGen utils.TokenGenerator
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *web.RegisterRequest) (*web.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user := &domain.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Role:      "user",
		CreatedAt: time.Now().Local().Format(time.RFC3339),
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenGen.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &web.RegisterResponse{
		UserID:       createdUser.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *web.LoginRequest) (*web.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.userRepo.GetUser(ctx, req.Username)
	if err != nil || user == nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, err := s.tokenGen.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &web.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, accessToken string) error {
	_, err := s.tokenGen.ValidateToken(accessToken)
	if err != nil {
		return err
	}

	return nil
}
