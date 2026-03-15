package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"thomas-backend/internal/apperror"
	"thomas-backend/pkg/jwtutil"
	"thomas-backend/pkg/password"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*TokenResponse, error)
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
}

type service struct {
	repo         Repository
	validate     *validator.Validate
	tokenManager *jwtutil.Manager
	logger       *zap.Logger
}

func NewService(repo Repository, validate *validator.Validate, tokenManager *jwtutil.Manager, logger *zap.Logger) Service {
	return &service{
		repo:         repo,
		validate:     validate,
		tokenManager: tokenManager,
		logger:       logger,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*TokenResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, apperror.New(
			apperror.CodeInvalidInput,
			"invalid register payload",
			http.StatusBadRequest,
			fmt.Errorf("validation: %w", err),
		)
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not register user",
			http.StatusInternalServerError,
			fmt.Errorf("hashing password: %w", err),
		)
	}

	createdUser, err := s.repo.CreateUser(ctx, req.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.New(
				apperror.CodeConflict,
				"email already exists",
				http.StatusConflict,
				fmt.Errorf("creating user with duplicate email: %w", err),
			)
		}
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not register user",
			http.StatusInternalServerError,
			fmt.Errorf("creating user: %w", err),
		)
	}

	token, err := s.tokenManager.Generate(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not issue token",
			http.StatusInternalServerError,
			fmt.Errorf("generating token: %w", err),
		)
	}

	return &TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.tokenManager.TTL().Seconds()),
		User: UserResponse{
			ID:    createdUser.ID,
			Email: createdUser.Email,
		},
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, apperror.New(
			apperror.CodeInvalidInput,
			"invalid login payload",
			http.StatusBadRequest,
			fmt.Errorf("validating login payload: %w", err),
		)
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("login failed", zap.String("email", req.Email), zap.Error(fmt.Errorf("getting user by email: %w", err)))
		return nil, apperror.New(
			apperror.CodeUnauthorized,
			"invalid credentials",
			http.StatusUnauthorized,
			fmt.Errorf("finding login user: %w", err),
		)
	}

	if err := password.Compare(user.PasswordHash, req.Password); err != nil {
		s.logger.Warn("login failed due to password mismatch", zap.String("email", req.Email))
		return nil, apperror.New(
			apperror.CodeUnauthorized,
			"invalid credentials",
			http.StatusUnauthorized,
			fmt.Errorf("checking password: %w", err),
		)
	}

	token, err := s.tokenManager.Generate(user.ID, user.Email)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeInternal,
			"could not issue token",
			http.StatusInternalServerError,
			fmt.Errorf("generating token: %w", err),
		)
	}

	return &TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.tokenManager.TTL().Seconds()),
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	}, nil
}
