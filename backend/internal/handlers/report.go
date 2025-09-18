package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/middleware"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/services"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/errors"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/types"
)

// ReportHandler handles report HTTP requests
type ReportHandler struct {
	reportRepo      models.ReportRepository
	authService     *services.AuthService
	aiService       *services.AIService
	uploadDirectory string
	maxFileSize     int64
}

// NewReportHandler creates a new report handler
func NewReportHandler(
	reportRepo models.ReportRepository,
	authService *services.AuthService,
	aiService *services.AIService,
	uploadDir string,
	maxFileSize int64,
) *ReportHandler {
	return &ReportHandler{
		reportRepo:      reportRepo,
		authService:     authService,
		aiService:       aiService,
		uploadDirectory: uploadDir,
		maxFileSize:     maxFileSize,
	}
}

// UploadReportHandler handles file upload requests
// POST /api/reports
func (rh *ReportHandler) UploadReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user from context (set by auth middleware)
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse multipart form with size limit
	err := r.ParseMultipartForm(rh.maxFileSize)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "File too large or invalid form data")
		return
	}

	// Get the uploaded file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "No file provided or invalid file field")
		return
	}
	defer file.Close()

	// Validate file type and size
	if err := rh.validateFile(fileHeader); err != nil {
		handleServiceError(w, err)
		return
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(rh.uploadDirectory, 0755); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	uniqueFilename := rh.generateUniqueFilename(fileHeader.Filename)
	filePath := filepath.Join(rh.uploadDirectory, uniqueFilename)

	// Save file to disk
	if err := rh.saveFile(file, filePath); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Create report record in database
	report := &models.Report{
		UserID:           user.ID,
		OriginalFilename: fileHeader.Filename,
		FilePath:         filePath,
		FileType:         fileHeader.Header.Get("Content-Type"),
		FileSize:         fileHeader.Size,
		ProcessingStatus: "pending",
	}

	if err := rh.reportRepo.Create(report); err != nil {
		// Clean up uploaded file on database error
		os.Remove(filePath)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to save report metadata")
		return
	}

	// Trigger async AI processing
	go rh.processReportAsync(report)

	// Return success response
	response := types.UploadResponse{
		Message:  "File uploaded successfully and queued for processing",
		Success:  true,
		ReportID: report.ID,
	}

	writeJSONResponse(w, http.StatusCreated, response)
}

// GetReportsHandler retrieves user's reports with pagination
// GET /api/reports
func (rh *ReportHandler) GetReportsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse pagination parameters
	limit, offset := rh.parsePaginationParams(r)

	// Get reports from database
	reports, err := rh.reportRepo.GetByUserID(user.ID, limit, offset)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve reports")
		return
	}

	// Convert to response format
	reportResponses := make([]types.Report, len(reports))
	for i, report := range reports {
		reportResponses[i] = types.Report{
			ID:                report.ID,
			UserID:           report.UserID,
			OriginalFilename: report.OriginalFilename,
			FilePath:         report.FilePath,
			FileType:         report.FileType,
			SimplifiedSummary: report.SimplifiedSummary,
			UploadDate:       report.UploadDate,
			ProcessedAt:      report.ProcessedAt,
		}
	}

	response := types.ReportListResponse{
		Reports: reportResponses,
		Total:   len(reportResponses),
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// GetReportHandler retrieves a specific report by ID
// GET /api/reports/{id}
func (rh *ReportHandler) GetReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract report ID from URL
	vars := mux.Vars(r)
	reportID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	// Get report from database
	report, err := rh.reportRepo.GetByID(reportID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve report")
		return
	}

	if report == nil {
		writeErrorResponse(w, http.StatusNotFound, "Report not found")
		return
	}

	// Check if user owns this report
	if report.UserID != user.ID {
		writeErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Convert to response format
	reportResponse := types.Report{
		ID:                report.ID,
		UserID:           report.UserID,
		OriginalFilename: report.OriginalFilename,
		FilePath:         report.FilePath,
		FileType:         report.FileType,
		SimplifiedSummary: report.SimplifiedSummary,
		UploadDate:       report.UploadDate,
		ProcessedAt:      report.ProcessedAt,
	}

	writeJSONResponse(w, http.StatusOK, reportResponse)
}

// DeleteReportHandler deletes a report and its file
// DELETE /api/reports/{id}
func (rh *ReportHandler) DeleteReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract report ID from URL
	vars := mux.Vars(r)
	reportID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	// Get report to check ownership and get file path
	report, err := rh.reportRepo.GetByID(reportID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve report")
		return
	}

	if report == nil {
		writeErrorResponse(w, http.StatusNotFound, "Report not found")
		return
	}

	// Check if user owns this report
	if report.UserID != user.ID {
		writeErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete from database first
	if err := rh.reportRepo.Delete(reportID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete report")
		return
	}

	// Delete file from filesystem (ignore errors for cleanup)
	os.Remove(report.FilePath)

	response := map[string]any{
		"message": "Report deleted successfully",
		"success": true,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// validateFile checks file type and size constraints
func (rh *ReportHandler) validateFile(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > rh.maxFileSize {
		return errors.NewValidationError("File size exceeds maximum limit of 20MB")
	}

	// Check file extension
	filename := strings.ToLower(fileHeader.Filename)
	allowedExtensions := []string{".pdf", ".txt", ".docx", ".doc"}

	isAllowed := false
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(filename, ext) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return errors.NewValidationError("File type not supported. Please upload PDF, TXT, or DOCX files only")
	}

	// Additional content-type validation
	contentType := fileHeader.Header.Get("Content-Type")
	allowedTypes := []string{
		"application/pdf",
		"text/plain",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/msword",
	}

	isValidContentType := false
	for _, allowedType := range allowedTypes {
		if strings.Contains(contentType, allowedType) {
			isValidContentType = true
			break
		}
	}

	if !isValidContentType {
		return errors.NewValidationError("Invalid file content type")
	}

	return nil
}

// generateUniqueFilename creates a unique filename to prevent conflicts
func (rh *ReportHandler) generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)

	// Use timestamp and a portion of original name for uniqueness
	timestamp := time.Now().Unix()

	// Sanitize filename (remove special characters)
	safeFilename := strings.ReplaceAll(nameWithoutExt, " ", "_")
	safeFilename = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return -1
	}, safeFilename)

	return fmt.Sprintf("%d_%s%s", timestamp, safeFilename, ext)
}

