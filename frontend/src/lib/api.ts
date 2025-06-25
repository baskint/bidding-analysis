// src/lib/api.ts
const API_BASE_URL = process.env.NODE_ENV === 'production'
  ? 'https://bidding-analysis-539382269313.us-central1.run.app'
  : 'http://localhost:8080';

export interface User {
  id: string;
  username: string;
  created_at: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

// Helper function for making API requests
async function apiRequest<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;

  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: 'Network error' }));
    throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
  }

  return response.json();
}

// Auth API functions
export async function loginUser(username: string, password: string): Promise<AuthResponse> {
  return apiRequest<AuthResponse>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
}

export async function registerUser(username: string, password: string): Promise<AuthResponse> {
  return apiRequest<AuthResponse>('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
}

export async function getCurrentUser(): Promise<{ user: User }> {
  return apiRequest<{ user: User }>('/api/auth/me');
}

// Bidding analysis types
export interface BidData {
  id: string;
  campaign_id: string;
  amount: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CampaignData {
  id: string;
  name: string;
  budget: number;
  status: string;
  performance_metrics: {
    impressions: number;
    clicks: number;
    conversions: number;
    cost: number;
  };
  created_at: string;
}

export interface AnalyticsData {
  total_campaigns: number;
  total_bids: number;
  total_revenue: number;
  performance_summary: {
    click_through_rate: number;
    conversion_rate: number;
    cost_per_click: number;
    return_on_ad_spend: number;
  };
  time_period: {
    start_date: string;
    end_date: string;
  };
}

// Bidding analysis API functions
export async function getBidData(): Promise<BidData[]> {
  return apiRequest<BidData[]>('/api/bids');
}

export async function getCampaignData(): Promise<CampaignData[]> {
  return apiRequest<CampaignData[]>('/api/campaigns');
}

export async function getAnalytics(): Promise<AnalyticsData> {
  return apiRequest<AnalyticsData>('/api/analytics');
}
