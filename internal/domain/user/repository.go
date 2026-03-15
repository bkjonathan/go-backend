package user

import (
	"context"
	"fmt"
	"thomas-backend/internal/database"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, limit, offset int32) ([]User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, id int64, email, passwordHash string) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, limit, offset int32) ([]User, error) {
	var rows []database.User
	err := r.db.WithContext(ctx).
		Order("id ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&rows).
		Error
	if err != nil {
		return nil, fmt.Errorf("listing users from repository: %w", err)
	}

	users := make([]User, 0, len(rows))
	for _, row := range rows {
		users = append(users, mapDBUser(row))
	}

	return users, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var row database.User
	err := r.db.WithContext(ctx).First(&row, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting user by id from repository: %w", err)
	}

	mapped := mapDBUser(row)
	return &mapped, nil
}

func (r *repository) Update(ctx context.Context, id int64, email, passwordHash string) (*User, error) {
	var updated database.User
	if err := r.db.WithContext(ctx).First(&updated, id).Error; err != nil {
		return nil, fmt.Errorf("updating user in repository: %w", err)
	}

	updated.Email = email
	updated.PasswordHash = passwordHash

	err := r.db.WithContext(ctx).Save(&updated).Error
	if err != nil {
		return nil, fmt.Errorf("updating user in repository: %w", err)
	}

	mapped := mapDBUser(updated)
	return &mapped, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&database.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("deleting user from repository: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user %d not found for delete: %w", id, gorm.ErrRecordNotFound)
	}
	return nil
}

func mapDBUser(u database.User) User {
	return User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
