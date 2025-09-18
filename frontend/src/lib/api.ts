// API Client for Medical Report Backend Integration

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// API Response Types
export interface ApiResponse<T = any> {
  data?: T;
  message?: string;
  error?: boolean;
  status?: number;
}

export interface AuthResponse {
  token?: string;
  message: string;
  success?: boolean;
}

export interface User {
  id: number;
  email: string;
  full_name: string;
  created_at: string;
  updated_at: string;
}

export interface SignupRequest {
  full_name: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface Report {
  id: number;
  user_id: number;
  original_filename: string;
  file_path: string;
  file_type: string;
  simplified_summary: string;
  upload_date: string;
  processed_at?: string;
}

export interface HealthMetric {
  name: string;
  value: string;
  unit: string;
  score: number; // 0-100 for speedometer
  status: "normal" | "warning" | "critical";
  range_min: number;
  range_max: number;
  description: string;
}

// API Error Class
export class ApiError extends Error {
  constructor(
    public status: number,
    public message: string,
    public response?: any
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// HTTP Client with automatic token handling
class HttpClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private getToken(): string | null {
    return localStorage.getItem('token');
  }

  private getHeaders(includeAuth = false): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (includeAuth) {
      const token = this.getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    const contentType = response.headers.get('Content-Type');
    const isJson = contentType?.includes('application/json');

    let data: any;
    if (isJson) {
      data = await response.json();
    } else {
      data = await response.text();
    }

    if (!response.ok) {
      const errorMessage = data?.message || data?.error || `HTTP ${response.status}`;
      throw new ApiError(response.status, errorMessage, data);
    }

    return data;
  }

  async get<T>(endpoint: string, options: { auth?: boolean } = {}): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: 'GET',
      headers: this.getHeaders(options.auth),
    });

    return this.handleResponse<T>(response);
  }

  async post<T>(
    endpoint: string,
    body?: any,
    options: { auth?: boolean } = {}
  ): Promise<T> {
    const isFormData = body instanceof FormData;

    const headers = this.getHeaders(options.auth);
    if (isFormData) {
      // Remove Content-Type header for FormData - browser will set it automatically with boundary
      delete (headers as any)['Content-Type'];
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: 'POST',
      headers,
      body: isFormData ? body : JSON.stringify(body),
    });

    return this.handleResponse<T>(response);
  }

  async put<T>(
    endpoint: string,
    body?: any,
    options: { auth?: boolean } = {}
  ): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: 'PUT',
      headers: this.getHeaders(options.auth),
      body: JSON.stringify(body),
    });

    return this.handleResponse<T>(response);
  }

  async delete<T>(
    endpoint: string,
    options: { auth?: boolean } = {}
  ): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: 'DELETE',
      headers: this.getHeaders(options.auth),
    });

    return this.handleResponse<T>(response);
  }
}

// Initialize HTTP client
const httpClient = new HttpClient(API_BASE_URL);

// Authentication API
export const authApi = {
  async signup(data: SignupRequest): Promise<AuthResponse> {
    return httpClient.post<AuthResponse>('/api/auth/signup', data);
  },

  async login(data: LoginRequest): Promise<AuthResponse> {
    return httpClient.post<AuthResponse>('/api/auth/login', data);
  },

  async logout(): Promise<AuthResponse> {
    const result = await httpClient.post<AuthResponse>('/api/auth/logout', {}, { auth: true });
    // Clear token from localStorage
    localStorage.removeItem('token');
    return result;
  },

  async getMe(): Promise<User> {
    return httpClient.get<User>('/api/auth/me', { auth: true });
  },

  async refreshToken(): Promise<AuthResponse> {
    return httpClient.post<AuthResponse>('/api/auth/refresh', {}, { auth: true });
  },

  // Helper methods
  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  },

  setToken(token: string): void {
    localStorage.setItem('token', token);
  },

  removeToken(): void {
    localStorage.removeItem('token');
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  }
};

// Reports API (placeholder for future implementation)
export const reportsApi = {
  async upload(file: File): Promise<Report> {
    const formData = new FormData();
    formData.append('file', file);

    return httpClient.post<Report>('/api/reports', formData, { auth: true });
  },

  async getAll(): Promise<Report[]> {
    return httpClient.get<Report[]>('/api/reports', { auth: true });
  },

  async getById(id: number): Promise<Report> {
    return httpClient.get<Report>(`/api/reports/${id}`, { auth: true });
  },

  async delete(id: number): Promise<void> {
    return httpClient.delete<void>(`/api/reports/${id}`, { auth: true });
  },

  async getSummary(id: number): Promise<{ report: Report; summary: string }> {
    return httpClient.get<{ report: Report; summary: string }>(`/api/reports/${id}/summary`, { auth: true });
  },

  async getHealthMetrics(id: number): Promise<{ report_id: number; metrics: HealthMetric[]; status: string }> {
    return httpClient.get<{ report_id: number; metrics: HealthMetric[]; status: string }>(`/api/reports/${id}/metrics`, { auth: true });
  }
};

// Chat API (placeholder for future implementation)
export const chatApi = {
  async sendMessage(reportId: number, message: string): Promise<any> {
    return httpClient.post(
      `/api/reports/${reportId}/chat`,
      { message },
      { auth: true }
    );
  },

  async getHistory(reportId: number): Promise<any[]> {
    return httpClient.get(`/api/reports/${reportId}/chat`, { auth: true });
  },

  async deleteMessage(reportId: number, messageId: number): Promise<void> {
    return httpClient.delete(`/api/reports/${reportId}/chat/${messageId}`, { auth: true });
  }
};

// Health check API
export const healthApi = {
  async check(): Promise<{ status: string; service: string; version: string }> {
    return httpClient.get('/health');
  }
};

// Token management utilities
export const tokenManager = {
  isTokenExpired(token: string): boolean {
    try {
      // JWT tokens have 3 parts separated by dots
      const parts = token.split('.');
      if (parts.length !== 3) return true;

      // Decode payload (second part)
      const payload = JSON.parse(atob(parts[1]));

      // Check expiration time (exp is in seconds, Date.now() is in milliseconds)
      const currentTime = Math.floor(Date.now() / 1000);
      return payload.exp < currentTime;
    } catch {
      return true;
    }
  },

  async autoRefreshToken(): Promise<boolean> {
    const token = authApi.getToken();
    if (!token) return false;

    try {
      if (this.isTokenExpired(token)) {
        const response = await authApi.refreshToken();
        if (response.token) {
          authApi.setToken(response.token);
          return true;
        }
      }
      return true;
    } catch {
      authApi.removeToken();
      return false;
    }
  }
};

// Export default API object
export default {
  auth: authApi,
  reports: reportsApi,
  chat: chatApi,
  health: healthApi,
  tokenManager,
  ApiError
};