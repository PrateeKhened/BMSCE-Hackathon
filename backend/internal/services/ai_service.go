package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/ledongthuc/pdf"
	"google.golang.org/api/option"
)

// HealthMetric represents a single health parameter with scoring
type HealthMetric struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`       // Can be string or number from AI
	Unit        string      `json:"unit"`
	Score       float64     `json:"score"`       // 0-100 score for speedometer
	Status      string      `json:"status"`      // "normal", "warning", "critical"
	RangeMin    float64     `json:"range_min"`   // Normal range minimum
	RangeMax    float64     `json:"range_max"`   // Normal range maximum
	Description string      `json:"description"` // Explanation for user
}

// GetValueAsString converts the value to string format for display
func (h *HealthMetric) GetValueAsString() string {
	switch v := h.Value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.1f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// AnalysisResult contains the complete AI analysis
type AnalysisResult struct {
	Summary         string          `json:"summary"`
	SimpleSummary   string          `json:"simple_summary"`
	HealthMetrics   []HealthMetric  `json:"health_metrics"`
	KeyFindings     []string        `json:"key_findings"`
	Recommendations []string        `json:"recommendations"`
	RiskLevel       string          `json:"risk_level"` // "low", "medium", "high"
}

// AIService handles AI-powered report analysis using Gemini
type AIService struct {
	client     *genai.Client
	model      *genai.GenerativeModel
	apiKey     string
	maxTokens  int32
}

// NewAIService creates a new AI service instance
func NewAIService(apiKey string) (*AIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Configure the model for medical report analysis
	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.3) // Lower temperature for more consistent medical analysis
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(2048)

	// Set safety settings for medical content
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	return &AIService{
		client:    client,
		model:     model,
		apiKey:    apiKey,
		maxTokens: 2048,
	}, nil
}

// AnalyzeReport processes a medical report file and returns comprehensive analysis
func (ai *AIService) AnalyzeReport(filePath, fileType string) (string, error) {
	fmt.Println("--- AI Service: AnalyzeReport ---")
	fmt.Println("File path:", filePath)
	fmt.Println("File type:", fileType)

	// Extract text content from file
	content, err := ai.extractTextFromFile(filePath, fileType)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from file: %w", err)
	}
	fmt.Println("Extracted content length:", len(content))

	// Generate comprehensive analysis
	analysis, err := ai.generateAnalysis(content)
	if err != nil {
		return "", fmt.Errorf("failed to generate AI analysis: %w", err)
	}

	// Convert to JSON for storage
	analysisJSON, err := json.Marshal(analysis)
	if err != nil {
		return "", fmt.Errorf("failed to serialize analysis: %w", err)
	}

	return string(analysisJSON), nil
}

// extractTextFromFile extracts text content based on file type
func (ai *AIService) extractTextFromFile(filePath, fileType string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt":
		return ai.extractFromTXT(filePath)
	case ".pdf":
		return ai.extractFromPDF(filePath)
	case ".docx", ".doc":
		return ai.extractFromDOCX(filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

// extractFromTXT reads plain text files
func (ai *AIService) extractFromTXT(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// extractFromPDF extracts text from PDF files using ledongthuc/pdf library
func (ai *AIService) extractFromPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var textContent strings.Builder
	totalPages := r.NumPage()

	// Extract text from all pages
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		content, err := page.GetPlainText(nil)
		if err != nil {
			// Log error but continue with other pages
			fmt.Printf("Warning: Failed to extract text from page %d: %v\n", pageNum, err)
			continue
		}

		textContent.WriteString(content)
		textContent.WriteString("\n")
	}

	extractedText := textContent.String()
	if strings.TrimSpace(extractedText) == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return extractedText, nil
}

// extractFromDOCX extracts text from DOCX files (placeholder - requires DOCX library)
func (ai *AIService) extractFromDOCX(filePath string) (string, error) {
	// TODO: Implement DOCX text extraction using a library like gingfrederik/docx
	// For now, return placeholder text
	return "DOCX text extraction not yet implemented. Please use TXT format for testing.", nil
}

// generateAnalysis uses Gemini to analyze medical report content
func (ai *AIService) generateAnalysis(content string) (*AnalysisResult, error) {
	ctx := context.Background()

	// Create comprehensive prompt for medical analysis
	prompt := ai.buildAnalysisPrompt(content)
	fmt.Println("--- AI Service: Prompt ---")
	fmt.Println(prompt)

	// Generate response from Gemini
	resp, err := ai.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// Extract text from response
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}
	fmt.Println("--- AI Service: Response ---")
	fmt.Println(responseText)

	// Parse the structured response
	analysis, err := ai.parseAnalysisResponse(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse analysis response: %w", err)
	}

	return analysis, nil
}

// loadPromptTemplate loads the medical analysis prompt template from file
func (ai *AIService) loadPromptTemplate() (string, error) {
	promptPath := "prompts/medical_analysis_prompt.txt"
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		// Fallback to embedded prompt if file doesn't exist
		return ai.getDefaultPromptTemplate(), nil
	}
	return string(promptBytes), nil
}

// getDefaultPromptTemplate returns a fallback prompt if file loading fails
func (ai *AIService) getDefaultPromptTemplate() string {
	return `You are a medical AI assistant specialized in analyzing medical reports and lab results. Please analyze the following medical report and provide a comprehensive analysis in JSON format.

Medical Report Content:
{{REPORT_CONTENT}}

