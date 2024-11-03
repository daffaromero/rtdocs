package service

import (
	"context"
	"rtdocs/model/domain"
	"rtdocs/repository"

	"github.com/google/uuid"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, newUser *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, updatedUser *domain.User) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *userService) CreateUser(ctx context.Context, newUser *domain.User) (*domain.User, error) {
	newUser.ID = uuid.New().String()
	return s.repo.CreateUser(ctx, newUser)
}

func (s *userService) UpdateUser(ctx context.Context, updatedUser *domain.User) (*domain.User, error) {
	return s.repo.UpdateUser(ctx, updatedUser)
}
