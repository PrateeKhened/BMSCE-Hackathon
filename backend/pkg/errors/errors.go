package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Authentication errors
var (
	ErrInvalidCredentials = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid email or password",
		Type:    "AUTH_ERROR",
	}

	ErrUserNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "User not found",
		Type:    "AUTH_ERROR",
	}

	ErrUserAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "User with this email already exists",
		Type:    "AUTH_ERROR",
	}

	ErrInvalidToken = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid or expired token",
		Type:    "AUTH_ERROR",
	}

	ErrTokenMissing = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Authorization token missing",
		Type:    "AUTH_ERROR",
	}
)

// File upload errors
var (
	ErrFileTooBig = &AppError{
		Code:    http.StatusRequestEntityTooLarge,
		Message: "File size exceeds maximum limit",
		Type:    "UPLOAD_ERROR",
	}

	ErrInvalidFileType = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid file type",
		Type:    "UPLOAD_ERROR",
	}

	ErrFileUploadFailed = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Failed to upload file",
		Type:    "UPLOAD_ERROR",
	}
)

// Database errors
var (
	ErrDatabaseConnection = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Database connection error",
		Type:    "DATABASE_ERROR",
	}

	ErrRecordNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Record not found",
		Type:    "DATABASE_ERROR",
	}
)

// Validation errors
var (
	ErrInvalidInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid input data",
		Type:    "VALIDATION_ERROR",
	}

	ErrMissingRequiredField = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Missing required field",
		Type:    "VALIDATION_ERROR",
	}
)

// NewValidationError creates a new validation error with custom message
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Type:    "VALIDATION_ERROR",
	}
}

// AI processing errors
var (
	ErrAIProcessingFailed = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "AI processing failed",
		Type:    "AI_ERROR",
	}

	ErrReportNotProcessed = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Report has not been processed yet",
		Type:    "AI_ERROR",
	}
)