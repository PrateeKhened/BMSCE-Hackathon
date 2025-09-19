import React, { useState, useEffect } from 'react';
import Login from '../components/Login';
import FileUpload from '../components/FileUpload';
import ReportAnalysisDisplay from '../components/ReportAnalysisDisplay';
import { apiService } from '../services/apiService';

type AppState = 'login' | 'upload' | 'analysis' | 'loading';

const MedicalAnalysisDemo: React.FC = () => {
  const [appState, setAppState] = useState<AppState>('login');
  const [currentReportId, setCurrentReportId] = useState<number | null>(null);

  // Check if user is already logged in
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      setAppState('upload');
    }
  }, []);

  const handleLoginSuccess = () => {
    setAppState('upload');
  };

  const handleUploadSuccess = (reportId: number) => {
    setCurrentReportId(reportId);
    setAppState('loading');

    // Wait for AI processing (15 seconds) then show analysis
    setTimeout(() => {
      setAppState('analysis');
    }, 15000);
  };

  const handleLogout = () => {
    apiService.logout();
    setAppState('login');
    setCurrentReportId(null);
  };

  const handleNewUpload = () => {
    setAppState('upload');
    setCurrentReportId(null);
  };

  const renderHeader = () => {
    if (appState === 'login') return null;

    return (
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <span className="text-2xl mr-2">🏥</span>
              <h1 className="text-xl font-semibold text-gray-900">
                Medical Report Analyzer
              </h1>
            </div>
            <div className="flex items-center space-x-4">
              {(appState === 'analysis' || appState === 'loading') && (
                <button
                  onClick={handleNewUpload}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  New Upload
                </button>
              )}
              <button
                onClick={handleLogout}
                className="text-gray-600 hover:text-gray-800 px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>
    );
  };

  const renderContent = () => {
    switch (appState) {
      case 'login':
        return <Login onLoginSuccess={handleLoginSuccess} />;

      case 'upload':
        return (
          <div className="min-h-screen bg-gray-50 py-12">
            <FileUpload onUploadSuccess={handleUploadSuccess} />
          </div>
        );

      case 'loading':
        return (
          <div className="min-h-screen bg-gray-50 flex items-center justify-center">
            <div className="bg-white rounded-xl shadow-lg p-12 text-center max-w-md">
              <div className="text-6xl mb-6">🔬</div>
              <h2 className="text-2xl font-bold text-gray-900 mb-4">
                AI is Analyzing Your Report
              </h2>
              <div className="flex justify-center mb-6">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
              </div>
              <p className="text-gray-600 mb-4">
                Our AI is carefully analyzing your medical report and extracting health metrics.
              </p>
              <div className="bg-blue-50 rounded-lg p-4">
                <p className="text-sm text-blue-800">
                  This usually takes about 15 seconds. Your results will appear automatically.
                </p>
              </div>
            </div>
          </div>
        );

      case 'analysis':
        return currentReportId ? (
          <ReportAnalysisDisplay reportId={currentReportId} />
        ) : (
          <div className="min-h-screen bg-gray-50 flex items-center justify-center">
            <div className="bg-white rounded-xl shadow-lg p-8 text-center">
              <h2 className="text-xl font-semibold text-gray-800">No Report Selected</h2>
              <button
                onClick={handleNewUpload}
                className="mt-4 bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors"
              >
                Upload New Report
              </button>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {renderHeader()}
      {renderContent()}
    </div>
  );
};

export default MedicalAnalysisDemo;