// frontend/src/lib/api/alerts.ts
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
  return data.result?.data || data;
};

// Types
export type AlertType = 
  | "fraud"
  | "budget"
  | "performance"
  | "model"
  | "system"
  | "campaign";

export type AlertSeverity = "low" | "medium" | "high" | "critical";
export type AlertStatus = "unread" | "read" | "acknowledged" | "resolved" | "dismissed";

export interface Alert {
  id: string;
  type: AlertType;
  severity: AlertSeverity;
  status: AlertStatus;
  title: string;
  message: string;
  campaign_id?: string;
  campaign_name?: string;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
  acknowledged_at?: string;
  resolved_at?: string;
  notes?: string;
}

export interface AlertOverview {
  total_alerts: number;
  unread_alerts: number;
  critical_alerts: number;
  alerts_by_type: Record<AlertType, number>;
  alerts_by_severity: Record<AlertSeverity, number>;
  recent_trend: {
    date: string;
    count: number;
  }[];
}

export interface AlertRule {
  id: string;
  name: string;
  description: string;
  type: AlertType;
  severity: AlertSeverity;
  enabled: boolean;
  conditions: {
    metric: string;
    operator: "gt" | "lt" | "eq" | "gte" | "lte";
    threshold: number;
    duration_minutes?: number;
  }[];
  notification_channels: ("email" | "dashboard" | "webhook")[];
  created_at: string;
  updated_at: string;
}

// API Functions

/**
 * Get all alerts with optional filtering
 */
export async function getAlerts(params?: {
  type?: AlertType;
  severity?: AlertSeverity;
  status?: AlertStatus;
  campaign_id?: string;
  start_date?: string;
  end_date?: string;
  limit?: number;
  offset?: number;
}): Promise<Alert[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.getAlerts`,
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
 * Get alert overview and statistics
 */
export async function getAlertOverview(days: number = 30): Promise<AlertOverview> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.getOverview`,
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
 * Update alert status
 */
export async function updateAlertStatus(
  alertId: string,
  status: AlertStatus,
  notes?: string
): Promise<{ success: boolean; message: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.updateStatus`,
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
 * Bulk update alert statuses
 */
export async function bulkUpdateAlerts(
  alertIds: string[],
  status: AlertStatus
): Promise<{ success: boolean; message: string; updated_count: number }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.bulkUpdate`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({
        alert_ids: alertIds,
        status,
      }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Get alert rules
 */
export async function getAlertRules(): Promise<AlertRule[]> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.getRules`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({}),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Create alert rule
 */
export async function createAlertRule(
  rule: Omit<AlertRule, "id" | "created_at" | "updated_at">
): Promise<{ success: boolean; rule_id: string; message: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.createRule`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(rule),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Update alert rule
 */
export async function updateAlertRule(
  ruleId: string,
  updates: Partial<Omit<AlertRule, "id" | "created_at" | "updated_at">>
): Promise<{ success: boolean; message: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.updateRule`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({
        rule_id: ruleId,
        ...updates,
      }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

/**
 * Delete alert rule
 */
export async function deleteAlertRule(
  ruleId: string
): Promise<{ success: boolean; message: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/alerts.deleteRule`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ rule_id: ruleId }),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

// Helper functions

export function getSeverityColor(severity: AlertSeverity): string {
  switch (severity) {
    case "critical":
      return "text-red-600";
    case "high":
      return "text-orange-600";
    case "medium":
      return "text-yellow-600";
    case "low":
      return "text-blue-600";
    default:
      return "text-slate-600";
  }
}

export function getSeverityBadgeColor(severity: AlertSeverity): string {
  switch (severity) {
    case "critical":
      return "bg-red-100 text-red-800 border-red-200";
    case "high":
      return "bg-orange-100 text-orange-800 border-orange-200";
    case "medium":
      return "bg-yellow-100 text-yellow-800 border-yellow-200";
    case "low":
      return "bg-blue-100 text-blue-800 border-blue-200";
    default:
      return "bg-slate-100 text-slate-800 border-slate-200";
  }
}

export function getStatusColor(status: AlertStatus): string {
  switch (status) {
    case "unread":
      return "text-red-600";
    case "read":
      return "text-yellow-600";
    case "acknowledged":
      return "text-blue-600";
    case "resolved":
      return "text-green-600";
    case "dismissed":
      return "text-gray-600";
    default:
      return "text-slate-600";
  }
}

export function getStatusBadgeColor(status: AlertStatus): string {
  switch (status) {
    case "unread":
      return "bg-red-100 text-red-800 border-red-200";
    case "read":
      return "bg-yellow-100 text-yellow-800 border-yellow-200";
    case "acknowledged":
      return "bg-blue-100 text-blue-800 border-blue-200";
    case "resolved":
      return "bg-green-100 text-green-800 border-green-200";
    case "dismissed":
      return "bg-gray-100 text-gray-800 border-gray-200";
    default:
      return "bg-slate-100 text-slate-800 border-slate-200";
  }
}

export function getTypeIcon(type: AlertType): string {
  switch (type) {
    case "fraud":
      return "Shield";
    case "budget":
      return "DollarSign";
    case "performance":
      return "TrendingDown";
    case "model":
      return "Brain";
    case "system":
      return "AlertTriangle";
    case "campaign":
      return "Target";
    default:
      return "Bell";
  }
}

export function getTypeColor(type: AlertType): string {
  switch (type) {
    case "fraud":
      return "text-red-600";
    case "budget":
      return "text-green-600";
    case "performance":
      return "text-purple-600";
    case "model":
      return "text-blue-600";
    case "system":
      return "text-orange-600";
    case "campaign":
      return "text-indigo-600";
    default:
      return "text-slate-600";
  }
}

export function formatAlertType(type: AlertType): string {
  return type.charAt(0).toUpperCase() + type.slice(1);
}
