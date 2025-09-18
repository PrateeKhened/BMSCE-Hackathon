-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    original_filename TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    simplified_summary TEXT,
    processing_status TEXT DEFAULT 'pending' CHECK (processing_status IN ('pending', 'processing', 'completed', 'failed')),
    upload_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on user_id for faster user report queries
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);

-- Create index on upload_date for chronological sorting
CREATE INDEX IF NOT EXISTS idx_reports_upload_date ON reports(upload_date);

-- Create index on processing_status for filtering by status
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(processing_status);

-- Create composite index for user reports by date (most common query)
CREATE INDEX IF NOT EXISTS idx_reports_user_date ON reports(user_id, upload_date DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_reports_user_date;
DROP INDEX IF EXISTS idx_reports_status;
DROP INDEX IF EXISTS idx_reports_upload_date;
DROP INDEX IF EXISTS idx_reports_user_id;
DROP TABLE IF EXISTS reports;
-- +goose StatementEnd
