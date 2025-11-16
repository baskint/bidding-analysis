// frontend/src/app/dashboard/settings/components/NotificationSettings.tsx
'use client';

import React, { useState } from 'react';
import { UserSettings, UserSettingsUpdate } from '@/lib/api/settings';
import { Mail, MessageSquare, Webhook, Clock } from 'lucide-react';

interface NotificationSettingsProps {
  settings: UserSettings;
  onUpdate: (update: UserSettingsUpdate) => Promise<void>;
}

export function NotificationSettings({ settings, onUpdate }: NotificationSettingsProps) {
  const [emailNotifications, setEmailNotifications] = useState(settings.email_notifications);
  const [slackNotifications, setSlackNotifications] = useState(settings.slack_notifications);
  const [webhookNotifications, setWebhookNotifications] = useState(settings.webhook_notifications);
  const [alertFrequency, setAlertFrequency] = useState(settings.alert_frequency);
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    setSaving(true);
    try {
      await onUpdate({
        email_notifications: emailNotifications,
        slack_notifications: slackNotifications,
        webhook_notifications: webhookNotifications,
        alert_frequency: alertFrequency,
      });
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Email Notifications */}
      <div className="flex items-center justify-between">
        <div className="flex items-start gap-3">
          <Mail className="w-5 h-5 text-gray-600 dark:text-slate-400 mt-0.5" />
          <div>
            <div className="font-medium text-gray-900 dark:text-slate-100">Email Notifications</div>
            <div className="text-sm text-gray-600 dark:text-slate-400">
              Receive alerts and updates via email
            </div>
          </div>
        </div>
        <label className="relative inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            checked={emailNotifications}
            onChange={(e) => setEmailNotifications(e.target.checked)}
            className="sr-only peer"
          />
          <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
        </label>
      </div>

      {/* Slack Notifications */}
      <div className="flex items-center justify-between">
        <div className="flex items-start gap-3">
          <MessageSquare className="w-5 h-5 text-gray-600 dark:text-slate-400 mt-0.5" />
          <div>
            <div className="font-medium text-gray-900 dark:text-slate-100">Slack Notifications</div>
            <div className="text-sm text-gray-600 dark:text-slate-400">
              Send alerts to your Slack workspace
            </div>
          </div>
        </div>
        <label className="relative inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            checked={slackNotifications}
            onChange={(e) => setSlackNotifications(e.target.checked)}
            className="sr-only peer"
          />
          <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
        </label>
      </div>

      {/* Webhook Notifications */}
      <div className="flex items-center justify-between">
        <div className="flex items-start gap-3">
          <Webhook className="w-5 h-5 text-gray-600 dark:text-slate-400 mt-0.5" />
          <div>
            <div className="font-medium text-gray-900 dark:text-slate-100">Webhook Notifications</div>
            <div className="text-sm text-gray-600 dark:text-slate-400">
              Send alerts to custom webhook endpoints
            </div>
          </div>
        </div>
        <label className="relative inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            checked={webhookNotifications}
            onChange={(e) => setWebhookNotifications(e.target.checked)}
            className="sr-only peer"
          />
          <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
        </label>
      </div>

      {/* Alert Frequency */}
      <div className="pt-4 border-t border-gray-200 dark:border-slate-700">
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          <div className="flex items-center gap-2">
            <Clock className="w-4 h-4" />
            Alert Frequency
          </div>
        </label>
        <select
          value={alertFrequency}
          onChange={(e) => setAlertFrequency(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
        >
          <option value="realtime">Real-time (Instant)</option>
          <option value="hourly">Hourly Digest</option>
          <option value="daily">Daily Summary</option>
          <option value="weekly">Weekly Report</option>
        </select>
        <p className="text-sm text-gray-500 dark:text-slate-400 mt-2">
          Choose how often you want to receive notification alerts
        </p>
      </div>

      {/* Save Button */}
      <div className="flex justify-end pt-4 border-t border-gray-200 dark:border-slate-700">
        <button
          onClick={handleSave}
          disabled={saving}
          className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {saving ? 'Saving...' : 'Save Changes'}
        </button>
      </div>
    </div>
  );
}
