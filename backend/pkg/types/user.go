package types

import "time"

type User struct {
	ID            int       `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"` // Never expose password in JSON
	FullName      string    `json:"full_name" db:"full_name"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required,min=2"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}