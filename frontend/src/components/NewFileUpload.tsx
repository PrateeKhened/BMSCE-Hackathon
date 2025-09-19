import React, { useState } from 'react';
import { newApiService } from '../services/newApiService';

interface NewFileUploadProps {
  onUploadSuccess: (reportId: number) => void;
}

const NewFileUpload: React.FC<NewFileUploadProps> = ({ onUploadSuccess }) => {
  const [file, setFile] = useState<File | null>(null);
  const [description, setDescription] = useState('');
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [dragActive, setDragActive] = useState(false);

  const handleFileSelect = (selectedFile: File) => {
    const allowedTypes = ['application/pdf', 'text/plain', 'application/msword',
                         'application/vnd.openxmlformats-officedocument.wordprocessingml.document'];

    if (!allowedTypes.includes(selectedFile.type)) {
      setError('Please upload a PDF, TXT, DOC, or DOCX file');
      return;
    }

    if (selectedFile.size > 10 * 1024 * 1024) { // 10MB limit
      setError('File size must be less than 10MB');
      return;
    }

    setFile(selectedFile);
    setError(null);
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

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFileSelect(e.dataTransfer.files[0]);
    }
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleFileSelect(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      setError('Please select a file to upload');
      return;
    }

    try {
      setUploading(true);
      setError(null);
      console.log('🚀 Starting upload for file:', file.name);

      const result = await newApiService.uploadReport(file, description || 'Medical report upload');
      console.log('✅ Upload successful:', result);

      // Call the success callback with the report ID
      onUploadSuccess(result.report_id);

    } catch (err) {
      console.error('❌ Upload failed:', err);
      setError(err instanceof Error ? err.message : 'Upload failed');
    } finally {
      setUploading(false);
    }
  };

  const resetForm = () => {
    setFile(null);
    setDescription('');
    setError(null);
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-200">
        <h2 className="text-3xl font-bold text-gray-900 mb-6 text-center flex items-center justify-center">
          <span className="text-blue-600 mr-3">📄</span>
          Upload Medical Report
        </h2>

        {/* File Upload Area */}
        <div
          className={`relative border-2 border-dashed rounded-xl p-8 text-center transition-all duration-200 ${
            dragActive
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
          }`}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
        >
          <input
            type="file"
            id="file-upload"
            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            accept=".pdf,.txt,.doc,.docx"
            onChange={handleFileInput}
            disabled={uploading}
          />

          {file ? (
            <div className="space-y-4">
              <div className="text-green-600 text-5xl mb-3">✅</div>
              <h3 className="text-lg font-semibold text-gray-900">File Selected</h3>
              <div className="bg-gray-50 rounded-lg p-4 border">
                <p className="font-medium text-gray-800">{file.name}</p>
                <p className="text-sm text-gray-600">
                  {(file.size / 1024 / 1024).toFixed(2)} MB • {file.type}
                </p>
              </div>
              <button
                onClick={resetForm}
                className="text-blue-600 hover:text-blue-700 font-medium text-sm"
                disabled={uploading}
              >
                Choose a different file
              </button>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="text-gray-400 text-5xl mb-3">📁</div>
              <h3 className="text-lg font-semibold text-gray-900">
                Drag and drop your medical report here
              </h3>
              <p className="text-gray-600">
                or <span className="text-blue-600 font-medium">click to browse</span>
              </p>
              <p className="text-sm text-gray-500">
                Supports PDF, TXT, DOC, DOCX (max 10MB)
              </p>
            </div>
          )}
        </div>

        {/* Description Field */}
        <div className="mt-6">
          <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
            Description (Optional)
          </label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Add any additional notes about this report..."
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
            rows={3}
            disabled={uploading}
          />
        </div>

        {/* Error Message */}
        {error && (
          <div className="mt-4 bg-red-50 border border-red-200 rounded-lg p-4 flex items-start">
            <span className="text-red-500 mr-2 mt-0.5">⚠️</span>
            <p className="text-red-700 font-medium">{error}</p>
          </div>
        )}

        {/* Upload Button */}
        <div className="mt-8">
          <button
            onClick={handleUpload}
            disabled={!file || uploading}
            className={`w-full py-4 px-6 rounded-lg font-semibold text-white transition-all duration-200 ${
              !file || uploading
                ? 'bg-gray-300 cursor-not-allowed'
                : 'bg-blue-600 hover:bg-blue-700 active:bg-blue-800 shadow-lg hover:shadow-xl'
            }`}
          >
            {uploading ? (
              <div className="flex items-center justify-center">
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                Processing Report...
              </div>
            ) : (
              <div className="flex items-center justify-center">
                <span className="mr-2">🚀</span>
                Upload & Analyze Report
              </div>
            )}
          </button>
        </div>

        {/* Help Text */}
        <div className="mt-6 text-center">
          <p className="text-sm text-gray-500">
            Your medical report will be analyzed using AI to extract key health metrics and insights.
          </p>
        </div>
      </div>
    </div>
  );
};

export default NewFileUpload;