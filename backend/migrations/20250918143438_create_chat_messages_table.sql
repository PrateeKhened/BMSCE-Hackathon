-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL,
    user_message TEXT NOT NULL,
    ai_response TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (report_id) REFERENCES reports(id) ON DELETE CASCADE
);

-- Create index on report_id for faster chat history queries
CREATE INDEX IF NOT EXISTS idx_chat_messages_report_id ON chat_messages(report_id);

-- Create index on created_at for chronological ordering
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);

-- Create composite index for report chat history (most common query)
CREATE INDEX IF NOT EXISTS idx_chat_messages_report_date ON chat_messages(report_id, created_at ASC);

-- Create index on non-deleted messages for active chat queries
CREATE INDEX IF NOT EXISTS idx_chat_messages_active ON chat_messages(is_deleted) WHERE is_deleted = FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chat_messages_active;
DROP INDEX IF EXISTS idx_chat_messages_report_date;
DROP INDEX IF EXISTS idx_chat_messages_created_at;
DROP INDEX IF EXISTS idx_chat_messages_report_id;
DROP TABLE IF EXISTS chat_messages;
-- +goose StatementEnd