Please provide your analysis in the following JSON structure:
{
  "summary": "Detailed medical summary for healthcare professionals",
  "simple_summary": "Easy-to-understand summary for patients (avoid medical jargon)",
  "health_metrics": [
    {
      "name": "Parameter name (e.g., Blood Glucose, Cholesterol)",
      "value": "Measured value (can be number or string)",
      "unit": "Unit of measurement",
      "score": "Score from 0-100 (100 = optimal, 0 = critical)",
      "status": "normal/warning/critical",
      "range_min": "Normal range minimum value",
      "range_max": "Normal range maximum value",
      "description": "Simple explanation of what this means"
    }
  ],
  "key_findings": ["List of important findings"],
  "recommendations": ["List of actionable recommendations"],
  "risk_level": "low/medium/high"
}

Guidelines:
1. Extract all measurable parameters (blood tests, vitals, etc.)
2. Provide scores based on how close values are to optimal ranges
3. Use simple language in simple_summary and descriptions
4. Be accurate but not alarming in tone
5. Include lifestyle recommendations when appropriate
6. If no specific values are found, focus on general health insights
7. For numeric values, you can return them as numbers in the JSON

Respond only with valid JSON.`
}

// buildAnalysisPrompt creates a comprehensive prompt for medical analysis
func (ai *AIService) buildAnalysisPrompt(content string) string {
	promptTemplate, err := ai.loadPromptTemplate()
	if err != nil {
		// Use default template if loading fails
		promptTemplate = ai.getDefaultPromptTemplate()
	}

	// Replace placeholder with actual content
	prompt := strings.ReplaceAll(promptTemplate, "{{REPORT_CONTENT}}", content)
	return prompt
}

// parseAnalysisResponse parses the AI response into structured data
func (ai *AIService) parseAnalysisResponse(response string) (*AnalysisResult, error) {
	// Clean response (remove markdown formatting if present)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Try to find JSON within the response (sometimes AI adds extra text)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart >= 0 && jsonEnd > jsonStart {
		response = response[jsonStart:jsonEnd+1]
	}

	var analysis AnalysisResult
	err := json.Unmarshal([]byte(response), &analysis)
	if err != nil {
		// Log the actual response for debugging
		fmt.Printf("Failed to parse JSON response: %s\nError: %v\n", response, err)

		// If JSON parsing fails, create a fallback analysis with the raw response
		return &AnalysisResult{
			Summary:       "AI analysis completed. Raw response formatting required improvement.",
			SimpleSummary: fmt.Sprintf("Analysis: %s", ai.extractSimpleSummary(response)),
			HealthMetrics: ai.extractHealthMetrics(response),
			KeyFindings:   []string{"Report analysis completed", "Response parsing needed enhancement"},
			Recommendations: []string{"Consult with your healthcare provider for personalized advice"},
			RiskLevel:     "medium",
		}, nil
	}

	// Validate and enhance the analysis
	ai.validateAndEnhanceAnalysis(&analysis)

	return &analysis, nil
}

// validateAndEnhanceAnalysis ensures the analysis meets quality standards
func (ai *AIService) validateAndEnhanceAnalysis(analysis *AnalysisResult) {
	// Ensure all required fields have content
	if analysis.Summary == "" {
		analysis.Summary = "Medical analysis completed."
	}
	if analysis.SimpleSummary == "" {
		analysis.SimpleSummary = "Your report has been analyzed. Please discuss with your healthcare provider."
	}
	if analysis.RiskLevel == "" {
		analysis.RiskLevel = "medium"
	}

	// Validate health metrics scores
	for i := range analysis.HealthMetrics {
		metric := &analysis.HealthMetrics[i]

		// Ensure score is within valid range
		if metric.Score < 0 {
			metric.Score = 0
		} else if metric.Score > 100 {
			metric.Score = 100
		}

		// Validate status matches score
		if metric.Status == "" {
			if metric.Score >= 80 {
				metric.Status = "normal"
			} else if metric.Score >= 50 {
				metric.Status = "warning"
			} else {
				metric.Status = "critical"
			}
		}
	}

	// Ensure we have at least one recommendation
	if len(analysis.Recommendations) == 0 {
		analysis.Recommendations = []string{
			"Regular health check-ups with your healthcare provider",
			"Maintain a balanced diet and regular exercise",
			"Follow any prescribed treatments consistently",
		}
	}
}

// GetHealthMetrics extracts health metrics from analysis for speedometer display
func (ai *AIService) GetHealthMetrics(analysisJSON string) ([]HealthMetric, error) {
	var analysis AnalysisResult
	err := json.Unmarshal([]byte(analysisJSON), &analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse analysis: %w", err)
	}

	return analysis.HealthMetrics, nil
}

// Close cleanly shuts down the AI service
func (ai *AIService) Close() error {
	if ai.client != nil {
		return ai.client.Close()
	}
	return nil
}

// extractSimpleSummary extracts a simple summary from raw AI response
func (ai *AIService) extractSimpleSummary(response string) string {
	// Try to extract meaningful content from the response
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 20 && !strings.HasPrefix(line, "{") && !strings.HasPrefix(line, "}") {
			// Return first meaningful line as summary
			return line
		}
	}
	return "Your medical report has been analyzed. Please consult your healthcare provider."
}

// extractHealthMetrics attempts to extract health metrics from raw response
func (ai *AIService) extractHealthMetrics(response string) []HealthMetric {
	// For now, return empty metrics if JSON parsing fails
	// In the future, we could implement regex parsing to extract numeric values
	return []HealthMetric{}
}

// Helper function to determine file content type from extension
func getContentTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".doc":
		return "application/msword"
	default:
		return "application/octet-stream"
	}
}