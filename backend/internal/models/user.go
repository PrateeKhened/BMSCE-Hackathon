package models

import (
	"database/sql"
	"time"
)

// User represents a user in our system
// Decision: Using struct tags for both JSON and database mapping
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

// UserRepository defines the interface for user database operations
// Decision: Using repository pattern for better testability and separation of concerns
type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id int) error
	List(limit, offset int) ([]*User, error)
}

// SQLUserRepository implements UserRepository using SQL database
type SQLUserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &SQLUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *SQLUserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, email_verified, is_active)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at`

	// Decision: Using RETURNING clause to get generated ID and timestamps
	row := r.db.QueryRow(query, user.Email, user.PasswordHash, user.FullName, user.EmailVerified, user.IsActive)
	return row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by their ID
func (r *SQLUserRepository) GetByID(id int) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, full_name, email_verified, is_active, created_at, updated_at
		FROM users
		WHERE id = ? AND is_active = TRUE`

	// Decision: Only return active users in standard queries
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.EmailVerified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil for not found, not an error
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by their email address
func (r *SQLUserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, full_name, email_verified, is_active, created_at, updated_at
		FROM users
		WHERE email = ? AND is_active = TRUE`

	row := r.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.EmailVerified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update modifies an existing user
func (r *SQLUserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET email = ?, full_name = ?, email_verified = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND is_active = TRUE`

	// Decision: Not allowing password updates here - separate method for security
	result, err := r.db.Exec(query, user.Email, user.FullName, user.EmailVerified, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // User not found or not active
	}

	return nil
}

// Delete soft deletes a user (sets is_active to FALSE)
func (r *SQLUserRepository) Delete(id int) error {
	query := `UPDATE users SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	// Decision: Soft delete to preserve data integrity with reports and chat history
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// List retrieves a paginated list of users
func (r *SQLUserRepository) List(limit, offset int) ([]*User, error) {
	query := `
		SELECT id, email, password_hash, full_name, email_verified, is_active, created_at, updated_at
		FROM users
		WHERE is_active = TRUE
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
			&user.EmailVerified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}