# Medical Analysis Prompts

This directory contains configurable prompt templates for the AI medical report analysis system.

## Files

### `medical_analysis_prompt.txt`
The main system prompt template used by Gemini AI to analyze medical reports and extract health metrics.

## How to Modify the System Prompt

1. **Edit the prompt file**: Open `medical_analysis_prompt.txt` in any text editor
2. **Use placeholder**: Keep `{{REPORT_CONTENT}}` where the medical report content should be inserted
3. **Restart the server**: The server loads the prompt file on each analysis request, so changes take effect immediately

## Prompt Structure

The current prompt includes:

- **Role definition**: Establishes the AI as a medical assistant
- **JSON schema**: Defines the expected output structure
- **Guidelines**: 10 specific rules for analysis quality
- **Output format**: Ensures consistent speedometer data (0-100 scores)

## Key Features

- **Flexible value types**: Supports both numeric and string values from AI
- **Speedometer scoring**: All metrics scored 0-100 for UI display
- **Status classification**: normal/warning/critical based on scores
- **Patient-friendly**: Simple explanations alongside technical summaries
- **Fallback system**: Built-in default prompt if file loading fails

## Testing Changes

Upload a medical report via the API to test prompt modifications:

```bash
# Upload report
curl -X POST http://localhost:8080/api/reports \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@sample_report.txt" \
  -F "description=Test prompt changes"

# Check results
curl -X GET http://localhost:8080/api/reports/{id}/metrics \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Best Practices

1. **Keep JSON schema intact**: Don't modify the required fields
2. **Preserve guidelines**: The 10 guidelines ensure quality analysis
3. **Test thoroughly**: Upload different report types after changes
4. **Backup originals**: Keep copies of working prompts before major changes
5. **Validate JSON**: Ensure the AI response remains valid JSON format