package models

import (
	"database/sql"
	"time"
)

// ChatMessage represents a chat message in our system
type ChatMessage struct {
	ID          int       `json:"id" db:"id"`
	ReportID    int       `json:"report_id" db:"report_id"`
	UserMessage string    `json:"user_message" db:"user_message"`
	AIResponse  string    `json:"ai_response" db:"ai_response"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
}

// ChatMessageRepository defines the interface for chat message database operations
type ChatMessageRepository interface {
	Create(message *ChatMessage) error
	GetByID(id int) (*ChatMessage, error)
	GetByReportID(reportID int, limit, offset int) ([]*ChatMessage, error)
	Update(message *ChatMessage) error
	SoftDelete(id int) error
	HardDelete(id int) error
	GetChatHistory(reportID int) ([]*ChatMessage, error)
}

// SQLChatMessageRepository implements ChatMessageRepository using SQL database
type SQLChatMessageRepository struct {
	db *sql.DB
}

// NewChatMessageRepository creates a new chat message repository
func NewChatMessageRepository(db *sql.DB) ChatMessageRepository {
	return &SQLChatMessageRepository{db: db}
}

// Create inserts a new chat message into the database
func (r *SQLChatMessageRepository) Create(message *ChatMessage) error {
	query := `
		INSERT INTO chat_messages (report_id, user_message, ai_response)
		VALUES (?, ?, ?)
		RETURNING id, created_at`

	// Decision: Auto-generate timestamps and ID, is_deleted defaults to FALSE
	row := r.db.QueryRow(query, message.ReportID, message.UserMessage, message.AIResponse)
	return row.Scan(&message.ID, &message.CreatedAt)
}

// GetByID retrieves a chat message by its ID
func (r *SQLChatMessageRepository) GetByID(id int) (*ChatMessage, error) {
	message := &ChatMessage{}
	query := `
		SELECT id, report_id, user_message, ai_response, created_at, is_deleted
		FROM chat_messages
		WHERE id = ? AND is_deleted = FALSE`

	// Decision: Only return non-deleted messages by default
	row := r.db.QueryRow(query, id)
	err := row.Scan(&message.ID, &message.ReportID, &message.UserMessage,
		&message.AIResponse, &message.CreatedAt, &message.IsDeleted)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetByReportID retrieves chat messages for a specific report with pagination
func (r *SQLChatMessageRepository) GetByReportID(reportID int, limit, offset int) ([]*ChatMessage, error) {
	query := `
		SELECT id, report_id, user_message, ai_response, created_at, is_deleted
		FROM chat_messages
		WHERE report_id = ? AND is_deleted = FALSE
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?`

	// Decision: Order by created_at ASC to show chat history chronologically
	rows, err := r.db.Query(query, reportID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*ChatMessage
	for rows.Next() {
		message := &ChatMessage{}
		err := rows.Scan(&message.ID, &message.ReportID, &message.UserMessage,
			&message.AIResponse, &message.CreatedAt, &message.IsDeleted)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// Update modifies an existing chat message
func (r *SQLChatMessageRepository) Update(message *ChatMessage) error {
	query := `
		UPDATE chat_messages
		SET user_message = ?, ai_response = ?
		WHERE id = ? AND is_deleted = FALSE`

	// Decision: Only allow updating message content, not metadata
	result, err := r.db.Exec(query, message.UserMessage, message.AIResponse, message.ID)
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

// SoftDelete marks a chat message as deleted
func (r *SQLChatMessageRepository) SoftDelete(id int) error {
	query := `UPDATE chat_messages SET is_deleted = TRUE WHERE id = ?`

	// Decision: Soft delete to preserve chat history for analysis
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

// HardDelete permanently removes a chat message
func (r *SQLChatMessageRepository) HardDelete(id int) error {
	query := `DELETE FROM chat_messages WHERE id = ?`

	// Decision: Hard delete for admin cleanup or GDPR compliance
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

// GetChatHistory retrieves all chat messages for a report (for AI context)
func (r *SQLChatMessageRepository) GetChatHistory(reportID int) ([]*ChatMessage, error) {
	query := `
		SELECT id, report_id, user_message, ai_response, created_at, is_deleted
		FROM chat_messages
		WHERE report_id = ? AND is_deleted = FALSE
		ORDER BY created_at ASC`

	// Decision: No pagination for chat history - AI needs full context
	rows, err := r.db.Query(query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*ChatMessage
	for rows.Next() {
		message := &ChatMessage{}
		err := rows.Scan(&message.ID, &message.ReportID, &message.UserMessage,
			&message.AIResponse, &message.CreatedAt, &message.IsDeleted)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}