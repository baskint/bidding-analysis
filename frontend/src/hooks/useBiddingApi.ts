// src/hooks/useBiddingApi.ts
import { useState, useEffect } from 'react';
import {
  getBidData,
  getCampaignData,
  getAnalytics,
  processBid,
  type BidData,
  type CampaignData,
  type AnalyticsData
} from '@/lib/api';

// Helper function to handle authentication errors
const handleAuthError = (error: Error) => {
  if (error.message.includes('401') || error.message.includes('Unauthorized')) {
    // Token expired or invalid - redirect to login
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
    window.location.href = '/login';
    return true;
  }
  return false;
};

// Hook for fetching bid data
export function useBidData() {
  const [data, setData] = useState<BidData[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchBidData = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const result = await getBidData();
      setData(result);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch bid data');

      // Handle authentication errors
      if (handleAuthError(error)) {
        return;
      }

      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  // Clear error
  const clearError = (): void => {
    setError(null);
  };

  return {
    data,
    loading,
    error,
    fetchBidData,
    refetch: fetchBidData,
    clearError,
  };
}

// Hook for fetching campaign data
export function useCampaignData() {
  const [data, setData] = useState<CampaignData[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchCampaignData = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const result = await getCampaignData();
      setData(result);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch campaign data');

      // Handle authentication errors
      if (handleAuthError(error)) {
        return;
      }

      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  // Clear error
  const clearError = (): void => {
    setError(null);
  };

  return {
    data,
    loading,
    error,
    fetchCampaignData,
    refetch: fetchCampaignData,
    clearError,
  };
}

// Hook for fetching analytics
export function useAnalytics() {
  const [data, setData] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAnalytics = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const result = await getAnalytics();
      setData(result);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch analytics');

      // Handle authentication errors
      if (handleAuthError(error)) {
        return;
      }

      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  // Auto-fetch on mount only if authenticated
  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      fetchAnalytics();
    }
  }, []);

  // Clear error
  const clearError = (): void => {
    setError(null);
  };

  return {
    data,
    loading,
    error,
    fetchAnalytics,
    refetch: fetchAnalytics,
    clearError,
  };
}

// Hook for processing bids
export function useBidProcessor() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const processBidRequest = async (bidData: Partial<BidData>): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await processBid(bidData);
      return true;
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to process bid');

      // Handle authentication errors
      if (handleAuthError(error)) {
        return false;
      }

      setError(error.message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  // Clear error
  const clearError = (): void => {
    setError(null);
  };

  return {
    loading,
    error,
    processBid: processBidRequest,
    clearError,
  };
}
