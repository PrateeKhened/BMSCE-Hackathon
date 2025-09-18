package types

import "time"

type Report struct {
	ID                int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	OriginalFilename string    `json:"original_filename" db:"original_filename"`
	FilePath         string    `json:"file_path" db:"file_path"`
	FileType         string    `json:"file_type" db:"file_type"`
	SimplifiedSummary string   `json:"simplified_summary" db:"simplified_summary"`
	UploadDate       time.Time `json:"upload_date" db:"upload_date"`
	ProcessedAt      *time.Time `json:"processed_at" db:"processed_at"`
}

type UploadRequest struct {
	File        []byte `json:"file"`
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}

type UploadResponse struct {
	Message  string `json:"message"`
	Success  bool   `json:"success"`
	ReportID int    `json:"report_id,omitempty"`
}

type ReportSummaryResponse struct {
	Report  Report `json:"report"`
	Summary string `json:"summary"`
}

type ChatMessage struct {
	ID         int       `json:"id" db:"id"`
	ReportID   int       `json:"report_id" db:"report_id"`
	UserMessage string   `json:"user_message" db:"user_message"`
	AIResponse string    `json:"ai_response" db:"ai_response"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ChatRequest struct {
	ReportID int    `json:"report_id" validate:"required"`
	Message  string `json:"message" validate:"required,min=1"`
}

type ChatResponse struct {
	Message   string        `json:"message"`
	Success   bool          `json:"success"`
	ChatData  *ChatMessage  `json:"chat_data,omitempty"`
}

type ReportListResponse struct {
	Reports []Report `json:"reports"`
	Total   int      `json:"total"`
}