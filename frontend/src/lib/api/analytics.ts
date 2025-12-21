// frontend/src/lib/api/analytics.ts
/**
 * Analytics API functions
 * Uses shared utilities and types from @/lib/utils and @/lib/types
 */

import { apiPost } from '@/lib/utils';
import type {
  PerformanceMetrics,
  KeywordAnalysis,
  DeviceBreakdown,
  GeoBreakdown,
  HourlyPerformance,
  DailyTrend,
  CompetitiveAnalysis,
  DateRangeParams,
} from '@/lib/types';

/**
 * Get overall performance metrics
 */
export async function getPerformanceOverview(
  params?: DateRangeParams,
): Promise<PerformanceMetrics> {
  return apiPost<PerformanceMetrics>('/trpc/analytics.getPerformanceOverview', params || {});
}

/**
 * Get keyword performance analysis
 */
export async function getKeywordAnalysis(
  params?: DateRangeParams & { limit?: number },
): Promise<KeywordAnalysis[]> {
  return apiPost<KeywordAnalysis[]>('/trpc/analytics.getKeywordAnalysis', params || {});
}

/**
 * Get device breakdown
 */
export async function getDeviceBreakdown(
  params?: DateRangeParams,
): Promise<DeviceBreakdown[]> {
  return apiPost<DeviceBreakdown[]>('/trpc/analytics.getDeviceBreakdown', params || {});
}

/**
 * Get geographic breakdown
 */
export async function getGeoBreakdown(
  params?: DateRangeParams & { limit?: number },
): Promise<GeoBreakdown[]> {
  return apiPost<GeoBreakdown[]>('/trpc/analytics.getGeoBreakdown', params || {});
}

/**
 * Get hourly performance
 */
export async function getHourlyPerformance(
  params?: DateRangeParams,
): Promise<HourlyPerformance[]> {
  return apiPost<HourlyPerformance[]>('/trpc/analytics.getHourlyPerformance', params || {});
}

/**
 * Get daily performance trends
 */
export async function getDailyTrends(
  params?: DateRangeParams,
): Promise<DailyTrend[]> {
  return apiPost<DailyTrend[]>('/trpc/analytics.getDailyTrends', params || {});
}

/**
 * Get competitive analysis
 */
export async function getCompetitiveAnalysis(
  params?: DateRangeParams,
): Promise<CompetitiveAnalysis[]> {
  return apiPost<CompetitiveAnalysis[]>('/trpc/analytics.getCompetitiveAnalysis', params || {});
}

// Re-export types for convenience
export type {
  PerformanceMetrics,
  KeywordAnalysis,
  DeviceBreakdown,
  GeoBreakdown,
  HourlyPerformance,
  DailyTrend,
  CompetitiveAnalysis,
};
