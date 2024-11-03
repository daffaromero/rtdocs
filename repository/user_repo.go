package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rtdocs/model/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (q *userRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		log.Println("User ID is required")
		return nil, nil
	}
	query := "SELECT * FROM users WHERE id = $1"

	var user domain.User
	row := q.db.QueryRow(ctx, query, id)

	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

func (q *userRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := "SELECT * FROM users"
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (q *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := "INSERT INTO users (id, username, password) VALUES ($1, $2, $3) RETURNING id"
	row := q.db.QueryRow(ctx, query, user.ID, user.Username, user.Password)

	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (q *userRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := "UPDATE users SET username = $2, password = $3 WHERE id = $1 RETURNING id, username, password"
	row := q.db.QueryRow(ctx, query, user.ID, user.Username, user.Password)

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return user, nil
}
