package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"thomas-backend/internal/apperror"
	"thomas-backend/pkg/password"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	List(ctx context.Context, limit, offset int32) ([]UserResponse, error)
	GetByID(ctx context.Context, id int64) (*UserResponse, error)
	GetMe(ctx context.Context, id int64) (*UserResponse, error)
	Update(ctx context.Context, id int64, req UpdateRequest) (*UserResponse, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo     Repository
	validate *validator.Validate
	logger   *zap.Logger
}

func NewService(repo Repository, validate *validator.Validate, logger *zap.Logger) Service {
	return &service{
		repo:     repo,
		validate: validate,
		logger:   logger,
	}
}
func (s *service) List(ctx context.Context, limit, offset int32) ([]UserResponse, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not list users",
			http.StatusInternalServerError,
			fmt.Errorf("listing users in service: %w", err),
		)
	}

	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeNotFound,
				"user not found",
				http.StatusNotFound,
				fmt.Errorf("getting user by id: %w", err),
			)
		}
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not get user",
			http.StatusInternalServerError,
			fmt.Errorf("getting user by id in service: %w", err),
		)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *service) GetMe(ctx context.Context, id int64) (*UserResponse, error) {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}
	return user, nil
}

func (s *service) Update(ctx context.Context, id int64, req UpdateRequest) (*UserResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, apperror.New(
			apperror.CodeInvalidInput,
			"invalid update payload",
			http.StatusBadRequest,
			fmt.Errorf("validating update payload: %w", err),
		)
	}

	if req.Email == nil && req.Password == nil {
		return nil, apperror.New(
			apperror.CodeInvalidInput,
			"at least one field is required",
			http.StatusBadRequest,
			fmt.Errorf("update payload missing fields"),
		)
	}

	currentUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeNotFound,
				"user not found",
				http.StatusNotFound,
				fmt.Errorf("getting user before update: %w", err),
			)
		}
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not update user",
			http.StatusInternalServerError,
			fmt.Errorf("getting current user before update: %w", err),
		)
	}

	email := currentUser.Email
	if req.Email != nil {
		email = *req.Email
	}

	passwordHash := currentUser.PasswordHash
	if req.Password != nil {
		hashed, hashErr := password.Hash(*req.Password)
		if hashErr != nil {
			return nil, apperror.New(
				apperror.CodeInternal,
				"could not update user",
				http.StatusInternalServerError,
				fmt.Errorf("hashing updated password: %w", hashErr),
			)
		}
		passwordHash = hashed
	}

	updated, err := s.repo.Update(ctx, id, email, passwordHash)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.New(
				apperror.CodeConflict,
				"email already exists",
				http.StatusConflict,
				fmt.Errorf("updating user with duplicate email: %w", err),
			)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeNotFound,
				"user not found",
				http.StatusNotFound,
				fmt.Errorf("updating missing user: %w", err),
			)
		}
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not update user",
			http.StatusInternalServerError,
			fmt.Errorf("updating user in service: %w", err),
		)
	}

	response := updated.ToResponse()
	return &response, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(
				apperror.CodeNotFound,
				"user not found",
				http.StatusNotFound,
				fmt.Errorf("deleting missing user: %w", err),
			)
		}
		s.logger.Error("delete user failed", zap.Int64("id", id), zap.Error(fmt.Errorf("deleting user in service: %w", err)))
		return apperror.New(
			apperror.CodeInternal,
			"could not delete user",
			http.StatusInternalServerError,
			fmt.Errorf("deleting user in service: %w", err),
		)
	}
	return nil
}
