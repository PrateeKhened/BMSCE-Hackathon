// Simple, working API service for Medical Report Analysis
const API_BASE_URL = 'http://localhost:3001';

// Types
export interface HealthMetric {
  name: string;
  value: number | string;
  unit: string;
  score: number;
  status: 'normal' | 'warning' | 'critical';
  range_min: number;
  range_max: number;
  description: string;
}

export interface AnalysisResult {
  summary: string;
  simple_summary: string;
  health_metrics: HealthMetric[];
  key_findings: string[];
  recommendations: string[];
  risk_level: 'low' | 'medium' | 'high';
}

export interface Report {
  id: number;
  user_id: number;
  original_filename: string;
  file_path: string;
  file_type: string;
  simplified_summary: string;
  upload_date: string;
  processed_at: string;
}

export interface AuthResponse {
  token: string;
  user: {
    id: number;
    email: string;
    full_name: string;
  };
}

class NewApiService {
  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('token');
    return {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` })
    };
  }

  // Authentication
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const data = await response.json();
    localStorage.setItem('token', data.token);
    return data;
  }

  async signup(email: string, password: string, fullName: string): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/api/auth/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email,
        password,
        full_name: fullName
      })
    });

    if (!response.ok) {
      throw new Error('Signup failed');
    }

    const data = await response.json();
    localStorage.setItem('token', data.token);
    return data;
  }

  logout(): void {
    localStorage.removeItem('token');
  }

  // File Upload
  async uploadReport(file: File, description: string): Promise<{ report_id: number; message: string }> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('description', description);

    const token = localStorage.getItem('token');
    const response = await fetch(`${API_BASE_URL}/api/reports`, {
      method: 'POST',
      headers: {
        ...(token && { 'Authorization': `Bearer ${token}` })
      },
      body: formData
    });

    if (!response.ok) {
      throw new Error('Upload failed');
    }

    return response.json();
  }

  // Get Report Details
  async getReport(reportId: number): Promise<{ report: Report }> {
    const response = await fetch(`${API_BASE_URL}/api/reports/${reportId}`, {
      method: 'GET',
      headers: this.getAuthHeaders()
    });

    if (!response.ok) {
      throw new Error('Failed to get report');
    }

    return response.json();
  }

  // Get AI Analysis - THIS IS THE KEY METHOD
  async getAnalysis(reportId: number): Promise<AnalysisResult> {
    console.log(`🔍 Fetching analysis for report ${reportId}`);

    const response = await fetch(`${API_BASE_URL}/api/reports/${reportId}/summary`, {
      method: 'GET',
      headers: this.getAuthHeaders()
    });

    if (!response.ok) {
      console.error(`❌ Analysis fetch failed: ${response.status} ${response.statusText}`);
      throw new Error('Failed to get analysis');
    }

    const data = await response.json();
    console.log('📊 Raw analysis response:', data);

    // The backend returns: { report: {...}, summary: "JSON string" }
    // We need to parse the summary field which is a JSON string
    if (!data.summary) {
      throw new Error('No summary found in response');
    }

    console.log('📊 Raw summary string:', data.summary);
    console.log('📊 Summary type:', typeof data.summary);
    console.log('📊 Summary length:', data.summary.length);

    let analysis: AnalysisResult;
    try {
      analysis = JSON.parse(data.summary);
      console.log('✅ Parsed analysis:', analysis);
      console.log('✅ Health metrics in parsed:', analysis.health_metrics?.length || 0);
    } catch (parseError) {
      console.error('❌ JSON Parse Error:', parseError);
      console.error('❌ Failed to parse summary:', data.summary.substring(0, 200) + '...');
      throw new Error(`Failed to parse analysis JSON: ${parseError instanceof Error ? parseError.message : 'Unknown error'}`);
    }

    return analysis;
  }
}

export const newApiService = new NewApiService();