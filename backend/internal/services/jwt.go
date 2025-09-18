package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents our custom JWT claims
// Decision: Embed jwt.RegisteredClaims for standard fields (exp, iat, etc.)
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secret     []byte        // Secret key for signing tokens
	expiration time.Duration // Token expiration time
}

// NewJWTService creates a new JWT service
// Decision: Accept secret and expiration as parameters for configuration flexibility
func NewJWTService(secret string, expiration time.Duration) *JWTService {
	return &JWTService{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

// GenerateToken creates a new JWT token for a user
// Decision: Accept userID and email as separate params for type safety
func (js *JWTService) GenerateToken(userID int, email string) (string, error) {
	// Decision: Set token expiration from current time + configured duration
	expirationTime := time.Now().Add(js.expiration)

	// Decision: Create custom claims with user information
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "medical-report-backend", // Decision: Identify our service
		},
	}

	// Decision: Use HS256 signing method (HMAC with SHA-256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Decision: Sign the token with our secret key
	tokenString, err := token.SignedString(js.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT token
// Decision: Return claims if valid, error if invalid/expired
func (js *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Decision: Parse token with custom claims struct
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Decision: Verify the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return js.secret, nil
	})

	if err != nil {
		return nil, err
	}

	// Decision: Extract claims if token is valid
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token for existing valid token
// Decision: Useful for extending user sessions without re-authentication
func (js *JWTService) RefreshToken(tokenString string) (string, error) {
	// Decision: First validate the existing token
	claims, err := js.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Decision: Generate new token with same user information
	return js.GenerateToken(claims.UserID, claims.Email)
}

// GetUserFromToken extracts user information from a token
// Decision: Convenience method for middleware to get user context
func (js *JWTService) GetUserFromToken(tokenString string) (userID int, email string, err error) {
	claims, err := js.ValidateToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	return claims.UserID, claims.Email, nil
}