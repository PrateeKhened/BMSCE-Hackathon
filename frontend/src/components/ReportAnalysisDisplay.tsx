import React, { useState, useEffect } from 'react';
import { AnalysisResult, Report, apiService } from '../services/apiService';
import HealthMetricCard from './HealthMetricCard';

interface ReportAnalysisDisplayProps {
  reportId: number;
}

const ReportAnalysisDisplay: React.FC<ReportAnalysisDisplayProps> = ({ reportId }) => {
  const [report, setReport] = useState<Report | null>(null);
  const [analysis, setAnalysis] = useState<AnalysisResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchReportAnalysis = async () => {
      try {
        setLoading(true);
        setError(null);

        // Get the report details
        const reportData = await apiService.getReport(reportId);
        setReport(reportData.report);

        // Get the AI analysis from the summary endpoint
        const summaryData = await apiService.getReportSummary(reportId);
        const analysisData = JSON.parse(summaryData.summary);
        setAnalysis(analysisData);

      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load report analysis');
        console.error('Error fetching report analysis:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchReportAnalysis();
  }, [reportId]);

  const getRiskLevelColor = (riskLevel: string): string => {
    switch (riskLevel) {
      case 'low': return 'bg-green-100 text-green-800 border-green-200';
      case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'high': return 'bg-red-100 text-red-800 border-red-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="bg-white rounded-lg shadow-lg p-8 text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <h2 className="text-xl font-semibold text-gray-800 mb-2">Analyzing Your Report</h2>
          <p className="text-gray-600">This usually takes about 15 seconds...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="bg-white rounded-lg shadow-lg p-8 text-center max-w-md">
          <div className="text-red-500 text-5xl mb-4">⚠️</div>
          <h2 className="text-xl font-semibold text-gray-800 mb-2">Error Loading Report</h2>
          <p className="text-gray-600 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  if (!report || !analysis) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="bg-white rounded-lg shadow-lg p-8 text-center">
          <h2 className="text-xl font-semibold text-gray-800">No Data Available</h2>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between">
            <div className="mb-4 lg:mb-0">
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Medical Report Analysis</h1>
              <p className="text-gray-600">
                <span className="font-medium">File:</span> {report.original_filename}
              </p>
              <p className="text-gray-600">
                <span className="font-medium">Uploaded:</span> {formatDate(report.upload_date)}
              </p>
              {report.processed_at && (
                <p className="text-gray-600">
                  <span className="font-medium">Processed:</span> {formatDate(report.processed_at)}
                </p>
              )}
            </div>

            {/* Risk Level Badge */}
            <div className={`inline-flex items-center px-4 py-2 rounded-full border ${getRiskLevelColor(analysis.risk_level)}`}>
              <span className="font-medium">
                Risk Level: {analysis.risk_level.toUpperCase()}
              </span>
            </div>
          </div>
        </div>

        {/* Simple Summary */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 flex items-center">
            <span className="text-blue-600 mr-2">📋</span>
            Summary
          </h2>
          <p className="text-lg text-gray-700 leading-relaxed">
            {analysis.simple_summary}
          </p>
        </div>

        {/* Health Metrics Grid */}
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-6 flex items-center">
            <span className="text-green-600 mr-2">📊</span>
            Health Metrics
          </h2>

          {analysis.health_metrics.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {analysis.health_metrics.map((metric, index) => (
                <HealthMetricCard key={index} metric={metric} />
              ))}
            </div>
          ) : (
            <div className="bg-white rounded-xl shadow-lg p-8 text-center">
              <p className="text-gray-600">No specific health metrics were extracted from this report.</p>
            </div>
          )}
        </div>

        {/* Key Findings and Recommendations */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          {/* Key Findings */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h3 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
              <span className="text-yellow-600 mr-2">🔍</span>
              Key Findings
            </h3>
            <ul className="space-y-3">
              {analysis.key_findings.map((finding, index) => (
                <li key={index} className="flex items-start">
                  <span className="text-blue-600 mr-2 mt-1">•</span>
                  <span className="text-gray-700">{finding}</span>
                </li>
              ))}
            </ul>
          </div>

          {/* Recommendations */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h3 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
              <span className="text-purple-600 mr-2">💡</span>
              Recommendations
            </h3>
            <ul className="space-y-3">
              {analysis.recommendations.map((recommendation, index) => (
                <li key={index} className="flex items-start">
                  <span className="text-green-600 mr-2 mt-1">✓</span>
                  <span className="text-gray-700">{recommendation}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Detailed Summary for Healthcare Professionals */}
        <div className="bg-white rounded-xl shadow-lg p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 flex items-center">
            <span className="text-red-600 mr-2">🏥</span>
            Detailed Medical Summary
          </h2>
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">
              {analysis.summary}
            </p>
          </div>
        </div>

        {/* Footer */}
        <div className="mt-8 text-center">
          <p className="text-gray-500 text-sm">
            This analysis is generated by AI and should not replace professional medical advice.
            Always consult with your healthcare provider for medical decisions.
          </p>
        </div>
      </div>
    </div>
  );
};

export default ReportAnalysisDisplay;