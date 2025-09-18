package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type UploadConfig struct {
	MaxFileSize int64
	UploadPath  string
	AllowedTypes []string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "localhost"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite3"),
			DSN:    getEnv("DB_DSN", "./medical_reports.db"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		},
		Upload: UploadConfig{
			MaxFileSize: getInt64Env("MAX_FILE_SIZE", 10*1024*1024), // 10MB default
			UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
			AllowedTypes: []string{"application/pdf", "text/plain", "image/jpeg", "image/png"},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}