import React, { useState, useRef } from 'react';
import { apiService } from '../services/apiService';

interface FileUploadProps {
  onUploadSuccess: (reportId: number) => void;
}

const FileUpload: React.FC<FileUploadProps> = ({ onUploadSuccess }) => {
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [dragActive, setDragActive] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = (selectedFile: File) => {
    // Validate file type
    const allowedTypes = ['text/plain', 'application/pdf', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'];
    if (!allowedTypes.includes(selectedFile.type)) {
      setError('Please upload a TXT, PDF, or DOCX file only.');
      return;
    }

    // Validate file size (20MB)
    const maxSize = 20 * 1024 * 1024; // 20MB in bytes
    if (selectedFile.size > maxSize) {
      setError('File size must be less than 20MB.');
      return;
    }

    setFile(selectedFile);
    setError(null);
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0];
    if (selectedFile) {
      handleFileSelect(selectedFile);
    }
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    const droppedFile = e.dataTransfer.files?.[0];
    if (droppedFile) {
      handleFileSelect(droppedFile);
    }
  };

  const handleUpload = async () => {
    if (!file) return;

    setUploading(true);
    setUploadProgress(0);
    setError(null);

    try {
      // Simulate upload progress
      const progressInterval = setInterval(() => {
        setUploadProgress(prev => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return 90;
          }
          return prev + 10;
        });
      }, 200);

      const description = `Medical report uploaded on ${new Date().toLocaleDateString()}`;
      const result = await apiService.uploadReport(file, description);

      clearInterval(progressInterval);
      setUploadProgress(100);

      // Wait a moment for visual feedback, then proceed
      setTimeout(() => {
        onUploadSuccess(result.report_id);
      }, 500);

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
      setUploading(false);
      setUploadProgress(0);
    }
  };

  const resetUpload = () => {
    setFile(null);
    setUploading(false);
    setUploadProgress(0);
    setError(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white rounded-xl shadow-lg p-8">
        <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
          Upload Medical Report
        </h2>

        {!file ? (
          <div
            className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
              dragActive
                ? 'border-blue-500 bg-blue-50'
                : 'border-gray-300 hover:border-blue-400 hover:bg-gray-50'
            }`}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
          >
            <div className="text-6xl mb-4">📄</div>
            <h3 className="text-lg font-semibold text-gray-700 mb-2">
              Drop your medical report here
            </h3>
            <p className="text-gray-500 mb-4">
              or click to browse files
            </p>
            <button
              onClick={() => fileInputRef.current?.click()}
              className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors font-medium"
            >
              Choose File
            </button>
            <input
              ref={fileInputRef}
              type="file"
              onChange={handleFileChange}
              accept=".txt,.pdf,.docx"
              className="hidden"
            />
            <p className="text-xs text-gray-400 mt-3">
              Supports: TXT, PDF, DOCX files up to 20MB
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {/* Selected File Info */}
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="text-2xl">📄</div>
                  <div>
                    <p className="font-medium text-gray-800">{file.name}</p>
                    <p className="text-sm text-gray-500">{formatFileSize(file.size)}</p>
                  </div>
                </div>
                <button
                  onClick={resetUpload}
                  disabled={uploading}
                  className="text-red-600 hover:text-red-800 disabled:opacity-50"
                >
                  ✕
                </button>
              </div>
            </div>

            {/* Upload Progress */}
            {uploading && (
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Uploading...</span>
                  <span className="text-gray-600">{uploadProgress}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${uploadProgress}%` }}
                  ></div>
                </div>
                <p className="text-sm text-gray-500 text-center">
                  🔬 AI will analyze your report after upload...
                </p>
              </div>
            )}

            {/* Error Message */}
            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <div className="flex items-center">
                  <span className="text-red-500 text-xl mr-2">⚠️</span>
                  <span className="text-red-700">{error}</span>
                </div>
              </div>
            )}

            {/* Upload Button */}
            {!uploading && (
              <button
                onClick={handleUpload}
                className="w-full bg-green-600 text-white py-3 px-6 rounded-lg hover:bg-green-700 transition-colors font-medium text-lg"
              >
                Upload & Analyze Report
              </button>
            )}
          </div>
        )}

        {/* Information */}
        <div className="mt-6 p-4 bg-blue-50 rounded-lg">
          <h4 className="font-medium text-blue-900 mb-2">What happens next?</h4>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>• Your report will be securely uploaded</li>
            <li>• AI will analyze the medical data (~15 seconds)</li>
            <li>• You'll see health metrics with speedometer displays</li>
            <li>• Get easy-to-understand summaries and recommendations</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default FileUpload;