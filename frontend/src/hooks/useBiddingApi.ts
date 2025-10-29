import { useState, useCallback, useEffect } from 'react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Types matching your backend
interface BidSubmitRequest {
  campaign_id: string;
  user_id: string;
  bid_price: number;
  floor_price: number;
  device_type: string;
  os?: string;
  browser?: string;
  country: string;
  region?: string;
  city?: string;
  keywords?: string[];
  segment_id?: string;
  segment_category?: string;
}

interface BidSubmitResponse {
  bid_event_id: string;
  status: 'won' | 'lost';
  processed_at: string;
  win_probability: number;
  message: string;
}

interface BidPredictionRequest {
  campaign_id: string;
  user_segment: string;
  device_type: string;
  country: string;
  floor_price: number;
  keywords?: string[];
  engagement_score?: number;
  conversion_probability?: number;
}

interface BidPredictionResponse {
  predicted_bid: number;
  confidence: number;
  strategy: string;
  fraud_risk: boolean;
  reasoning: string;
}

interface BidStreamData {
  bid_event_id: string;
  campaign_id: string;
  bid_price: number;
  win_price?: number;
  won: boolean;
  converted: boolean;
  timestamp: string;
  segment_category: string;
  device_type: string;
  country: string;
}

interface CampaignTargeting {
  countries?: string[];
  segment_categories?: string[];
  device_types?: string[];
}

interface CampaignCreateRequest {
  name: string;
  budget: number;
  start_date?: string; // ISO date
  end_date?: string; // ISO date
  targeting?: CampaignTargeting;
  bidding_strategy?: 'manual' | 'auto' | 'target_cpa';
  bid_amount?: number;
  metadata?: Record<string, unknown>;
}

// interface CampaignCreateResponse {
//   campaign_id: string;
//   name: string;
//   status: 'active' | 'paused' | 'archived';
//   created_at: string;
//   budget: number;
//   targeting?: CampaignTargeting;
//   bidding_strategy?: string;
// }

// Helper function to get auth token
function getAuthToken(): string | null {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('auth_token');
  }
  return null;
}

// Custom hook for bidding operations
export function useBiddingApi() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const makeAuthenticatedRequest = async (url: string, options: RequestInit = {}) => {
    const token = getAuthToken();

    const headers = {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    };

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Authentication required. Please log in.');
      }
      throw new Error(`Request failed: ${response.statusText}`);
    }

    return response.json();
  };

  const submitBid = useCallback(async (bidData: BidSubmitRequest): Promise<BidSubmitResponse | null> => {
    setLoading(true);
    setError(null);

    try {
      const result = await makeAuthenticatedRequest(`${API_BASE_URL}/trpc/bidding.submit`, {
        method: 'POST',
        body: JSON.stringify(bidData),
      });

      // Handle tRPC response format
      return result.result?.data || result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to submit bid';
      setError(errorMessage);
      console.error('Bid submission error:', err);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const predictBid = useCallback(async (predictionData: BidPredictionRequest): Promise<BidPredictionResponse | null> => {
    setLoading(true);
    setError(null);

    try {
      const result = await makeAuthenticatedRequest(`${API_BASE_URL}/trpc/bidding.predict`, {
        method: 'POST',
        body: JSON.stringify(predictionData),
      });

      return result.result?.data || result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to predict bid';
      setError(errorMessage);
      console.error('Bid prediction error:', err);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const getBidStream = useCallback(async (): Promise<BidStreamData[]> => {
    try {
      const result = await makeAuthenticatedRequest(`${API_BASE_URL}/trpc/bidding.stream`);
      return result.result?.data || result || [];
    } catch (err) {
      console.error('Bid stream error:', err);
      return [];
    }
  }, []);

  const getDashboardMetrics = useCallback(async () => {
    try {
      const result = await makeAuthenticatedRequest(`${API_BASE_URL}/trpc/analytics.getDashboardMetrics`);
      return result.result?.data || result;
    } catch (err) {
      console.error('Dashboard metrics error:', err);
      return null;
    }
  }, []);

  const createCampaign = useCallback(async (campaignData: CampaignCreateRequest) => {
    try {
      const result = await makeAuthenticatedRequest(`${API_BASE_URL}/trpc/campaign.create`, {
        method: 'POST',
        body: JSON.stringify(campaignData),
      });
      return result.result?.data || result;
    } catch (err) {
      console.error('Create campaign error:', err);
      throw err;
    }
  }, []);

  const listCampaigns = useCallback(async () => {
    try {
      const result = await makeAuthenticatedRequest(`${API_BASE_URL}/trpc/campaign.list`);
      return result.result?.data || result || [];
    } catch (err) {
      console.error('List campaigns error:', err);
      return [];
    }
  }, []);

  return {
    submitBid,
    predictBid,
    getBidStream,
    getDashboardMetrics,
    createCampaign,
    listCampaigns,
    loading,
    error,
    clearError: () => setError(null),
  };
}

// Hook for real-time bid stream
export function useBidStream(intervalMs: number = 5000) {
  const [bidStream, setBidStream] = useState<BidStreamData[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const { getBidStream } = useBiddingApi();

  const startStream = useCallback(() => {
    setIsStreaming(true);
  }, []);

  const stopStream = useCallback(() => {
    setIsStreaming(false);
  }, []);

  useEffect(() => {
    if (!isStreaming) return;

    const fetchBids = async () => {
      const bids = await getBidStream();
      setBidStream(bids);
    };

    // Initial fetch
    fetchBids();

    // Set up interval
    const interval = setInterval(fetchBids, intervalMs);

    return () => clearInterval(interval);
  }, [isStreaming, getBidStream, intervalMs]);

  return {
    bidStream,
    isStreaming,
    startStream,
    stopStream,
  };
}

// Hook for bid form management
export function useBidForm() {
  const [formData, setFormData] = useState<Partial<BidSubmitRequest>>({
    device_type: 'mobile',
    country: 'US',
    segment_category: 'standard',
  });

  const updateField = useCallback(<K extends keyof BidSubmitRequest>(
    field: K,
    value: BidSubmitRequest[K]
  ) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  }, []);

  const resetForm = useCallback(() => {
    setFormData({
      device_type: 'mobile',
      country: 'US',
      segment_category: 'standard',
    });
  }, []);

  return {
    formData,
    updateField,
    resetForm,
    isValid: formData.campaign_id && formData.bid_price && formData.floor_price,
  };
}
