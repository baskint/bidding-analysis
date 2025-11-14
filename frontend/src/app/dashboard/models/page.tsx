'use client';

import { useState, useEffect, useCallback } from 'react';
import {
  listMLModels,
  createMLModel,
  updateMLModel,
  deleteMLModel,
  setDefaultMLModel,
  type MLModel,
  type MLModelCreate,
  type MLModelUpdate,
} from '@/lib/api/ml_models';
import ModelCard from './components/ModalCard';
import ModelFormModal from './components/ModelFormModal';

const MODEL_STATUSES = [
  { value: 'active', label: 'Active', color: 'bg-green-100 text-green-800' },
  { value: 'inactive', label: 'Inactive', color: 'bg-gray-100 text-gray-800' },
  { value: 'training', label: 'Training', color: 'bg-blue-100 text-blue-800' },
  { value: 'testing', label: 'Testing', color: 'bg-yellow-100 text-yellow-800' },
];

export default function ModelsPage() {
  const [models, setModels] = useState<MLModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedModel, setSelectedModel] = useState<MLModel | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  // Wrap loadModels in useCallback to prevent infinite re-renders
  const loadModels = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await listMLModels(page, 20);

      // Your API returns the data directly, not wrapped in result
      setModels(response.models || []);
      setTotalPages(response.total_pages || 1);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load models');
      console.error('Failed to load models:', err);
    } finally {
      setLoading(false);
    }
  }, [page]); // Add page as dependency since it's used inside

  useEffect(() => {
    loadModels();
  }, [loadModels]); // Now loadModels is stable due to useCallback

  const handleCreate = async (data: MLModelCreate | MLModelUpdate): Promise<void> => {
    try {
      // Since this is for creation, we need to ensure required fields are present
      // Cast to MLModelCreate and provide defaults for required fields
      const createData: MLModelCreate = {
        name: data.name || '', // Provide default if undefined
        type: (data as MLModelCreate).type || 'bidding_optimizer', // Cast and provide default
        version: data.version || '1.0.0', // Provide default
        provider: (data as MLModelCreate).provider || 'openai', // Cast and provide default
        description: data.description,
        endpoint: data.endpoint,
        api_key: data.api_key,
        config: data.config,
        is_default: data.is_default,
        status: data.status,
      };

      await createMLModel(createData);
      setShowCreateModal(false);
      loadModels();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to create model');
    }
  };

  const handleUpdate = async (id: string, data: MLModelUpdate) => {
    try {
      await updateMLModel(id, data);
      setShowEditModal(false);
      setSelectedModel(null);
      loadModels();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update model');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this model?')) return;

    try {
      await deleteMLModel(id);
      loadModels();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete model');
    }
  };

  const handleSetDefault = async (id: string) => {
    try {
      await setDefaultMLModel(id);
      loadModels();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to set default model');
    }
  };

  const getStatusColor = (status: string) => {
    const statusObj = MODEL_STATUSES.find((s) => s.value === status);
    return statusObj?.color || 'bg-gray-100 text-gray-800';
  };

  if (loading && models.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600">Loading models...</div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-slate-100">ML Models</h1>
          <p className="text-gray-600 mt-1">
            Configure and manage your machine learning models for bidding optimization
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          + Create Model
        </button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      {/* Models Grid */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {models.map((model) => (
          <ModelCard
            key={model.id}
            model={model}
            onEdit={() => {
              setSelectedModel(model);
              setShowEditModal(true);
            }}
            onDelete={() => handleDelete(model.id)}
            onSetDefault={() => handleSetDefault(model.id)}
            getStatusColor={getStatusColor}
          />
        ))}
      </div>

      {models.length === 0 && !loading && (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg mb-4">No models configured yet</p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            Create Your First Model
          </button>
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center gap-2 mt-6">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-4 py-2 border rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
          >
            Previous
          </button>
          <span className="px-4 py-2">
            Page {page} of {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className="px-4 py-2 border rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
          >
            Next
          </button>
        </div>
      )}

      {/* Create Modal */}
      {showCreateModal && (
        <ModelFormModal
          title="Create New Model"
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreate}
        />
      )}

      {/* Edit Modal */}
      {showEditModal && selectedModel && (
        <ModelFormModal
          title="Edit Model"
          model={selectedModel}
          onClose={() => {
            setShowEditModal(false);
            setSelectedModel(null);
          }}
          onSubmit={(data) => handleUpdate(selectedModel.id, data)}
        />
      )}
    </div>
  );
}
