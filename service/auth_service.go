package service

import (
	"context"
	"rtdocs/model/domain"
	"rtdocs/model/web"
	"rtdocs/repository"
	"rtdocs/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	Login(ctx context.Context, req web.LoginRequest) (*web.LoginResponse, error)
	Logout(ctx context.Context, accessToken string) error
}

type authService struct {
	userRepo repository.UserRepository
	tokenGen utils.TokenGenerator
}

func NewAuthService(userRepo repository.UserRepository, tokenGen utils.TokenGenerator) AuthService {
	return &authService{
		userRepo: userRepo,
		tokenGen: tokenGen,
	}
}

func (s *authService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.ID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *authService) Login(ctx context.Context, req web.LoginRequest) (*web.LoginResponse, error) {
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
