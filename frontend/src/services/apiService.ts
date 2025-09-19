// API Service for Medical Report Backend Integration

const API_BASE_URL = 'http://localhost:3001';

export interface HealthMetric {
  name: string;
  value: number | string;
  unit: string;
  score: number; // 0-100 for speedometer
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
    email_verified: boolean;
    is_active: boolean;
    created_at: string;
    updated_at: string;
  };
}

class ApiService {
  private getHeaders(includeAuth = true): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (includeAuth) {
      const token = localStorage.getItem('token');
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('token');
    return {
      'Authorization': `Bearer ${token}`,
    };
  }

  // Authentication
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
      method: 'POST',
      headers: this.getHeaders(false),
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const data = await response.json();
    localStorage.setItem('token', data.token);
    return data;
  }

  async signup(email: string, password: string, name: string): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/api/auth/signup`, {
      method: 'POST',
      headers: this.getHeaders(false),
      body: JSON.stringify({ email, password, name }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Signup failed');
    }

    const data = await response.json();
    localStorage.setItem('token', data.token);
    return data;
  }

  logout() {
    localStorage.removeItem('token');
  }

  // File Upload
  async uploadReport(file: File, description: string): Promise<{ report_id: number }> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('description', description);

    const response = await fetch(`${API_BASE_URL}/api/reports`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: formData,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Upload failed');
    }

    return await response.json();
  }

  // Get Reports
  async getAllReports(): Promise<{ reports: Report[]; total: number }> {
    const response = await fetch(`${API_BASE_URL}/api/reports`, {
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch reports');
    }

    return await response.json();
  }

  async getReport(reportId: number): Promise<{ report: Report }> {
    const response = await fetch(`${API_BASE_URL}/api/reports/${reportId}`, {
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch report');
    }

    return await response.json();
  }

  // Get Analysis Summary
  async getReportSummary(reportId: number): Promise<{ report: Report; summary: string }> {
    const response = await fetch(`${API_BASE_URL}/api/reports/${reportId}/summary`, {
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch report summary');
    }

    return await response.json();
  }

  // Get Health Metrics for Speedometer
  async getHealthMetrics(reportId: number): Promise<{
    metrics: HealthMetric[];
    report_id: number;
    status: string;
  }> {
    const response = await fetch(`${API_BASE_URL}/api/reports/${reportId}/metrics`, {
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch health metrics');
    }

    return await response.json();
  }

  // Parse Analysis from Report
  parseAnalysis(report: Report): AnalysisResult {
    try {
      return JSON.parse(report.simplified_summary);
    } catch (error) {
      console.error('Failed to parse analysis:', error);
      // Return fallback structure
      return {
        summary: 'Analysis parsing failed',
        simple_summary: 'Your report has been processed but could not be displayed properly.',
        health_metrics: [],
        key_findings: ['Report processed'],
        recommendations: ['Please contact support for assistance'],
        risk_level: 'medium',
      };
    }
  }

  // Health Check
  async healthCheck(): Promise<{ status: string }> {
    const response = await fetch(`${API_BASE_URL}/health`);
    return await response.json();
  }
}

export const apiService = new ApiService();
export default apiService;