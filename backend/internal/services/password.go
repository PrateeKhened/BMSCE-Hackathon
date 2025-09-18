package services

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and verification
// Decision: Using a service struct to allow for future configuration
type PasswordService struct {
	cost int // bcrypt cost factor
}

// NewPasswordService creates a new password service
// Decision: DefaultCost (10) provides good security/performance balance
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost, // Cost of 10
	}
}

// NewPasswordServiceWithCost creates a password service with custom cost
// Decision: Allow custom cost for testing (lower) or high-security (higher)
func NewPasswordServiceWithCost(cost int) *PasswordService {
	return &PasswordService{
		cost: cost,
	}
}

// HashPassword hashes a plain text password using bcrypt
// Decision: Return error to handle bcrypt failures (memory allocation, etc.)
func (ps *PasswordService) HashPassword(password string) (string, error) {
	// Decision: Convert to bytes as bcrypt works with []byte
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), ps.cost)
	if err != nil {
		return "", err
	}

	// Decision: Return as string for easier database storage
	return string(bytes), nil
}

// CheckPassword verifies a password against its hash
// Decision: Return bool for simple usage, log errors internally if needed
func (ps *PasswordService) CheckPassword(password, hash string) bool {
	// Decision: bcrypt.CompareHashAndPassword handles all the comparison logic
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	// Decision: Return false for any error (invalid hash, wrong password, etc.)
	return err == nil
}

// GetCost returns the current cost factor
// Decision: Useful for debugging and ensuring correct configuration
func (ps *PasswordService) GetCost() int {
	return ps.cost
}