package auth

import (
	"context"
	"fmt"
	"thomas-backend/internal/database"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, email, passwordHash string) (*User, error) {
	created := database.User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	err := r.db.WithContext(ctx).Create(&created).Error
	if err != nil {
		return nil, fmt.Errorf("creating user in repository: %w", err)
	}

	return &User{
		ID:           created.ID,
		Email:        created.Email,
		PasswordHash: created.PasswordHash,
	}, nil

}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var result database.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&result).
		Error
	if err != nil {
		return nil, fmt.Errorf("getting user by email from repository: %w", err)
	}

	return &User{
		ID:           result.ID,
		Email:        result.Email,
		PasswordHash: result.PasswordHash,
	}, nil
}
