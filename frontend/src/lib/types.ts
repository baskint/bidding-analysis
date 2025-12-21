// frontend/src/lib/types.ts
/**
 * Shared TypeScript types and interfaces
 * Centralized type definitions to avoid duplication
 */

// ============================================================================
// Common API Response Types
// ============================================================================

export interface ApiResponse<T> {
  result: {
    data: T;
    type: 'data';
  };
}

export interface ApiError {
  error: string;
  code?: string;
  details?: Record<string, unknown>;
}

// ============================================================================
// Domain Models
// ============================================================================

export interface User {
  id: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface Campaign {
  id: string;
  user_id: string;
  name: string;
  budget: number;
  daily_budget: number;
  target_cpa: number;
  status: 'active' | 'paused' | 'completed';
  created_at: string;
  updated_at: string;
}

export interface BidEvent {
  bid_event_id: string;
  campaign_id: string;
  user_id: string;
  bid_price: number;
  win_price?: number;
  floor_price: number;
  won: boolean;
  converted: boolean;
  timestamp: string;
  segment_category: string;
  segment_id: string;
  device_type: 'mobile' | 'desktop' | 'tablet';
  os: string;
  browser: string;
  country: string;
  keywords?: string[];
}

export interface MLModel {
  id: string;
  user_id: string;
  name: string;
  version: string;
  model_type: string;
  status: 'training' | 'active' | 'inactive' | 'failed';
  accuracy?: number;
  created_at: string;
  updated_at: string;
}

export interface FraudAlert {
  id: string;
  bid_event_id: string;
  fraud_type: string;
  confidence_score: number;
  details: Record<string, unknown>;
  status: 'pending' | 'acknowledged' | 'dismissed';
  created_at: string;
}

// ============================================================================
// Analytics Types
// ============================================================================

export interface PerformanceMetrics {
  total_bids: number;
  won_bids: number;
  conversions: number;
  total_spend: number;
  revenue: number;
  win_rate: number;
  conversion_rate: number;
  average_bid: number;
  cpa: number;
  roas: number;
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
  region?: string;
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

export interface CompetitiveAnalysis {
  segment_category: string;
  our_win_rate: number;
  market_average_bid: number;
  our_average_bid: number;
  average_floor_price: number;
  competition_intensity: number;
  total_opportunities: number;
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

// ============================================================================
// Prediction Types
// ============================================================================

export interface PredictionRequest {
  campaign_id: string;
  floor_price: number;
  user_segment: string;
  device_type: string;
  country: string;
  keywords?: string[];
  engagement_score: number;
  conversion_probability: number;
}

export interface PredictionResult {
  predicted_bid: number;
  confidence: number;
  strategy: 'ml_optimized' | 'rule_based_fallback';
  fraud_risk: boolean;
  reasoning: string;
}

// ============================================================================
// Settings Types
// ============================================================================

export interface UserSettings {
  email_notifications: boolean;
  slack_notifications: boolean;
  fraud_alert_threshold: number;
  budget_alert_threshold: number;
  performance_alert_threshold: number;
}

// ============================================================================
// Utility Types
// ============================================================================

export interface DateRangeParams {
  start_date?: string;
  end_date?: string;
}

export interface PaginationParams {
  page?: number;
  limit?: number;
}

export type SortOrder = 'asc' | 'desc';

export interface SortParams {
  sort_by?: string;
  sort_order?: SortOrder;
}
