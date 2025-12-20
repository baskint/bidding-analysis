// frontend/src/lib/api/analytics.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Helper function to get auth headers
const getAuthHeaders = (): Record<string, string> => {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("auth_token") : null;
  return {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
  };
};

// Helper function to handle API responses
const handleResponse = async (response: Response) => {
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `API call failed: ${response.status}`);
  }
  const data = await response.json();
  // tRPC wraps responses in result.data
  return data.result?.data || data;
};

// Types
export interface PerformanceMetrics {
  total_bids: number;
  won_bids: number;
  conversions: number;
  total_spend: number;
  revenue: number;
  win_rate: number;
  conversion_rate: number;
  average_bid: number;
  cpa: number; // Cost Per Acquisition
  roas: number; // Return On Ad Spend
  fraud_detections: number;
}

export interface KeywordAnalysis {
  keyword: string;
  totalBids: number;
  wonBids: number;
  conversions: number;
  spend: number;
  revenue: number;
  winRate: number;
  conversionRate: number;
  cpa: number;
  roas: number;
}

export interface DeviceBreakdown {
  deviceType: string;
  totalBids: number;
  wonBids: number;
  conversions: number;
  spend: number;
  winRate: number;
  conversionRate: number;
  averageBid: number;
}

export interface GeoBreakdown {
  country: string;
  region: string;
  totalBids: number;
  wonBids: number;
  conversions: number;
  spend: number;
  winRate: number;
  conversionRate: number;
}

export interface HourlyPerformance {
  hour: number;
  total_bids: number;
  won_bids: number;
  conversions: number;
  spend: number;
  win_rate: number;
  conversion_rate: number;
  average_bid: number;
}

export interface DailyTrend {
  date: string;
  total_bids: number;
  won_bids: number;
  conversions: number;
  spend: number;
  revenue: number;
  win_rate: number;
  conversion_rate: number;
  cpa: number;
}

export interface CompetitiveAnalysis {
  segment_category: string;
  our_win_rate: number;
  market_average_bid: number;
  our_average_bid: number;
  average_floor_price: number;
  competition_intensity: number;
  total_opportunities: number;
}

export interface CampaignComparison {
  campaign_id: string;
  campaign_name: string;
  total_bids: number;
  won_bids: number;
  conversions: number;
  spend: number;
  win_rate: number;
  conversion_rate: number;
}

export interface DateRangeParams {
  start_date?: string; // YYYY-MM-DD format
  end_date?: string; // YYYY-MM-DD format
}

// API Functions
/**
 * Get overall performance metrics
 */
export async function getPerformanceOverview(
  params?: DateRangeParams,
): Promise<PerformanceMetrics> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getPerformanceOverview`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
      cache: 'no-store', // Important: prevent Next.js caching
      signal: AbortSignal.timeout(30000), // 30 second timeout
    },
  );
  return handleResponse(response);
}

/**
 * Get keyword performance analysis
 */
export async function getKeywordAnalysis(
  params?: DateRangeParams & { limit?: number },
): Promise<KeywordAnalysis[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getKeywordAnalysis`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
      cache: 'no-store', // Important: prevent Next.js caching
      signal: AbortSignal.timeout(30000), // 30 second timeout
    },
  );
  return handleResponse(response);
}

/**
 * Get device breakdown
 */
export async function getDeviceBreakdown(
  params?: DateRangeParams,
): Promise<DeviceBreakdown[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getDeviceBreakdown`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
      cache: 'no-store', // Important: prevent Next.js caching
      signal: AbortSignal.timeout(30000), // 30 second timeout
    },
  );
  return handleResponse(response);
}

/**
 * Get geographic breakdown
 */
export async function getGeoBreakdown(
  params?: DateRangeParams & { limit?: number },
): Promise<GeoBreakdown[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getGeoBreakdown`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
    },
  );
  return handleResponse(response);
}

/**
 * Get hourly performance
 */
export async function getHourlyPerformance(
  params?: DateRangeParams,
): Promise<HourlyPerformance[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getHourlyPerformance`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
    },
  );
  return handleResponse(response);
}

/**
 * Get daily performance trends
 */
export async function getDailyTrends(
  params?: DateRangeParams,
): Promise<DailyTrend[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getDailyTrends`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
    },
  );
  return handleResponse(response);
}

/**
 * Get competitive analysis
 */
export async function getCompetitiveAnalysis(
  params?: DateRangeParams,
): Promise<CompetitiveAnalysis[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getCompetitiveAnalysis`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
    },
  );
  return handleResponse(response);
}

/**
 * Get campaign comparison
 */
export async function getCampaignComparison(
  params?: DateRangeParams & { campaign_ids?: string[] },
): Promise<CampaignComparison[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/analytics.getCampaignComparison`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
    },
  );
  return handleResponse(response);
}
