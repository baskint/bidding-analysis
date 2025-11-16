// frontend/src/app/dashboard/settings/page.tsx
'use client';

import React, { useState, useEffect } from 'react';
import {
  getUserSettings,
  updateUserSettings,
  UserSettings as UserSettingsType,
  UserSettingsUpdate,
} from '@/lib/api/settings';
import { SettingsSection } from './components/SettingsSection';
import { ProfileSettings } from './components/ProfileSettings';
import { NotificationSettings } from './components/NotificationSettings';
import { Integrations } from './components/Integrations';
import { APISettings } from './components/APISettings';
import { AlertThresholds } from './components/AlertThresholds';
import {
  User,
  Bell,
  Plug,
  Key,
  AlertTriangle,
  Settings as SettingsIcon,
} from 'lucide-react';

type TabId = 'profile' | 'notifications' | 'integrations' | 'api' | 'alerts';

interface Tab {
  id: TabId;
  label: string;
  icon: React.ReactNode;
}

const TABS: Tab[] = [
  { id: 'profile', label: 'Profile', icon: <User className="w-4 h-4" /> },
  { id: 'notifications', label: 'Notifications', icon: <Bell className="w-4 h-4" /> },
  { id: 'integrations', label: 'Integrations', icon: <Plug className="w-4 h-4" /> },
  { id: 'api', label: 'API Access', icon: <Key className="w-4 h-4" /> },
  { id: 'alerts', label: 'Alert Thresholds', icon: <AlertTriangle className="w-4 h-4" /> },
];

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<TabId>('profile');
  const [settings, setSettings] = useState<UserSettingsType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadSettings();
  }, []);

  const loadSettings = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await getUserSettings();
      setSettings(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load settings');
      console.error('Failed to load settings:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = async (update: UserSettingsUpdate) => {
    try {
      const updated = await updateUserSettings(update);
      setSettings(updated);
      alert('Settings updated successfully!');
    } catch (err) {
      console.error('Failed to update settings:', err);
      alert('Failed to update settings');
      throw err;
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <SettingsIcon className="w-12 h-12 text-gray-400 mx-auto mb-4 animate-spin" />
          <div className="text-lg text-gray-600 dark:text-slate-400">Loading settings...</div>
        </div>
      </div>
    );
  }

  if (error || !settings) {
    return (
      <div className="max-w-2xl mx-auto">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-red-900 dark:text-red-100 mb-2">
            Error Loading Settings
          </h3>
          <p className="text-red-700 dark:text-red-300">{error}</p>
          <button
            onClick={loadSettings}
            className="mt-4 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-slate-100">Settings</h1>
        <p className="text-gray-600 dark:text-slate-400 mt-1">
          Manage your account settings, integrations, and preferences
        </p>
      </div>

      {/* Tabs Navigation */}
      <div className="bg-white dark:bg-slate-800 rounded-lg shadow-sm border border-gray-200 dark:border-slate-700 overflow-hidden">
        <div className="border-b border-gray-200 dark:border-slate-700">
          <div className="flex overflow-x-auto">
            {TABS.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 px-6 py-4 text-sm font-medium border-b-2 transition-colors whitespace-nowrap ${
                  activeTab === tab.id
                    ? 'border-blue-600 text-blue-600 dark:text-blue-400'
                    : 'border-transparent text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200 hover:border-gray-300 dark:hover:border-slate-600'
                }`}
              >
                {tab.icon}
                {tab.label}
              </button>
            ))}
          </div>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {activeTab === 'profile' && (
            <SettingsSection
              title="Profile Information"
              description="Update your personal information and preferences"
            >
              <ProfileSettings settings={settings} onUpdate={handleUpdate} />
            </SettingsSection>
          )}

          {activeTab === 'notifications' && (
            <SettingsSection
              title="Notification Preferences"
              description="Choose how you want to receive alerts and notifications"
            >
              <NotificationSettings settings={settings} onUpdate={handleUpdate} />
            </SettingsSection>
          )}

          {activeTab === 'integrations' && (
            <SettingsSection
              title="Third-Party Integrations"
              description="Connect and manage external services"
            >
              <Integrations />
            </SettingsSection>
          )}

          {activeTab === 'api' && (
            <SettingsSection
              title="API Access"
              description="Manage your API credentials and access"
            >
              <APISettings settings={settings} onUpdate={loadSettings} />
            </SettingsSection>
          )}

          {activeTab === 'alerts' && (
            <SettingsSection
              title="Alert Thresholds"
              description="Configure when to receive different types of alerts"
            >
              <AlertThresholds settings={settings} onUpdate={handleUpdate} />
            </SettingsSection>
          )}
        </div>
      </div>

      {/* System Info Footer */}
      <div className="bg-gray-50 dark:bg-slate-800/50 rounded-lg p-4 border border-gray-200 dark:border-slate-700">
        <div className="flex items-center justify-between text-sm text-gray-600 dark:text-slate-400">
          <div>
            <span className="font-medium">Account Created:</span>{' '}
            {new Date(settings.created_at).toLocaleDateString()}
          </div>
          <div>
            <span className="font-medium">Last Updated:</span>{' '}
            {new Date(settings.updated_at).toLocaleDateString()}
          </div>
        </div>
      </div>
    </div>
  );
}