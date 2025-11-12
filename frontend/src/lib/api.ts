// src/lib/api.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Types
export interface User {
  id: string;
  username: string;
  created_at: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface BidData {
  id?: string;
  campaignId: string;
  userId: string;
  bidPrice: number;
  winPrice?: number;
  floorPrice?: number;
  won?: boolean;
  converted?: boolean;
  timestamp?: string;
  [key: string]: unknown;
}

export interface CampaignData {
  id: string;
  name: string;
  status: string;
  budget: number;
  spent: number;
  impressions: number;
  clicks: number;
  conversions: number;
  [key: string]: unknown;
}

export interface AnalyticsData {
  totalBids: number;
  winRate: number;
  conversionRate: number;
  averageBidPrice: number;
  fraudAlerts: number;
  modelAccuracy: number;
  [key: string]: unknown;
}

// Helper function to get auth headers
export const getAuthHeaders = (): Record<string, string> => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

// Helper function to handle API responses
export const handleResponse = async (response: Response) => {
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `API call failed: ${response.status}`);
  }
  return response.json();
};

// Auth API functions
export async function loginUser(username: string, password: string): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/trpc/auth.login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });

  return handleResponse(response);
}

export async function registerUser(username: string, password: string): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/trpc/auth.register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });

  return handleResponse(response);
}

export async function getCurrentUser(): Promise<{ user: User }> {
  const response = await fetch(`${API_BASE_URL}/trpc/auth.me`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({}),
  });

  return handleResponse(response);
}

// Bidding API functions
export async function getBidData(): Promise<BidData[]> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.getBidHistory`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({}),
  });

  return handleResponse(response);
}

export async function getCampaignData(): Promise<CampaignData[]> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.getStats`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({}),
  });

  return handleResponse(response);
}

export async function getAnalytics(): Promise<AnalyticsData> {
  const response = await fetch(`${API_BASE_URL}/trpc/analytics.getModelAccuracy`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({}),
  });

  return handleResponse(response);
}

export async function getFraudAlerts(): Promise<unknown> {
  const response = await fetch(`${API_BASE_URL}/trpc/analytics.getFraudAlerts`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({}),
  });

  return handleResponse(response);
}

export async function processBid(bidData: Partial<BidData>): Promise<unknown> {
  const response = await fetch(`${API_BASE_URL}/trpc/bidding.processBid`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(bidData),
  });

  return handleResponse(response);
}
