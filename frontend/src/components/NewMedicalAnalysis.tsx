import React, { useState, useEffect } from 'react';
import { AnalysisResult, Report, newApiService } from '../services/newApiService';
import SimpleSpeedometer from './SimpleSpeedometer';

interface NewMedicalAnalysisProps {
  reportId: number;
}

const NewMedicalAnalysis: React.FC<NewMedicalAnalysisProps> = ({ reportId }) => {
  const [report, setReport] = useState<Report | null>(null);
  const [analysis, setAnalysis] = useState<AnalysisResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAnalysis = async () => {
      try {
        setLoading(true);
        setError(null);
        console.log('🔍 Fetching analysis for report ID:', reportId);

        // Get report details
        const reportResponse = await newApiService.getReport(reportId);
        console.log('📄 Report data:', reportResponse);
        setReport(reportResponse.report);

        // Check if AI processing is complete
        if (!reportResponse.report.processed_at) {
          console.log('⏳ AI processing not complete yet, waiting...');

          // Poll for completion every 3 seconds
          const pollForCompletion = async (): Promise<void> => {
            return new Promise((resolve, reject) => {
              const checkProcessing = async () => {
                try {
                  const updatedReportResponse = await newApiService.getReport(reportId);
                  console.log('🔄 Checking processing status...', updatedReportResponse.report.processed_at);

                  if (updatedReportResponse.report.processed_at) {
                    console.log('✅ Processing complete!');
                    setReport(updatedReportResponse.report);
                    resolve();
                  } else {
                    console.log('⏳ Still processing, checking again in 3 seconds...');
                    setTimeout(checkProcessing, 3000);
                  }
                } catch (error) {
                  reject(error);
                }
              };

              checkProcessing();
            });
          };

          await pollForCompletion();
        }

        // Get AI analysis (now that processing is complete)
        const analysisData = await newApiService.getAnalysis(reportId);
        console.log('📊 Analysis data:', analysisData);
        console.log('📊 Health metrics count:', analysisData.health_metrics?.length || 0);
        console.log('📊 First metric:', analysisData.health_metrics?.[0]);

        // Validate analysis data
        if (!analysisData) {
          throw new Error('No analysis data received');
        }

        if (!analysisData.health_metrics || !Array.isArray(analysisData.health_metrics)) {
          throw new Error('Invalid health metrics data');
        }

        setAnalysis(analysisData);

      } catch (err) {
        console.error('❌ Error fetching analysis:', err);
        console.error('❌ Error stack:', err instanceof Error ? err.stack : 'No stack');
        setError(err instanceof Error ? err.message : 'Failed to load analysis');
      } finally {
        setLoading(false);
      }
    };

    if (reportId) {
      fetchAnalysis();
    }
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
        <div className="bg-white rounded-xl shadow-lg p-8 text-center max-w-md">
          <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600 mx-auto mb-6"></div>
          <h2 className="text-2xl font-semibold text-gray-800 mb-3">🤖 AI Analyzing Your Medical Report</h2>
          <p className="text-gray-600 mb-2">Our AI is extracting health metrics and insights...</p>
          <p className="text-sm text-gray-500 mb-3">This usually takes 15-30 seconds</p>
          <div className="bg-blue-50 rounded-lg p-4 text-left">
            <p className="text-sm text-blue-800 font-medium mb-2">Processing Steps:</p>
            <ul className="text-xs text-blue-700 space-y-1">
              <li>✓ File uploaded successfully</li>
              <li>🔄 Extracting text content</li>
              <li>🧠 AI analyzing health metrics</li>
              <li>📊 Generating speedometer visualizations</li>
            </ul>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="bg-white rounded-xl shadow-lg p-8 text-center max-w-md">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-xl font-semibold text-gray-800 mb-3">Analysis Error</h2>
          <p className="text-gray-600 mb-6">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors font-medium"
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
        <div className="bg-white rounded-xl shadow-lg p-8 text-center">
          <div className="text-gray-400 text-6xl mb-4">📋</div>
          <h2 className="text-xl font-semibold text-gray-800">No Analysis Available</h2>
          <p className="text-gray-600 mt-2">Please try uploading your report again</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">

        {/* Header Section */}
        <div className="bg-white rounded-2xl shadow-xl p-8 mb-8 border border-gray-200">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between">
            <div className="mb-6 lg:mb-0">
              <h1 className="text-4xl font-bold text-gray-900 mb-4 flex items-center">
                <span className="text-blue-600 mr-3">🏥</span>
                Medical Report Analysis
              </h1>
              <div className="space-y-2 text-gray-600">
                <p><span className="font-semibold">File:</span> {report.original_filename}</p>
                <p><span className="font-semibold">Uploaded:</span> {formatDate(report.upload_date)}</p>
                {report.processed_at && (
                  <p><span className="font-semibold">Processed:</span> {formatDate(report.processed_at)}</p>
                )}
              </div>
            </div>

            {/* Risk Level Badge */}
            <div className={`inline-flex items-center px-6 py-3 rounded-full border-2 ${getRiskLevelColor(analysis.risk_level)}`}>
              <span className="font-bold text-lg">
                Risk Level: {analysis.risk_level.toUpperCase()}
              </span>
            </div>
          </div>
        </div>

        {/* Quick Summary */}
        <div className="bg-white rounded-2xl shadow-xl p-8 mb-8 border border-gray-200">
          <h2 className="text-3xl font-bold text-gray-900 mb-6 flex items-center">
            <span className="text-green-600 mr-3">📋</span>
            Quick Summary
          </h2>
          <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-6 border-l-4 border-blue-500">
            <p className="text-lg text-gray-800 leading-relaxed font-medium">
              {analysis.simple_summary}
            </p>
          </div>
        </div>

        {/* Health Metrics Dashboard */}
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-gray-900 mb-8 flex items-center">
            <span className="text-purple-600 mr-3">📊</span>
            Health Metrics Dashboard
          </h2>

          {analysis.health_metrics && analysis.health_metrics.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {analysis.health_metrics.map((metric, index) => (
                <SimpleSpeedometer key={index} metric={metric} />
              ))}
            </div>
          ) : (
            <div className="bg-white rounded-2xl shadow-xl p-12 text-center border border-gray-200">
              <div className="text-gray-400 text-6xl mb-4">📊</div>
              <h3 className="text-xl font-semibold text-gray-700 mb-2">No Specific Metrics</h3>
              <p className="text-gray-600">No specific health metrics were extracted from this report.</p>
            </div>
          )}
        </div>

        {/* Key Findings and Recommendations */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">

          {/* Key Findings */}
          <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-200">
            <h3 className="text-2xl font-bold text-gray-900 mb-6 flex items-center">
              <span className="text-orange-600 mr-3">🔍</span>
              Key Findings
            </h3>
            <div className="space-y-4">
              {analysis.key_findings && analysis.key_findings.length > 0 ? (
                analysis.key_findings.map((finding, index) => (
                  <div key={index} className="flex items-start bg-orange-50 rounded-lg p-4 border-l-4 border-orange-400">
                    <span className="text-orange-600 mr-3 mt-1 font-bold">•</span>
                    <span className="text-gray-800 font-medium">{finding}</span>
                  </div>
                ))
              ) : (
                <p className="text-gray-600 italic">No specific findings identified.</p>
              )}
            </div>
          </div>

          {/* Recommendations */}
          <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-200">
            <h3 className="text-2xl font-bold text-gray-900 mb-6 flex items-center">
              <span className="text-green-600 mr-3">💡</span>
              Recommendations
            </h3>
            <div className="space-y-4">
              {analysis.recommendations && analysis.recommendations.length > 0 ? (
                analysis.recommendations.map((recommendation, index) => (
                  <div key={index} className="flex items-start bg-green-50 rounded-lg p-4 border-l-4 border-green-400">
                    <span className="text-green-600 mr-3 mt-1 font-bold">✓</span>
                    <span className="text-gray-800 font-medium">{recommendation}</span>
                  </div>
                ))
              ) : (
                <p className="text-gray-600 italic">No specific recommendations available.</p>
              )}
            </div>
          </div>
        </div>

        {/* Detailed Medical Summary */}
        <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-200">
          <h2 className="text-3xl font-bold text-gray-900 mb-6 flex items-center">
            <span className="text-red-600 mr-3">🏥</span>
            Detailed Medical Summary
          </h2>
          <div className="bg-gray-50 rounded-xl p-6 border border-gray-200">
            <p className="text-gray-800 leading-relaxed whitespace-pre-wrap font-mono text-sm">
              {analysis.summary}
            </p>
          </div>
        </div>

        {/* Important Disclaimer */}
        <div className="mt-8 bg-yellow-50 border-2 border-yellow-200 rounded-2xl p-6 text-center">
          <div className="text-yellow-600 text-2xl mb-2">⚠️</div>
          <p className="text-yellow-800 font-semibold mb-2">Important Medical Disclaimer</p>
          <p className="text-yellow-700 text-sm">
            This analysis is generated by AI and should not replace professional medical advice.
            Always consult with your healthcare provider for medical decisions and treatment plans.
          </p>
        </div>
      </div>
    </div>
  );
};

export default NewMedicalAnalysis;