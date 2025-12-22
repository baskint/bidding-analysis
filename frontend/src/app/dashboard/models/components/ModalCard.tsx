'use client';

import type { MLModel } from '@/lib/types';

const MODEL_TYPES = [
  { value: 'bidding_optimizer', label: 'Bidding Optimizer' },
  { value: 'fraud_detector', label: 'Fraud Detector' },
  { value: 'conversion_predictor', label: 'Conversion Predictor' },
  { value: 'audience_segmentation', label: 'Audience Segmentation' },
];

interface ModelCardProps {
  model: MLModel;
  onEdit: () => void;
  onDelete: () => void;
  onSetDefault: () => void;
  getStatusColor: (status: string) => string;
}

export default function ModelCard({
  model,
  onEdit,
  onDelete,
  onSetDefault,
  getStatusColor,
}: ModelCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-md border border-gray-200 p-6 hover:shadow-lg transition-shadow">
      <div className="flex justify-between items-start mb-4">
        <div className="flex-1">
          <h3 className="text-xl font-semibold text-gray-900 mb-1">{model.name}</h3>
          <p className="text-sm text-gray-600">
            {MODEL_TYPES.find((t) => t.value === model.type)?.label || model.type}
          </p>
        </div>
        {model.is_default && (
          <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs font-semibold rounded">
            DEFAULT
          </span>
        )}
      </div>

      <div className="space-y-2 mb-4">
        <div className="flex items-center gap-2">
          <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusColor(model.status)}`}>
            {model.status.toUpperCase()}
          </span>
          <span className="text-sm text-gray-600">v{model.version}</span>
        </div>

        <p className="text-sm text-gray-700 line-clamp-2">{model.description || 'No description'}</p>

        <div className="flex items-center gap-2 text-sm text-gray-600">
          <span className="font-medium">Provider:</span>
          <span>{model.provider}</span>
        </div>
      </div>

      {/* Metrics */}
      {model.metrics && Object.keys(model.metrics).length > 0 && (
        <div className="mb-4 p-3 bg-gray-50 rounded-lg">
          <p className="text-xs font-semibold text-gray-700 mb-2">Performance Metrics</p>
          <div className="grid grid-cols-2 gap-2 text-xs">
            {Object.entries(model.metrics).slice(0, 4).map(([key, value]) => (
              <div key={key}>
                <span className="text-gray-600">{key}:</span>{' '}
                <span className="font-medium">
                  {typeof value === 'number' ? value.toFixed(3) : String(value)}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Actions */}
      <div className="flex gap-2">
        <button
          onClick={onEdit}
          className="flex-1 px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors text-sm font-medium"
        >
          Edit
        </button>
        {!model.is_default && (
          <button
            onClick={onSetDefault}
            className="flex-1 px-3 py-2 bg-blue-50 text-blue-700 rounded hover:bg-blue-100 transition-colors text-sm font-medium"
          >
            Set Default
          </button>
        )}
        <button
          onClick={onDelete}
          className="px-3 py-2 bg-red-50 text-red-700 rounded hover:bg-red-100 transition-colors text-sm font-medium"
        >
          Delete
        </button>
      </div>

      <p className="text-xs text-gray-500 mt-3">
        Created: {new Date(model.created_at).toLocaleDateString()}
      </p>
    </div>
  );
}
