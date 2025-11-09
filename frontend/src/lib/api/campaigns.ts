// frontend/src/lib/api/campaigns.ts
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
export interface Campaign {
  id: string;
  name: string;
  user_id: string;
  status: "active" | "paused" | "archived";
  budget?: number;
  daily_budget?: number;
  target_cpa?: number;
  created_at: string;
  updated_at: string;
}

export interface CampaignSummary extends Campaign {
  total_bids: number;
  won_bids: number;
  conversions: number;
  total_spend: number;
  win_rate: number;
  conversion_rate: number;
  average_bid: number;
  cost_per_conversion: number;
  last_activity_at?: string;
}

export interface DailyMetric {
  date: string;
  total_bids: number;
  won_bids: number;
  conversions: number;
  total_spend: number;
  win_rate: number;
  conversion_rate: number;
}

export interface KeywordStat {
  keyword: string;
  bids: number;
  won_bids: number;
  conversions: number;
  spend: number;
  win_rate: number;
}

export interface DeviceStat {
  device_type: string;
  bids: number;
  won_bids: number;
  conversions: number;
  win_rate: number;
}

export interface GeoStat {
  country: string;
  bids: number;
  won_bids: number;
  conversions: number;
  win_rate: number;
}

export interface BidEvent {
  id: string;
  campaign_id: string;
  user_id: string;
  bid_price: number;
  win_price?: number;
  floor_price?: number;
  won: boolean;
  converted: boolean;
  segment_id?: string;
  segment_category?: string;
  country?: string;
  region?: string;
  city?: string;
  device_type?: string;
  os?: string;
  browser?: string;
  is_mobile: boolean;
  timestamp: string;
}

export interface CampaignDetail extends Campaign {
  // Campaign stats
  campaign_id: string;
  total_bids: number;
  won_bids: number;
  conversions: number;
  total_spend: number;
  average_bid: number;
  win_rate: number;
  conversion_rate: number;
  cost_per_conversion: number;
  return_on_ad_spend: number;

  // Detailed breakdowns
  daily_metrics?: DailyMetric[];
  top_keywords?: KeywordStat[];
  device_breakdown?: DeviceStat[];
  geo_breakdown?: GeoStat[];
  recent_bids?: BidEvent[];
}

export interface CreateCampaignInput {
  name: string;
  budget?: number;
  daily_budget?: number;
  target_cpa?: number;
}

export interface UpdateCampaignInput {
  id: string;
  name?: string;
  status?: "active" | "paused" | "archived";
  budget?: number;
  daily_budget?: number;
  target_cpa?: number;
}

// API Functions

/**
 * List all campaigns for the current user with summary metrics
 */
export async function listCampaigns(): Promise<CampaignSummary[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/campaign.listWithMetrics`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({}),
    },
  );
  return handleResponse(response);
}

/**
 * Get a single campaign with detailed metrics
 */
export async function getCampaign(id: string): Promise<CampaignDetail> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.get`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify({ id }),
  });
  return handleResponse(response);
}

/**
 * Create a new campaign
 */
export async function createCampaign(
  data: CreateCampaignInput,
): Promise<Campaign> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.create`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse(response);
}

/**
 * Update an existing campaign
 */
export async function updateCampaign(
  data: UpdateCampaignInput,
): Promise<Campaign> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.update`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse(response);
}

/**
 * Delete (archive) a campaign
 */
export async function deleteCampaign(
  id: string,
): Promise<{ success: boolean; message: string }> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.delete`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify({ id }),
  });
  return handleResponse(response);
}

/**
 * Pause an active campaign
 */
export async function pauseCampaign(id: string): Promise<Campaign> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.pause`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify({ id }),
  });
  return handleResponse(response);
}

/**
 * Activate a paused campaign
 */
export async function activateCampaign(id: string): Promise<Campaign> {
  const response = await fetch(`${API_BASE_URL}/trpc/campaign.activate`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify({ id }),
  });
  return handleResponse(response);
}

/**
 * Get daily metrics for a campaign
 */
export async function getCampaignDailyMetrics(
  id: string,
): Promise<DailyMetric[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/campaign.getDailyMetrics`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ id }),
    },
  );
  return handleResponse(response);
}
