import React from 'react';
import { HealthMetric } from '../services/newApiService';

interface SimpleSpeedometerProps {
  metric: HealthMetric;
}

const SimpleSpeedometer: React.FC<SimpleSpeedometerProps> = ({ metric }) => {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'normal': return '#10b981'; // green
      case 'warning': return '#f59e0b'; // orange
      case 'critical': return '#ef4444'; // red
      default: return '#6b7280'; // gray
    }
  };

  const color = getStatusColor(metric.status);
  const radius = 45;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (metric.score / 100) * circumference;

  return (
    <div className="bg-white rounded-xl shadow-lg p-6 text-center">
      {/* Speedometer Circle */}
      <div className="relative w-32 h-32 mx-auto mb-4">
        <svg width="128" height="128" className="transform -rotate-90">
          {/* Background Circle */}
          <circle
            cx="64"
            cy="64"
            r={radius}
            stroke="#e5e7eb"
            strokeWidth="8"
            fill="none"
          />
          {/* Progress Circle */}
          <circle
            cx="64"
            cy="64"
            r={radius}
            stroke={color}
            strokeWidth="8"
            fill="none"
            strokeDasharray={circumference}
            strokeDashoffset={strokeDashoffset}
            strokeLinecap="round"
            className="transition-all duration-1000 ease-out"
          />
        </svg>
        {/* Score Text */}
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="text-center">
            <div className="text-2xl font-bold" style={{ color }}>
              {metric.score}
            </div>
            <div className="text-xs text-gray-500">Score</div>
          </div>
        </div>
      </div>

      {/* Metric Info */}
      <h3 className="text-lg font-semibold text-gray-900 mb-2">
        {metric.name}
      </h3>

      <div className="text-xl font-bold text-gray-700 mb-1">
        {metric.value} {metric.unit}
      </div>

      <div
        className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium mb-3`}
        style={{
          backgroundColor: `${color}20`,
          color: color,
          border: `1px solid ${color}40`
        }}
      >
        {metric.status.toUpperCase()}
      </div>

      <p className="text-sm text-gray-600 leading-relaxed">
        {metric.description}
      </p>

      {/* Range Info */}
      {metric.range_min !== undefined && metric.range_max !== undefined && (
        <div className="mt-3 text-xs text-gray-500">
          Normal: {metric.range_min} - {metric.range_max} {metric.unit}
        </div>
      )}
    </div>
  );
};

export default SimpleSpeedometer;