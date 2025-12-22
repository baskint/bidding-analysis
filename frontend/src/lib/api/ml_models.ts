// frontend/src/lib/api/ml_models.ts
/**
 * ML Model management API functions
 */

import { apiGet, apiPost } from '@/lib/utils';
import type { MLModel } from '@/lib/types';

// ML Model specific types
export interface MLModelCreate {
  name: string;
  type: string;
  version: string;
  description?: string;
  provider: string;
  endpoint?: string;
  api_key?: string;
  config?: Record<string, unknown>;
  is_default?: boolean;
status?: string;
}

export interface MLModelUpdate {
  name?: string;
  version?: string;
  description?: string;
  status?: string;
  provider?: string;
  endpoint?: string;
  api_key?: string;
  config?: Record<string, unknown>;
  metrics?: Record<string, unknown>;
  is_default?: boolean;
}

export interface MLModelListResponse {
  models: MLModel[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// ML Model API functions

/**
 * List all ML models with pagination
 */
export async function listMLModels(page = 1, pageSize = 20): Promise<MLModelListResponse> {
  return apiGet<MLModelListResponse>(`/trpc/mlModel.list?page=${page}&pageSize=${pageSize}`);
}

/**
 * Get a single ML model by ID
 */
export async function getMLModel(id: string): Promise<MLModel> {
  return apiGet<MLModel>(`/trpc/mlModel.get?id=${id}`);
}

/**
 * Create a new ML model
 */
export async function createMLModel(data: MLModelCreate): Promise<MLModel> {
  return apiPost<MLModel>('/trpc/mlModel.create', data);
}

/**
 * Update an existing ML model
 */
export async function updateMLModel(id: string, data: MLModelUpdate): Promise<MLModel> {
  return apiPost<MLModel>('/trpc/mlModel.update', { id, ...data });
}

/**
 * Delete an ML model
 */
export async function deleteMLModel(id: string): Promise<{ success: boolean; message: string }> {
  return apiPost<{ success: boolean; message: string }>('/trpc/mlModel.delete', { id });
}

/**
 * Set a model as the default for its type
 */
export async function setDefaultMLModel(id: string): Promise<MLModel> {
  return apiPost<MLModel>('/trpc/mlModel.setDefault', { id });
}

/**
 * Get the default ML model for a specific type
 */
export async function getDefaultMLModel(type: string): Promise<MLModel> {
  return apiGet<MLModel>(`/trpc/mlModel.getDefault?type=${type}`);
}
