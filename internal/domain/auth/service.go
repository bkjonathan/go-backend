package auth

import (
	"context"
	"errors"
	"fmt"
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
		return nil, apperror.InvalidInput("invalid register payload", fmt.Errorf("validation: %w", err))
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, apperror.Internal("could not register user", fmt.Errorf("hashing password: %w", err))
	}

	createdUser, err := s.repo.CreateUser(ctx, req.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("email already exists", fmt.Errorf("creating user with duplicate email: %w", err))
		}
		return nil, apperror.Internal("could not register user", fmt.Errorf("creating user: %w", err))
	}

	return s.buildTokenResponse(createdUser.ID, createdUser.Email)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, apperror.InvalidInput("invalid login payload", fmt.Errorf("validating login payload: %w", err))
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("login failed", zap.String("email", req.Email), zap.Error(fmt.Errorf("getting user by email: %w", err)))
		return nil, apperror.Unauthorized("invalid credentials", fmt.Errorf("finding login user: %w", err))
	}

	if err := password.Compare(user.PasswordHash, req.Password); err != nil {
		s.logger.Warn("login failed due to password mismatch", zap.String("email", req.Email))
		return nil, apperror.Unauthorized("invalid credentials", fmt.Errorf("checking password: %w", err))
	}

	return s.buildTokenResponse(user.ID, user.Email)
}

// buildTokenResponse generates a JWT and wraps it in a TokenResponse.
// Extracted to avoid duplication between Register and Login.
func (s *service) buildTokenResponse(userID int64, email string) (*TokenResponse, error) {
	token, err := s.tokenManager.Generate(userID, email)
	if err != nil {
		return nil, apperror.Internal("could not issue token", fmt.Errorf("generating token: %w", err))
	}

	return &TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.tokenManager.TTL().Seconds()),
		User: UserResponse{
			ID:    userID,
			Email: email,
		},
	}, nil
}
