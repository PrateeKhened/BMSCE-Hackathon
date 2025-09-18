package models

import (
	"database/sql"
	"time"
)

// Report represents a medical report in our system
type Report struct {
	ID                int        `json:"id" db:"id"`
	UserID           int        `json:"user_id" db:"user_id"`
	OriginalFilename string     `json:"original_filename" db:"original_filename"`
	FilePath         string     `json:"file_path" db:"file_path"`
	FileType         string     `json:"file_type" db:"file_type"`
	FileSize         int64      `json:"file_size" db:"file_size"`
	SimplifiedSummary string    `json:"simplified_summary" db:"simplified_summary"`
	ProcessingStatus string     `json:"processing_status" db:"processing_status"`
	UploadDate       time.Time  `json:"upload_date" db:"upload_date"`
	ProcessedAt      *time.Time `json:"processed_at" db:"processed_at"` // Nullable
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// ReportRepository defines the interface for report database operations
type ReportRepository interface {
	Create(report *Report) error
	GetByID(id int) (*Report, error)
	GetByUserID(userID int, limit, offset int) ([]*Report, error)
	Update(report *Report) error
	UpdateProcessingStatus(id int, status string, summary string) error
	Delete(id int) error
	GetPendingReports(limit int) ([]*Report, error)
}

// SQLReportRepository implements ReportRepository using SQL database
type SQLReportRepository struct {
	db *sql.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *sql.DB) ReportRepository {
	return &SQLReportRepository{db: db}
}

// Create inserts a new report into the database
func (r *SQLReportRepository) Create(report *Report) error {
	query := `
		INSERT INTO reports (user_id, original_filename, file_path, file_type, file_size, processing_status)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id, upload_date, created_at, updated_at`

	// Decision: Set processing_status to 'pending' by default, timestamps auto-generated
	row := r.db.QueryRow(query, report.UserID, report.OriginalFilename,
		report.FilePath, report.FileType, report.FileSize, "pending")

	return row.Scan(&report.ID, &report.UploadDate, &report.CreatedAt, &report.UpdatedAt)
}

// GetByID retrieves a report by its ID
func (r *SQLReportRepository) GetByID(id int) (*Report, error) {
	report := &Report{}
	query := `
		SELECT id, user_id, original_filename, file_path, file_type, file_size,
			   simplified_summary, processing_status, upload_date, processed_at,
			   created_at, updated_at
		FROM reports
		WHERE id = ?`

	row := r.db.QueryRow(query, id)
	err := row.Scan(&report.ID, &report.UserID, &report.OriginalFilename,
		&report.FilePath, &report.FileType, &report.FileSize,
		&report.SimplifiedSummary, &report.ProcessingStatus, &report.UploadDate,
		&report.ProcessedAt, &report.CreatedAt, &report.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetByUserID retrieves reports for a specific user with pagination
func (r *SQLReportRepository) GetByUserID(userID int, limit, offset int) ([]*Report, error) {
	query := `
		SELECT id, user_id, original_filename, file_path, file_type, file_size,
			   simplified_summary, processing_status, upload_date, processed_at,
			   created_at, updated_at
		FROM reports
		WHERE user_id = ?
		ORDER BY upload_date DESC
		LIMIT ? OFFSET ?`

	// Decision: Order by upload_date DESC to show newest reports first
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*Report
	for rows.Next() {
		report := &Report{}
		err := rows.Scan(&report.ID, &report.UserID, &report.OriginalFilename,
			&report.FilePath, &report.FileType, &report.FileSize,
			&report.SimplifiedSummary, &report.ProcessingStatus, &report.UploadDate,
			&report.ProcessedAt, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

// Update modifies an existing report
func (r *SQLReportRepository) Update(report *Report) error {
	query := `
		UPDATE reports
		SET original_filename = ?, file_type = ?, file_size = ?,
			simplified_summary = ?, processing_status = ?, processed_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	result, err := r.db.Exec(query, report.OriginalFilename, report.FileType,
		report.FileSize, report.SimplifiedSummary, report.ProcessingStatus,
		report.ProcessedAt, report.ID)
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

// UpdateProcessingStatus updates the processing status and summary
// Decision: Separate method for AI processing updates to avoid race conditions
func (r *SQLReportRepository) UpdateProcessingStatus(id int, status string, summary string) error {
	query := `
		UPDATE reports
		SET processing_status = ?, simplified_summary = ?,
			processed_at = CASE WHEN ? = 'completed' THEN CURRENT_TIMESTAMP ELSE processed_at END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	// Decision: Set processed_at only when status is 'completed'
	result, err := r.db.Exec(query, status, summary, status, id)
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

// Delete removes a report from the database
func (r *SQLReportRepository) Delete(id int) error {
	query := `DELETE FROM reports WHERE id = ?`

	// Decision: Hard delete for reports since they're user-generated content
	// Chat messages will be cascade deleted due to foreign key constraint
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

// GetPendingReports retrieves reports that need AI processing
func (r *SQLReportRepository) GetPendingReports(limit int) ([]*Report, error) {
	query := `
		SELECT id, user_id, original_filename, file_path, file_type, file_size,
			   simplified_summary, processing_status, upload_date, processed_at,
			   created_at, updated_at
		FROM reports
		WHERE processing_status = 'pending'
		ORDER BY upload_date ASC
		LIMIT ?`

	// Decision: Process oldest pending reports first (FIFO)
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*Report
	for rows.Next() {
		report := &Report{}
		err := rows.Scan(&report.ID, &report.UserID, &report.OriginalFilename,
			&report.FilePath, &report.FileType, &report.FileSize,
			&report.SimplifiedSummary, &report.ProcessingStatus, &report.UploadDate,
			&report.ProcessedAt, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}