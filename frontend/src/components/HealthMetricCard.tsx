import React from 'react';
import { HealthMetric } from '../services/apiService';

interface HealthMetricCardProps {
  metric: HealthMetric;
}

const HealthMetricCard: React.FC<HealthMetricCardProps> = ({ metric }) => {
  const getColorByScore = (score: number): string => {
    if (score >= 80) return '#10B981'; // Green - Normal
    if (score >= 50) return '#F59E0B'; // Orange - Warning
    return '#EF4444'; // Red - Critical
  };

  const getColorByStatus = (status: string): string => {
    switch (status) {
      case 'normal': return '#10B981';
      case 'warning': return '#F59E0B';
      case 'critical': return '#EF4444';
      default: return '#6B7280';
    }
  };

  const circumference = 2 * Math.PI * 45; // radius = 45
  const strokeDasharray = circumference;
  const strokeDashoffset = circumference - (metric.score / 100) * circumference;

  return (
    <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition-shadow duration-300">
      {/* Header */}
      <div className="text-center mb-4">
        <h3 className="text-lg font-semibold text-gray-800 mb-1">{metric.name}</h3>
        <div className="text-2xl font-bold text-gray-900">
          {metric.value} <span className="text-sm font-normal text-gray-500">{metric.unit}</span>
        </div>
      </div>

      {/* Speedometer */}
      <div className="relative flex justify-center mb-4">
        <svg width="120" height="120" className="transform -rotate-90">
          {/* Background circle */}
          <circle
            cx="60"
            cy="60"
            r="45"
            stroke="#E5E7EB"
            strokeWidth="8"
            fill="none"
          />
          {/* Progress circle */}
          <circle
            cx="60"
            cy="60"
            r="45"
            stroke={getColorByScore(metric.score)}
            strokeWidth="8"
            fill="none"
            strokeDasharray={strokeDasharray}
            strokeDashoffset={strokeDashoffset}
            strokeLinecap="round"
            className="transition-all duration-1000 ease-out"
          />
        </svg>

        {/* Score overlay */}
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className="text-2xl font-bold text-gray-800">{metric.score}</span>
          <span className="text-xs text-gray-500 font-medium">/ 100</span>
        </div>
      </div>

      {/* Status Badge */}
      <div className="flex justify-center mb-3">
        <span
          className={`px-3 py-1 rounded-full text-xs font-semibold text-white`}
          style={{ backgroundColor: getColorByStatus(metric.status) }}
        >
          {metric.status.toUpperCase()}
        </span>
      </div>

      {/* Description */}
      <p className="text-sm text-gray-600 text-center mb-3 leading-relaxed">
        {metric.description}
      </p>

      {/* Normal Range */}
      <div className="bg-gray-50 rounded-lg p-3 text-center">
        <p className="text-xs text-gray-500 mb-1">Normal Range</p>
        <p className="text-sm font-medium text-gray-700">
          {metric.range_min} - {metric.range_max} {metric.unit}
        </p>
      </div>
    </div>
  );
};

export default HealthMetricCard;