// saveFile writes the uploaded file to disk
func (rh *ReportHandler) saveFile(src multipart.File, filePath string) error {
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// parsePaginationParams extracts limit and offset from query parameters
func (rh *ReportHandler) parsePaginationParams(r *http.Request) (limit, offset int) {
	// Default values
	limit = 20
	offset = 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}

// GetReportSummaryHandler returns the AI-generated summary and analysis
// GET /api/reports/{id}/summary
func (rh *ReportHandler) GetReportSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract report ID from URL
	vars := mux.Vars(r)
	reportID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	// Get report from database
	report, err := rh.reportRepo.GetByID(reportID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve report")
		return
	}

	if report == nil {
		writeErrorResponse(w, http.StatusNotFound, "Report not found")
		return
	}

	// Check if user owns this report
	if report.UserID != user.ID {
		writeErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if report has been processed
	if report.ProcessingStatus != "completed" {
		writeErrorResponse(w, http.StatusBadRequest, "Report is not ready yet")
		return
	}

	response := types.ReportSummaryResponse{
		Report: types.Report{
			ID:                report.ID,
			UserID:           report.UserID,
			OriginalFilename: report.OriginalFilename,
			FilePath:         report.FilePath,
			FileType:         report.FileType,
			SimplifiedSummary: report.SimplifiedSummary,
			UploadDate:       report.UploadDate,
			ProcessedAt:      report.ProcessedAt,
		},
		Summary: report.SimplifiedSummary,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// GetHealthMetricsHandler returns health metrics for speedometer display
// GET /api/reports/{id}/metrics
func (rh *ReportHandler) GetHealthMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract report ID from URL
	vars := mux.Vars(r)
	reportID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	// Get report from database
	report, err := rh.reportRepo.GetByID(reportID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve report")
		return
	}

	if report == nil {
		writeErrorResponse(w, http.StatusNotFound, "Report not found")
		return
	}

	// Check if user owns this report
	if report.UserID != user.ID {
		writeErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if report has been processed
	if report.ProcessingStatus != "completed" {
		writeErrorResponse(w, http.StatusBadRequest, "Report is not ready yet")
		return
	}

	// Check if AI service is available
	if rh.aiService == nil {
		writeErrorResponse(w, http.StatusServiceUnavailable, "AI service not available")
		return
	}

	// Extract health metrics from AI analysis
	healthMetrics, err := rh.aiService.GetHealthMetrics(report.SimplifiedSummary)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to extract health metrics")
		return
	}

	response := map[string]any{
		"report_id": report.ID,
		"metrics":   healthMetrics,
		"status":    "completed",
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// processReportAsync handles AI processing in background
func (rh *ReportHandler) processReportAsync(report *models.Report) {
	// Update status to processing
	rh.reportRepo.UpdateProcessingStatus(report.ID, "processing", "")

	// Check if AI service is available
	if rh.aiService == nil {
		rh.reportRepo.UpdateProcessingStatus(report.ID, "failed", "AI service not available - missing API key")
		return
	}

	// Extract text from file and get AI analysis
	summary, err := rh.aiService.AnalyzeReport(report.FilePath, report.FileType)
	if err != nil {
		// Update status to failed
		rh.reportRepo.UpdateProcessingStatus(report.ID, "failed", fmt.Sprintf("Processing failed: %v", err))
		return
	}

	// Update status to completed with summary
	rh.reportRepo.UpdateProcessingStatus(report.ID, "completed", summary)
}