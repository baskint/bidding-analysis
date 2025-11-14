// frontend/src/lib/api/fraud.ts
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
export interface FraudOverview {
  total_alerts: number;
  active_alerts: number;
  blocked_bids: number;
  amount_saved: number;
  threat_level: "low" | "medium" | "high" | "critical";
  alerts_by_type: Record<string, number>;
  top_affected_campaigns: CampaignRisk[];
}

export interface CampaignRisk {
  campaign_id: string;
  campaign_name: string;
  risk_score: number;
  fraud_attempts: number;
  threat_level: "low" | "medium" | "high" | "critical";
}

export interface FraudAlert {
  id: string;
  campaign_id: string;
  alert_type: string;
  severity: number;
  description: string;
  affected_user_ids: string[];
  detected_at: string;
  resolved_at?: string;
  status: "active" | "investigating" | "resolved" | "false_positive";
}

export interface FraudTrend {
  date: string;
  fraud_attempts: number;
  blocked_bids: number;
  amount_saved: number;
  alert_type: string;
}

export interface DeviceFraud {
  device_type: string;
  browser: string;
  os: string;
  total_bids: number;
  fraud_bids: number;
  fraud_rate: number;
}

export interface GeoFraud {
  country: string;
  region: string;
  city: string;
  total_bids: number;
  fraud_bids: number;
  fraud_rate: number;
}

// API Functions

/**
 * Get fraud overview metrics
 */
export async function getFraudOverview(days: number = 30): Promise<FraudOverview> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/fraud.getOverview`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ days }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Get fraud alerts with optional filtering
 */
export async function getRealFraudAlerts(params?: {
  status?: string;
  min_severity?: number;
  alert_type?: string;
  start_date?: string;
  end_date?: string;
  limit?: number;
}): Promise<FraudAlert[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/fraud.getAlerts`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params || {}),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Update fraud alert status
 */
export async function updateFraudAlert(
  alertId: string,
  status: string,
  notes?: string
): Promise<{ success: boolean; message: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/fraud.updateAlert`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({
        alert_id: alertId,
        status,
        notes,
      }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Get fraud trends over time
 */
export async function getFraudTrends(days: number = 30): Promise<FraudTrend[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/fraud.getTrends`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ days }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Get device-specific fraud analysis
 */
export async function getDeviceFraudAnalysis(days: number = 30): Promise<DeviceFraud[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/fraud.getDeviceAnalysis`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ days }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Get geographic fraud analysis
 */
export async function getGeoFraudAnalysis(days: number = 30): Promise<GeoFraud[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/fraud.getGeoAnalysis`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ days }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Create a new fraud alert
 */
export async function createFraudAlert(params: {
  campaign_id: string;
  alert_type: string;
  severity: number;
  description: string;
  affected_user_ids?: string[];
}): Promise<{ success: boolean; alert_id: string; message: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/fraud.createAlert`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(params),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

// Helper functions

export function getSeverityColor(severity: number): string {
  if (severity >= 8) return "text-red-600";
  if (severity >= 6) return "text-orange-600";
  if (severity >= 4) return "text-yellow-600";
  return "text-blue-600";
}

export function getSeverityBadgeColor(severity: number): string {
  if (severity >= 8) return "bg-red-100 text-red-800";
  if (severity >= 6) return "bg-orange-100 text-orange-800";
  if (severity >= 4) return "bg-yellow-100 text-yellow-800";
  return "bg-blue-100 text-blue-800";
}

export function getThreatLevelColor(level: string): string {
  switch (level) {
    case "critical":
      return "text-red-600";
    case "high":
      return "text-orange-600";
    case "medium":
      return "text-yellow-600";
    default:
      return "text-green-600";
  }
}

export function getThreatLevelBadgeColor(level: string): string {
  switch (level) {
    case "critical":
      return "bg-red-100 text-red-800";
    case "high":
      return "bg-orange-100 text-orange-800";
    case "medium":
      return "bg-yellow-100 text-yellow-800";
    default:
      return "bg-green-100 text-green-800";
  }
}

export function formatAlertType(type: string): string {
  return type
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

export function getStatusColor(status: string): string {
  switch (status) {
    case "active":
      return "text-red-600";
    case "investigating":
      return "text-yellow-600";
    case "resolved":
      return "text-green-600";
    case "false_positive":
      return "text-gray-600";
    default:
      return "text-slate-600";
  }
}

export function getStatusBadgeColor(status: string): string {
  switch (status) {
    case "active":
      return "bg-red-100 text-red-800";
    case "investigating":
      return "bg-yellow-100 text-yellow-800";
    case "resolved":
      return "bg-green-100 text-green-800";
    case "false_positive":
      return "bg-gray-100 text-gray-800";
    default:
      return "bg-slate-100 text-slate-800";
  }
}
