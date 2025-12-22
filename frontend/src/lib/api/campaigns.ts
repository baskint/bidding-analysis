// frontend/src/lib/api/campaigns.ts
/**
 * Campaign API functions
 */

import { apiPost } from '@/lib/utils';
import type { Campaign } from '@/lib/types';

// Campaign-specific types that extend the base Campaign type
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

export interface BidEventSummary {
  id: string;
  campaign_id: string;
  bid_price: number;
  win_price?: number;
  floor_price?: number;
  won: boolean;
  converted: boolean;
  device_type?: string;
  country?: string;
  city?: string;
  region?: string;
  timestamp: string;
}

export interface CampaignDetail extends Campaign {
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
  daily_metrics?: DailyMetric[];
  top_keywords?: KeywordStat[];
  device_breakdown?: DeviceStat[];
  geo_breakdown?: GeoStat[];
  recent_bids?: BidEventSummary[];
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
  return apiPost<CampaignSummary[]>('/trpc/campaign.listWithMetrics', {});
}

/**
 * Get a single campaign with detailed metrics
 */
export async function getCampaign(id: string): Promise<CampaignDetail> {
  return apiPost<CampaignDetail>('/trpc/campaign.get', { id });
}

/**
 * Create a new campaign
 */
export async function createCampaign(data: CreateCampaignInput): Promise<Campaign> {
  return apiPost<Campaign>('/trpc/campaign.create', data);
}

/**
 * Update an existing campaign
 */
export async function updateCampaign(data: UpdateCampaignInput): Promise<Campaign> {
  return apiPost<Campaign>('/trpc/campaign.update', data);
}

/**
 * Delete (archive) a campaign
 */
export async function deleteCampaign(id: string): Promise<{ success: boolean; message: string }> {
  return apiPost<{ success: boolean; message: string }>('/trpc/campaign.delete', { id });
}

/**
 * Pause an active campaign
 */
export async function pauseCampaign(id: string): Promise<Campaign> {
  return apiPost<Campaign>('/trpc/campaign.pause', { id });
}

/**
 * Activate a paused campaign
 */
export async function activateCampaign(id: string): Promise<Campaign> {
  return apiPost<Campaign>('/trpc/campaign.activate', { id });
}

/**
 * Get daily metrics for a campaign
 */
export async function getCampaignDailyMetrics(id: string): Promise<DailyMetric[]> {
  return apiPost<DailyMetric[]>('/trpc/campaign.getDailyMetrics', { id });
}
