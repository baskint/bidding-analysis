// frontend/src/app/dashboard/settings/components/SettingsSection.tsx
import React from 'react';

interface SettingsSectionProps {
  title: string;
  description?: string;
  children: React.ReactNode;
}

export function SettingsSection({ title, description, children }: SettingsSectionProps) {
  return (
    <div className="bg-white dark:bg-slate-800 rounded-lg shadow-sm border border-gray-200 dark:border-slate-700 overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-200 dark:border-slate-700">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-slate-100">{title}</h3>
        {description && (
          <p className="text-sm text-gray-600 dark:text-slate-400 mt-1">{description}</p>
        )}
      </div>
      <div className="px-6 py-4">{children}</div>
    </div>
  );
}
