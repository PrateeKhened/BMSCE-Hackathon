package services

import (
	"strings"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/errors"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/types"
)

// AuthService handles authentication business logic
// Decision: Use dependency injection for testability and flexibility
type AuthService struct {
	userRepo        models.UserRepository
	passwordService *PasswordService
	jwtService      *JWTService
}

// NewAuthService creates a new authentication service
// Decision: Inject all dependencies to allow for mocking in tests
func NewAuthService(
	userRepo models.UserRepository,
	passwordService *PasswordService,
	jwtService *JWTService,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// SignUp creates a new user account
// Decision: Accept signup request struct for validation and type safety
func (as *AuthService) SignUp(req *types.SignupRequest) (*types.LoginResponse, error) {
	// Decision: Validate email format before processing
	if !isValidEmail(req.Email) {
		return nil, errors.ErrInvalidInput
	}

	// Decision: Check minimum password length
	if len(req.Password) < 6 {
		return nil, errors.ErrInvalidInput
	}

	// Decision: Normalize email to lowercase for consistency
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Decision: Check if user already exists before processing
	existingUser, err := as.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.ErrDatabaseConnection
	}

	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// Decision: Hash password before storing
	hashedPassword, err := as.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, errors.ErrDatabaseConnection
	}

	// Decision: Create user with email verification disabled by default
	user := &models.User{
		Email:         email,
		PasswordHash:  hashedPassword,
		FullName:      strings.TrimSpace(req.FullName),
		EmailVerified: false, // Decision: Require email verification in future
		IsActive:      true,
	}

	// Decision: Create user in database
	err = as.userRepo.Create(user)
	if err != nil {
		return nil, errors.ErrDatabaseConnection
	}

	// Decision: Generate JWT token immediately after successful signup
	token, err := as.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.ErrDatabaseConnection
	}

	// Decision: Return user data and token for immediate login
	response := &types.LoginResponse{
		Token: token,
		User:  convertModelUserToTypeUser(user),
	}

	return response, nil
}

// Login authenticates a user and returns a JWT token
// Decision: Accept login request struct for validation
func (as *AuthService) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	// Decision: Validate input before processing
	if !isValidEmail(req.Email) || len(req.Password) == 0 {
		return nil, errors.ErrInvalidInput
	}

	// Decision: Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Decision: Get user from database
	user, err := as.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.ErrDatabaseConnection
	}

	// Decision: Return same error for both "user not found" and "wrong password"
	// This prevents user enumeration attacks
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Decision: Verify password using constant-time comparison
	if !as.passwordService.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.ErrInvalidCredentials
	}

	// Decision: Generate fresh JWT token on each login
	token, err := as.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.ErrDatabaseConnection
	}

	// Decision: Return user data and token
	response := &types.LoginResponse{
		Token: token,
		User:  convertModelUserToTypeUser(user),
	}

	return response, nil
}

// GetUserFromToken validates a JWT token and returns user information
// Decision: Useful for middleware to authenticate requests
func (as *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	// Decision: Validate token first
	userID, email, err := as.jwtService.GetUserFromToken(tokenString)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Decision: Get fresh user data from database (handles user deactivation)
	user, err := as.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.ErrDatabaseConnection
	}

	// Decision: Return error if user not found or deactivated
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// Decision: Verify email matches token (prevents token reuse after email change)
	if user.Email != email {
		return nil, errors.ErrInvalidToken
	}

	return user, nil
}

// RefreshToken generates a new token for valid existing token
// Decision: Extend user sessions without requiring re-authentication
func (as *AuthService) RefreshToken(tokenString string) (string, error) {
	// Decision: Validate current token and get user info
	_, err := as.GetUserFromToken(tokenString)
	if err != nil {
		return "", err
	}

	// Decision: Generate new token using JWT service
	newToken, err := as.jwtService.RefreshToken(tokenString)
	if err != nil {
		return "", errors.ErrInvalidToken
	}

	return newToken, nil
}

// isValidEmail performs basic email validation
// Decision: Simple validation for now, can be enhanced with regex if needed
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return len(email) > 0 && strings.Contains(email, "@") && strings.Contains(email, ".")
}

// convertModelUserToTypeUser converts models.User to types.User
// Decision: Keep models and API types separate for better abstraction
func convertModelUserToTypeUser(user *models.User) types.User {
	return types.User{
		ID:            user.ID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		FullName:      user.FullName,
		EmailVerified: user.EmailVerified,
		IsActive:      user.IsActive,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}