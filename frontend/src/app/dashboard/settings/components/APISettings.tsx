// frontend/src/app/dashboard/settings/components/APISettings.tsx
'use client';

import React, { useState } from 'react';
import { UserSettings, regenerateAPIKey } from '@/lib/api/settings';
import { Key, Copy, RefreshCw, Eye, EyeOff } from 'lucide-react';

interface APISettingsProps {
  settings: UserSettings;
  onUpdate: () => Promise<void>;
}

export function APISettings({ settings, onUpdate }: APISettingsProps) {
  const [showAPIKey, setShowAPIKey] = useState(false);
  const [regenerating, setRegenerating] = useState(false);

  const handleCopyAPIKey = () => {
    if (settings.api_key) {
      navigator.clipboard.writeText(settings.api_key);
      alert('API key copied to clipboard!');
    }
  };

  const handleRegenerateAPIKey = async () => {
    if (!confirm('Are you sure you want to regenerate your API key? This will invalidate the current key.')) {
      return;
    }

    setRegenerating(true);
    try {
      await regenerateAPIKey();
      await onUpdate();
      alert('API key regenerated successfully!');
    } catch (error) {
      console.error('Failed to regenerate API key:', error);
      alert('Failed to regenerate API key');
    } finally {
      setRegenerating(false);
    }
  };

  const maskAPIKey = (key: string) => {
    if (!key) return '';
    return `${key.substring(0, 8)}${'*'.repeat(Math.max(0, key.length - 12))}${key.substring(key.length - 4)}`;
  };

  return (
    <div className="space-y-6">
      {/* API Key Display */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          <div className="flex items-center gap-2">
            <Key className="w-4 h-4" />
            API Key
          </div>
        </label>
        <div className="flex gap-2">
          <div className="flex-1 relative">
            <input
              type="text"
              value={settings.api_key ? (showAPIKey ? settings.api_key : maskAPIKey(settings.api_key)) : 'No API key generated'}
              readOnly
              className="w-full px-4 py-2 pr-10 border border-gray-300 dark:border-slate-600 rounded-lg bg-gray-50 dark:bg-slate-700 dark:text-slate-100 font-mono text-sm"
            />
            {settings.api_key && (
              <button
                type="button"
                onClick={() => setShowAPIKey(!showAPIKey)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200"
              >
                {showAPIKey ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            )}
          </div>
          {settings.api_key && (
            <button
              onClick={handleCopyAPIKey}
              className="px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
              title="Copy to clipboard"
            >
              <Copy className="w-4 h-4" />
            </button>
          )}
          <button
            onClick={handleRegenerateAPIKey}
            disabled={regenerating}
            className="px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            title="Regenerate API key"
          >
            <RefreshCw className={`w-4 h-4 ${regenerating ? 'animate-spin' : ''}`} />
          </button>
        </div>
        <p className="text-sm text-gray-500 dark:text-slate-400 mt-2">
          Use this key to authenticate API requests. Keep it secure and never share it publicly.
        </p>
      </div>

      {/* API Rate Limit */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          API Rate Limit
        </label>
        <div className="px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-gray-50 dark:bg-slate-700">
          <div className="text-2xl font-bold text-gray-900 dark:text-slate-100">
            {settings.api_rate_limit.toLocaleString()}
          </div>
          <div className="text-sm text-gray-600 dark:text-slate-400">requests per hour</div>
        </div>
        <p className="text-sm text-gray-500 dark:text-slate-400 mt-2">
          Your current plan allows up to {settings.api_rate_limit.toLocaleString()} API requests per hour.
          Contact support to increase this limit.
        </p>
      </div>

      {/* API Documentation */}
      <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
        <h4 className="font-medium text-blue-900 dark:text-blue-100 mb-2">📚 API Documentation</h4>
        <p className="text-sm text-blue-800 dark:text-blue-200 mb-3">
          Learn how to integrate our API into your applications
        </p>
        <a
          href="/docs/api"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 text-sm text-blue-600 dark:text-blue-400 hover:underline"
        >
          View Documentation →
        </a>
      </div>

      {/* Example Usage */}
      <div>
        <h4 className="font-medium text-gray-900 dark:text-slate-100 mb-3">Example Usage</h4>
        <div className="bg-gray-900 dark:bg-slate-950 rounded-lg p-4 overflow-x-auto">
          <pre className="text-sm text-gray-100 font-mono">
{`# Submit a bid
curl -X POST https://api.yourplatform.com/v1/bids \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "campaign_id": "xxx-xxx-xxx",
    "bid_price": 2.50,
    "segment_id": "premium_users"
  }'`}
          </pre>
        </div>
      </div>
    </div>
  );
}
