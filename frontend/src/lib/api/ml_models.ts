// src/lib/api/ml_models.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
import { getAuthHeaders, handleResponse } from '../api';

// Types
export interface MLModel {
  id: string;
  user_id: string;
  name: string;
  type: string;  // 'bidding_optimizer', 'fraud_detector', 'conversion_predictor'
  version: string;
  description: string;
  status: string;  // 'active', 'inactive', 'training', 'testing'
  provider: string;  // 'openai', 'custom', 'tensorflow', 'pytorch'
  endpoint?: string;
  config: Record<string, unknown>;
  metrics: Record<string, unknown>;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

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
export async function listMLModels(page = 1, pageSize = 20): Promise<MLModelListResponse> {
  const response = await fetch(`${API_BASE_URL}/trpc/mlModel.list?page=${page}&pageSize=${pageSize}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  const result = await handleResponse(response);
  
  // Extract data from tRPC response format
  return result.result?.data || result;
}

export async function getMLModel(id: string): Promise<MLModel> {
  const response = await fetch(`${API_BASE_URL}/trpc/mlModel.get?id=${id}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  return handleResponse(response);
}

export async function createMLModel(data: MLModelCreate): Promise<MLModel> {
  const response = await fetch(`${API_BASE_URL}/trpc/mlModel.create`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });

  return handleResponse(response);
}

export async function updateMLModel(id: string, data: MLModelUpdate): Promise<MLModel> {
  const response = await fetch(`${API_BASE_URL}/trpc/mlModel.update`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({ id, ...data }),
  });

  return handleResponse(response);
}

export async function deleteMLModel(id: string): Promise<{ success: boolean; message: string }> {
  const response = await fetch(`${API_BASE_URL}/trpc/mlModel.delete`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({ id }),
  });

  return handleResponse(response);
}

export async function setDefaultMLModel(id: string): Promise<MLModel> {
  const response = await fetch(`${API_BASE_URL}/trpc/mlModel.setDefault`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({ id }),
  });

  return handleResponse(response);
}

export async function getDefaultMLModel(type: string): Promise<MLModel> {
  const response = await fetch(`${API_BASE_URL}/trpc/mlModel.getDefault?type=${type}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  return handleResponse(response);
}
