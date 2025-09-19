import React, { useState, useEffect } from 'react';
import NewFileUpload from './NewFileUpload';
import NewMedicalAnalysis from './NewMedicalAnalysis';
import { newApiService } from '../services/newApiService';

type AppState = 'login' | 'upload' | 'analysis';

const NewApp: React.FC = () => {
  const [currentState, setCurrentState] = useState<AppState>('login');
  const [reportId, setReportId] = useState<number | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [authError, setAuthError] = useState<string | null>(null);

  // Login form state
  const [email, setEmail] = useState('demo@example.com');
  const [password, setPassword] = useState('demo123');
  const [fullName, setFullName] = useState('Demo User');
  const [isSignup, setIsSignup] = useState(false);
  const [authLoading, setAuthLoading] = useState(false);

  useEffect(() => {
    // Check if user is already logged in
    const token = localStorage.getItem('token');
    if (token) {
      setIsAuthenticated(true);
      setCurrentState('upload');
    }
  }, []);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setAuthLoading(true);
    setAuthError(null);

    try {
      if (isSignup) {
        console.log('📝 Attempting signup...');
        await newApiService.signup(email, password, fullName);
        console.log('✅ Signup successful');
      } else {
        console.log('🔑 Attempting login...');
        await newApiService.login(email, password);
        console.log('✅ Login successful');
      }

      setIsAuthenticated(true);
      setCurrentState('upload');
    } catch (error) {
      console.error('❌ Auth error:', error);
      setAuthError(error instanceof Error ? error.message : 'Authentication failed');
    } finally {
      setAuthLoading(false);
    }
  };

  const handleLogout = () => {
    newApiService.logout();
    setIsAuthenticated(false);
    setCurrentState('login');
    setReportId(null);
  };

  const handleUploadSuccess = (uploadedReportId: number) => {
    console.log('🎉 Upload successful, report ID:', uploadedReportId);
    setReportId(uploadedReportId);
    setCurrentState('analysis');
  };

  const handleBackToUpload = () => {
    setCurrentState('upload');
    setReportId(null);
  };

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4">
        <div className="max-w-md w-full">
          <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-200">
            <div className="text-center mb-8">
              <h1 className="text-3xl font-bold text-gray-900 mb-2 flex items-center justify-center">
                <span className="text-blue-600 mr-3">🏥</span>
                Medical Report AI
              </h1>
              <p className="text-gray-600">
                {isSignup ? 'Create your account to get started' : 'Sign in to analyze your medical reports'}
              </p>
            </div>

            <form onSubmit={handleLogin} className="space-y-6">
              {isSignup && (
                <div>
                  <label htmlFor="fullName" className="block text-sm font-medium text-gray-700 mb-2">
                    Full Name
                  </label>
                  <input
                    type="text"
                    id="fullName"
                    value={fullName}
                    onChange={(e) => setFullName(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    required
                  />
                </div>
              )}

              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                  Email Address
                </label>
                <input
                  type="email"
                  id="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                  Password
                </label>
                <input
                  type="password"
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>

              {authError && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 flex items-start">
                  <span className="text-red-500 mr-2 mt-0.5">⚠️</span>
                  <p className="text-red-700 font-medium">{authError}</p>
                </div>
              )}

              <button
                type="submit"
                disabled={authLoading}
                className={`w-full py-3 px-4 rounded-lg font-semibold text-white transition-all duration-200 ${
                  authLoading
                    ? 'bg-gray-300 cursor-not-allowed'
                    : 'bg-blue-600 hover:bg-blue-700 active:bg-blue-800 shadow-lg hover:shadow-xl'
                }`}
              >
                {authLoading ? (
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                    {isSignup ? 'Creating Account...' : 'Signing In...'}
                  </div>
                ) : (
                  isSignup ? 'Create Account' : 'Sign In'
                )}
              </button>
            </form>

            <div className="mt-6 text-center">
              <button
                onClick={() => {
                  setIsSignup(!isSignup);
                  setAuthError(null);
                }}
                className="text-blue-600 hover:text-blue-700 font-medium"
                disabled={authLoading}
              >
                {isSignup
                  ? 'Already have an account? Sign in'
                  : "Don't have an account? Sign up"
                }
              </button>
            </div>

            {/* Demo credentials hint */}
            {!isSignup && (
              <div className="mt-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
                <p className="text-blue-800 text-sm">
                  <strong>Demo credentials:</strong><br />
                  Email: demo@example.com<br />
                  Password: demo123
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Header */}
      <div className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900 flex items-center">
                <span className="text-blue-600 mr-2">🏥</span>
                Medical Report AI
              </h1>
            </div>
            <div className="flex items-center space-x-4">
              {currentState === 'analysis' && (
                <button
                  onClick={handleBackToUpload}
                  className="text-blue-600 hover:text-blue-700 font-medium flex items-center"
                >
                  <span className="mr-1">←</span>
                  Upload New Report
                </button>
              )}
              <button
                onClick={handleLogout}
                className="text-gray-600 hover:text-gray-700 font-medium flex items-center"
              >
                <span className="mr-1">🚪</span>
                Logout
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="py-8">
        {currentState === 'upload' && (
          <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-8">
              <h2 className="text-4xl font-bold text-gray-900 mb-4">
                AI-Powered Medical Report Analysis
              </h2>
              <p className="text-xl text-gray-600 max-w-2xl mx-auto">
                Upload your medical report and get instant AI analysis with health metrics,
                key findings, and personalized recommendations.
              </p>
            </div>
            <NewFileUpload onUploadSuccess={handleUploadSuccess} />
          </div>
        )}

        {currentState === 'analysis' && reportId && (
          <NewMedicalAnalysis reportId={reportId} />
        )}
      </div>
    </div>
  );
};

export default NewApp;