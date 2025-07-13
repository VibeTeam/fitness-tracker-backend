package use_case

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/VibeTeam/fitness-tracker-backend/user/auth"
	"github.com/VibeTeam/fitness-tracker-backend/user/models"
	"github.com/VibeTeam/fitness-tracker-backend/user/repository"
)

// AuthService provides high-level authentication workflows (register/login/token handling).
type AuthService struct {
	repo        repository.UserRepository
	tokenManger *auth.Manager
}

// NewAuthService wires repository and token manager into a ready-to-use AuthService.
func NewAuthService(repo repository.UserRepository, tokenMgr *auth.Manager) *AuthService {
	return &AuthService{repo: repo, tokenManger: tokenMgr}
}

var (
	// ErrEmailAlreadyUsed is returned when attempting to register with an existing e-mail.
	ErrEmailAlreadyUsed = errors.New("email is already taken")
	// ErrInvalidCredentials is returned when login credentials do not match.
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// Register creates a new user and immediately returns freshly minted JWT pair.
func (s *AuthService) Register(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return "", "", ErrEmailAlreadyUsed
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = s.tokenManger.NewTokens(int32(user.ID))
	return
}

// Login verifies the supplied credentials and returns a new JWT pair upon success.
func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err = s.tokenManger.NewTokens(int32(user.ID))
	return
}

// Refresh validates the provided refresh token and issues a fresh access/refresh pair.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// Business logic resides in token manager; just proxy.
	return s.tokenManger.RefreshTokens(refreshToken)
}

// Validate parses the access token and returns the user ID if it is valid.
func (s *AuthService) Validate(ctx context.Context, accessToken string) (userID uint, err error) {
	id, err := s.tokenManger.ValidateAccessToken(accessToken)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
