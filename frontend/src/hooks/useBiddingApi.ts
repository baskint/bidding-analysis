// src/hooks/useBiddingApi.ts
import { useState, useEffect } from 'react';
import { getBidData, getCampaignData, getAnalytics, type BidData, type CampaignData, type AnalyticsData } from '@/lib/api';

// Hook for fetching bid data
export function useBidData() {
  const [data, setData] = useState<BidData[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchBidData = async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await getBidData();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch bid data');
    } finally {
      setLoading(false);
    }
  };

  return {
    data,
    loading,
    error,
    fetchBidData,
    refetch: fetchBidData,
  };
}

// Hook for fetching campaign data
export function useCampaignData() {
  const [data, setData] = useState<CampaignData[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchCampaignData = async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await getCampaignData();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch campaign data');
    } finally {
      setLoading(false);
    }
  };

  return {
    data,
    loading,
    error,
    fetchCampaignData,
    refetch: fetchCampaignData,
  };
}

// Hook for fetching analytics
export function useAnalytics() {
  const [data, setData] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAnalytics = async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await getAnalytics();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch analytics');
    } finally {
      setLoading(false);
    }
  };

  // Auto-fetch on mount
  useEffect(() => {
    fetchAnalytics();
  }, []);

  return {
    data,
    loading,
    error,
    fetchAnalytics,
    refetch: fetchAnalytics,
  };
}
