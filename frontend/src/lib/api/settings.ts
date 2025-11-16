// frontend/src/lib/api/settings.ts
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

// Specific integration config types
export interface GoogleAdsConfig {
  customer_id: string;
  developer_token: string;
  manager_account_id?: string;
  include_metrics: boolean;
  conversion_tracking?: boolean;
  campaign_types?: string[];
}

export interface FacebookAdsConfig {
  ad_account_id: string;
  access_level: string;
  include_deleted: boolean;
  data_retention_days?: number;
  include_insights?: boolean;
}

export interface SlackConfig {
  channel_id: string;
  bot_token: string;
  notifications_enabled: boolean;
  alert_types?: string[];
  mention_users?: string[];
}

export interface WebhookConfig {
  url: string;
  secret?: string;
  events: string[];
  retry_attempts?: number;
  timeout_ms?: number;
}

export interface GenericIntegrationConfig {
  [key: string]: string | number | boolean | string[] | null | undefined;
}

// Union type for all possible configs
export type IntegrationConfig = 
  | GoogleAdsConfig 
  | FacebookAdsConfig 
  | SlackConfig 
  | WebhookConfig 
  | GenericIntegrationConfig;

// Helper to check config type
export function isGoogleAdsConfig(config: IntegrationConfig): config is GoogleAdsConfig {
  return 'customer_id' in config && 'developer_token' in config;
}

export function isFacebookAdsConfig(config: IntegrationConfig): config is FacebookAdsConfig {
  return 'ad_account_id' in config && 'access_level' in config;
}

export function isSlackConfig(config: IntegrationConfig): config is SlackConfig {
  return 'channel_id' in config && 'bot_token' in config;
}

export function isWebhookConfig(config: IntegrationConfig): config is WebhookConfig {
  return 'url' in config && 'events' in config;
}

// Types
export interface UserSettings {
  id: string;
  user_id: string;
  full_name?: string;
  email?: string;
  phone?: string;
  timezone: string;
  language: string;
  email_notifications: boolean;
  slack_notifications: boolean;
  webhook_notifications: boolean;
  alert_frequency: string;
  fraud_alert_threshold: number;
  budget_alert_threshold: number;
  performance_alert_threshold: number;
  default_dashboard_view: string;
  default_date_range: string;
  dark_mode: boolean;
  api_key?: string;
  api_rate_limit: number;
  created_at: string;
  updated_at: string;
}

export interface UserSettingsUpdate {
  full_name?: string;
  email?: string;
  phone?: string;
  timezone?: string;
  language?: string;
  email_notifications?: boolean;
  slack_notifications?: boolean;
  webhook_notifications?: boolean;
  alert_frequency?: string;
  fraud_alert_threshold?: number;
  budget_alert_threshold?: number;
  performance_alert_threshold?: number;
  default_dashboard_view?: string;
  default_date_range?: string;
  dark_mode?: boolean;
}

export interface Integration {
  id: string;
  user_id: string;
  provider: string;
  integration_name: string;
  auth_type: string;
  access_token?: string;
  refresh_token?: string;
  api_key?: string;
  api_secret?: string;
  webhook_url?: string;
  token_expires_at?: string;
  config: IntegrationConfig; // Now uses the union type
  status: string;
  last_sync_at?: string;
  last_error?: string;
  created_at: string;
  updated_at: string;
}

export interface IntegrationCreate {
  provider: string;
  integration_name: string;
  auth_type: string;
  access_token?: string;
  refresh_token?: string;
  api_key?: string;
  api_secret?: string;
  webhook_url?: string;
  config?: IntegrationConfig; // Now uses the union type
}

export interface IntegrationUpdate {
  integration_name?: string;
  access_token?: string;
  refresh_token?: string;
  api_key?: string;
  api_secret?: string;
  webhook_url?: string;
  config?: IntegrationConfig; // Now uses the union type
  status?: string;
}

export interface BillingInfo {
  id: string;
  user_id: string;
  stripe_customer_id?: string;
  stripe_subscription_id?: string;
  plan_type: string;
  billing_cycle: string;
  monthly_bid_limit?: number;
  campaigns_limit?: number;
  ml_models_limit?: number;
  subscription_status: string;
  trial_ends_at?: string;
  current_period_start?: string;
  current_period_end?: string;
  created_at: string;
  updated_at: string;
}

// User Settings API
// ... (all your API functions remain the same, they'll automatically use the new types)
export async function getUserSettings(): Promise<UserSettings> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/settings.get`,
    {
      method: "GET",
      headers: getAuthHeaders(),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

export async function updateUserSettings(update: UserSettingsUpdate): Promise<UserSettings> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/settings.update`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(update),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

export async function regenerateAPIKey(): Promise<{ api_key: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/settings.regenerateAPIKey`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

// Integrations API
export async function listIntegrations(): Promise<{ integrations: Integration[] }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/integrations.list`,
    {
      method: "GET",
      headers: getAuthHeaders(),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

export async function getIntegration(id: string): Promise<Integration> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/integrations.get?id=${id}`,
    {
      method: "GET",
      headers: getAuthHeaders(),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

export async function createIntegration(data: IntegrationCreate): Promise<Integration> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/integrations.create`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

export async function updateIntegration(id: string, data: IntegrationUpdate): Promise<Integration> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/integrations.update?id=${id}`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

export async function deleteIntegration(id: string): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/integrations.delete?id=${id}`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

export async function testIntegration(id: string): Promise<{ message: string }> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/integrations.test?id=${id}`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

// Billing API
export async function getBillingInfo(): Promise<BillingInfo> {
  const response = await fetch(
    `${API_BASE_URL}/trpc/billing.get`,
    {
      method: "GET",
      headers: getAuthHeaders(),
      cache: 'no-store',
      signal: AbortSignal.timeout(30000),
    },
  );
  return handleResponse(response);
}

// Helper functions
export function getIntegrationStatusColor(status: string): string {
  switch (status) {
    case "active":
      return "text-green-600";
    case "error":
      return "text-red-600";
    case "disabled":
      return "text-gray-600";
    case "pending":
      return "text-yellow-600";
    default:
      return "text-slate-600";
  }
}

export function getIntegrationStatusBadgeColor(status: string): string {
  switch (status) {
    case "active":
      return "bg-green-100 text-green-800";
    case "error":
      return "bg-red-100 text-red-800";
    case "disabled":
      return "bg-gray-100 text-gray-800";
    case "pending":
      return "bg-yellow-100 text-yellow-800";
    default:
      return "bg-slate-100 text-slate-800";
  }
}

export function getPlanTypeColor(plan: string): string {
  switch (plan) {
    case "enterprise":
      return "text-purple-600";
    case "premium":
      return "text-blue-600";
    case "standard":
      return "text-green-600";
    case "free":
      return "text-gray-600";
    default:
      return "text-slate-600";
  }
}

export function getPlanTypeBadgeColor(plan: string): string {
  switch (plan) {
    case "enterprise":
      return "bg-purple-100 text-purple-800";
    case "premium":
      return "bg-blue-100 text-blue-800";
    case "standard":
      return "bg-green-100 text-green-800";
    case "free":
      return "bg-gray-100 text-gray-800";
    default:
      return "bg-slate-100 text-slate-800";
  }
}

export function formatAlertFrequency(frequency: string): string {
  switch (frequency) {
    case "realtime":
      return "Real-time";
    case "hourly":
      return "Hourly";
    case "daily":
      return "Daily";
    case "weekly":
      return "Weekly";
    default:
      return frequency
        .split("_")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ");
  }
}

export function formatProviderName(provider: string): string {
  return provider
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}
