// frontend/src/lib/api/settings.ts
/**
 * User settings and integrations API functions
 */

import { apiGet, apiPost, apiDelete } from '@/lib/utils';

// Integration Config Types
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
  config: IntegrationConfig;
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
  config?: IntegrationConfig;
}

export interface IntegrationUpdate {
  integration_name?: string;
  access_token?: string;
  refresh_token?: string;
  api_key?: string;
  api_secret?: string;
  webhook_url?: string;
  config?: IntegrationConfig;
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

/**
 * Get current user settings
 */
export async function getUserSettings(): Promise<UserSettings> {
  return apiGet<UserSettings>('/trpc/settings.get');
}

/**
 * Update user settings
 */
export async function updateUserSettings(update: UserSettingsUpdate): Promise<UserSettings> {
  return apiPost<UserSettings>('/trpc/settings.update', update);
}

/**
 * Regenerate API key
 */
export async function regenerateAPIKey(): Promise<{ api_key: string }> {
  return apiPost<{ api_key: string }>('/trpc/settings.regenerateAPIKey', {});
}

// Integrations API

/**
 * List all integrations
 */
export async function listIntegrations(): Promise<{ integrations: Integration[] }> {
  return apiGet<{ integrations: Integration[] }>('/trpc/integrations.list');
}

/**
 * Get a single integration
 */
export async function getIntegration(id: string): Promise<Integration> {
  return apiGet<Integration>(`/trpc/integrations.get?id=${id}`);
}

/**
 * Create a new integration
 */
export async function createIntegration(data: IntegrationCreate): Promise<Integration> {
  return apiPost<Integration>('/trpc/integrations.create', data);
}

/**
 * Update an integration
 */
export async function updateIntegration(id: string, data: IntegrationUpdate): Promise<Integration> {
  return apiPost<Integration>('/trpc/integrations.update', { id, ...data });
}

/**
 * Delete an integration
 */
export async function deleteIntegration(id: string): Promise<void> {
  return apiDelete<void>(`/trpc/integrations.delete?id=${id}`);
}

/**
 * Test an integration connection
 */
export async function testIntegration(id: string): Promise<{ message: string }> {
  return apiPost<{ message: string }>('/trpc/integrations.test', { id });
}

/**
 * Get billing information
 */
export async function getBillingInfo(): Promise<BillingInfo> {
  return apiGet<BillingInfo>('/trpc/billing.get');
}
