// frontend/src/app/dashboard/settings/components/ProfileSettings.tsx
'use client';

import React, { useState } from 'react';
import { UserSettings, UserSettingsUpdate } from '@/lib/api/settings';
import { User, Mail, Phone, Globe, Languages } from 'lucide-react';

interface ProfileSettingsProps {
  settings: UserSettings;
  onUpdate: (update: UserSettingsUpdate) => Promise<void>;
}

export function ProfileSettings({ settings, onUpdate }: ProfileSettingsProps) {
  const [fullName, setFullName] = useState(settings.full_name || '');
  const [email, setEmail] = useState(settings.email || '');
  const [phone, setPhone] = useState(settings.phone || '');
  const [timezone, setTimezone] = useState(settings.timezone || 'UTC');
  const [language, setLanguage] = useState(settings.language || 'en');
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    setSaving(true);
    try {
      await onUpdate({
        full_name: fullName,
        email,
        phone,
        timezone,
        language,
      });
    } finally {
      setSaving(false);
    }
  };

  const timezones = [
    'UTC',
    'America/New_York',
    'America/Chicago',
    'America/Denver',
    'America/Los_Angeles',
    'Europe/London',
    'Europe/Paris',
    'Asia/Tokyo',
    'Asia/Shanghai',
    'Australia/Sydney',
  ];

  const languages = [
    { code: 'en', name: 'English' },
    { code: 'es', name: 'Spanish' },
    { code: 'fr', name: 'French' },
    { code: 'de', name: 'German' },
    { code: 'ja', name: 'Japanese' },
    { code: 'zh', name: 'Chinese' },
  ];

  return (
    <div className="space-y-6">
      {/* Full Name */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          <div className="flex items-center gap-2">
            <User className="w-4 h-4" />
            Full Name
          </div>
        </label>
        <input
          type="text"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
          placeholder="John Doe"
        />
      </div>

      {/* Email */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          <div className="flex items-center gap-2">
            <Mail className="w-4 h-4" />
            Email Address
          </div>
        </label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
          placeholder="john@example.com"
        />
      </div>

      {/* Phone */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          <div className="flex items-center gap-2">
            <Phone className="w-4 h-4" />
            Phone Number
          </div>
        </label>
        <input
          type="tel"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
          placeholder="+1 (555) 123-4567"
        />
      </div>

      {/* Timezone */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          <div className="flex items-center gap-2">
            <Globe className="w-4 h-4" />
            Timezone
          </div>
        </label>
        <select
          value={timezone}
          onChange={(e) => setTimezone(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
        >
          {timezones.map((tz) => (
            <option key={tz} value={tz}>
              {tz}
            </option>
          ))}
        </select>
      </div>

      {/* Language */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
          <div className="flex items-center gap-2">
            <Languages className="w-4 h-4" />
            Language
          </div>
        </label>
        <select
          value={language}
          onChange={(e) => setLanguage(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
        >
          {languages.map((lang) => (
            <option key={lang.code} value={lang.code}>
              {lang.name}
            </option>
          ))}
        </select>
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
