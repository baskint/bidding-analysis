// frontend/src/lib/utils.ts
/**
 * Shared utility functions
 * Common helpers used across the application
 */

import type { ApiResponse, ApiError } from './types';

// ============================================================================
// API Configuration
// ============================================================================

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

/**
 * Get authentication headers with JWT token
 */
export function getAuthHeaders(): Record<string, string> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  };
}

// ============================================================================
// API Request Helpers
// ============================================================================

/**
 * Handle API response and extract data
 */
export async function handleApiResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `API call failed: ${response.status}`);
  }

  const data: ApiResponse<T> = await response.json();
  return data.result?.data || (data as any);
}

/**
 * Make a GET request to the API
 */
export async function apiGet<T>(endpoint: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });
  return handleApiResponse<T>(response);
}

/**
 * Make a POST request to the API
 */
export async function apiPost<T>(endpoint: string, body?: any): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: body ? JSON.stringify(body) : undefined,
  });
  return handleApiResponse<T>(response);
}

/**
 * Make a PUT request to the API
 */
export async function apiPut<T>(endpoint: string, body?: any): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: 'PUT',
    headers: getAuthHeaders(),
    body: body ? JSON.stringify(body) : undefined,
  });
  return handleApiResponse<T>(response);
}

/**
 * Make a DELETE request to the API
 */
export async function apiDelete<T>(endpoint: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
  });
  return handleApiResponse<T>(response);
}

// ============================================================================
// Date & Time Utilities
// ============================================================================

/**
 * Format a date for display
 */
export function formatDate(date: string | Date): string {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

/**
 * Format a datetime for display
 */
export function formatDateTime(date: string | Date): string {
  return new Date(date).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * Format time only
 */
export function formatTime(date: string | Date): string {
  return new Date(date).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

/**
 * Get relative time (e.g., "2 hours ago")
 */
export function getRelativeTime(date: string | Date): string {
  const now = new Date();
  const then = new Date(date);
  const seconds = Math.floor((now.getTime() - then.getTime()) / 1000);

  if (seconds < 60) return 'just now';
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;
  return formatDate(date);
}

/**
 * Calculate date range (for analytics)
 */
export function getDateRange(days: number): { start_date: string; end_date: string } {
  const end = new Date();
  const start = new Date();
  start.setDate(end.getDate() - days);

  return {
    start_date: start.toISOString().split('T')[0],
    end_date: end.toISOString().split('T')[0],
  };
}

// ============================================================================
// Number Formatting
// ============================================================================

/**
 * Format currency
 */
export function formatCurrency(value: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(value);
}

/**
 * Format percentage
 */
export function formatPercent(value: number, decimals = 1): string {
  return `${(value * 100).toFixed(decimals)}%`;
}

/**
 * Format large numbers with K, M, B suffixes
 */
export function formatCompactNumber(value: number): string {
  if (value >= 1000000000) return `${(value / 1000000000).toFixed(1)}B`;
  if (value >= 1000000) return `${(value / 1000000).toFixed(1)}M`;
  if (value >= 1000) return `${(value / 1000).toFixed(1)}K`;
  return value.toFixed(0);
}

// ============================================================================
// String Utilities
// ============================================================================

/**
 * Capitalize first letter
 */
export function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Convert snake_case to Title Case
 */
export function snakeToTitle(str: string): string {
  return str
    .split('_')
    .map(word => capitalize(word))
    .join(' ');
}

/**
 * Truncate text with ellipsis
 */
export function truncate(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength - 3) + '...';
}

// ============================================================================
// Validation Utilities
// ============================================================================

/**
 * Validate email format
 */
export function isValidEmail(email: string): boolean {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return regex.test(email);
}

/**
 * Validate number is positive
 */
export function isPositive(value: number): boolean {
  return value > 0;
}

/**
 * Validate number is in range
 */
export function isInRange(value: number, min: number, max: number): boolean {
  return value >= min && value <= max;
}

// ============================================================================
// Array Utilities
// ============================================================================

/**
 * Get unique values from array
 */
export function unique<T>(arr: T[]): T[] {
  return Array.from(new Set(arr));
}

/**
 * Sort array by key
 */
export function sortBy<T>(arr: T[], key: keyof T, order: 'asc' | 'desc' = 'asc'): T[] {
  return [...arr].sort((a, b) => {
    const aVal = a[key];
    const bVal = b[key];
    if (aVal < bVal) return order === 'asc' ? -1 : 1;
    if (aVal > bVal) return order === 'asc' ? 1 : -1;
    return 0;
  });
}

/**
 * Group array by key
 */
export function groupBy<T>(arr: T[], key: keyof T): Record<string, T[]> {
  return arr.reduce((groups, item) => {
    const groupKey = String(item[key]);
    return {
      ...groups,
      [groupKey]: [...(groups[groupKey] || []), item],
    };
  }, {} as Record<string, T[]>);
}

// ============================================================================
// Local Storage Utilities
// ============================================================================

/**
 * Safe localStorage get
 */
export function getFromStorage<T>(key: string, defaultValue: T): T {
  if (typeof window === 'undefined') return defaultValue;

  try {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : defaultValue;
  } catch {
    return defaultValue;
  }
}

/**
 * Safe localStorage set
 */
export function setInStorage<T>(key: string, value: T): void {
  if (typeof window === 'undefined') return;

  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.error('Failed to save to localStorage:', error);
  }
}

/**
 * Remove from localStorage
 */
export function removeFromStorage(key: string): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(key);
}

// ============================================================================
// Debounce & Throttle
// ============================================================================

/**
 * Debounce function calls
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;

  return function (...args: Parameters<T>) {
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

/**
 * Throttle function calls
 */
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean;

  return function (...args: Parameters<T>) {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

// ============================================================================
// Class Name Utilities
// ============================================================================

/**
 * Conditionally join class names
 */
export function cn(...classes: (string | boolean | undefined | null)[]): string {
  return classes.filter(Boolean).join(' ');
}
