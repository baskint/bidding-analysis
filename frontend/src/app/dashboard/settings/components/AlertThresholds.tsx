// frontend/src/app/dashboard/settings/components/AlertThresholds.tsx
'use client';

import React, { useState } from 'react';
import { UserSettings, UserSettingsUpdate } from '@/lib/api/settings';
import { AlertTriangle, TrendingDown, Shield } from 'lucide-react';

interface AlertThresholdsProps {
  settings: UserSettings;
  onUpdate: (update: UserSettingsUpdate) => Promise<void>;
}

export function AlertThresholds({ settings, onUpdate }: AlertThresholdsProps) {
  const [fraudThreshold, setFraudThreshold] = useState(settings.fraud_alert_threshold);
  const [budgetThreshold, setBudgetThreshold] = useState(settings.budget_alert_threshold);
  const [performanceThreshold, setPerformanceThreshold] = useState(settings.performance_alert_threshold);
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    setSaving(true);
    try {
      await onUpdate({
        fraud_alert_threshold: fraudThreshold,
        budget_alert_threshold: budgetThreshold,
        performance_alert_threshold: performanceThreshold,
      });
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <p className="text-gray-600 dark:text-slate-400">
        Configure when to receive alerts based on these threshold values
      </p>

      {/* Fraud Alert Threshold */}
      <div>
        <div className="flex items-center justify-between mb-2">
          <label className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-slate-300">
            <Shield className="w-4 h-4" />
            Fraud Detection Threshold
          </label>
          <span className="text-lg font-semibold text-gray-900 dark:text-slate-100">
            {(fraudThreshold * 100).toFixed(0)}%
          </span>
        </div>
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          value={fraudThreshold}
          onChange={(e) => setFraudThreshold(parseFloat(e.target.value))}
          className="w-full h-2 bg-gray-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer accent-red-600"
        />
        <div className="flex justify-between text-xs text-gray-500 dark:text-slate-400 mt-1">
          <span>Low Sensitivity (0%)</span>
          <span>High Sensitivity (100%)</span>
        </div>
        <p className="text-sm text-gray-600 dark:text-slate-400 mt-2">
          Alert when fraud probability exceeds {(fraudThreshold * 100).toFixed(0)}%
        </p>
      </div>

      {/* Budget Alert Threshold */}
      <div className="pt-4 border-t border-gray-200 dark:border-slate-700">
        <div className="flex items-center justify-between mb-2">
          <label className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-slate-300">
            <TrendingDown className="w-4 h-4" />
            Budget Alert Threshold
          </label>
          <span className="text-lg font-semibold text-gray-900 dark:text-slate-100">
            {(budgetThreshold * 100).toFixed(0)}%
          </span>
        </div>
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          value={budgetThreshold}
          onChange={(e) => setBudgetThreshold(parseFloat(e.target.value))}
          className="w-full h-2 bg-gray-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer accent-yellow-600"
        />
        <div className="flex justify-between text-xs text-gray-500 dark:text-slate-400 mt-1">
          <span>Early Warning (0%)</span>
          <span>Late Warning (100%)</span>
        </div>
        <p className="text-sm text-gray-600 dark:text-slate-400 mt-2">
          Alert when budget consumption reaches {(budgetThreshold * 100).toFixed(0)}%
        </p>
      </div>

      {/* Performance Alert Threshold */}
      <div className="pt-4 border-t border-gray-200 dark:border-slate-700">
        <div className="flex items-center justify-between mb-2">
          <label className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-slate-300">
            <AlertTriangle className="w-4 h-4" />
            Performance Alert Threshold
          </label>
          <span className="text-lg font-semibold text-gray-900 dark:text-slate-100">
            {(performanceThreshold * 100).toFixed(0)}%
          </span>
        </div>
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          value={performanceThreshold}
          onChange={(e) => setPerformanceThreshold(parseFloat(e.target.value))}
          className="w-full h-2 bg-gray-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer accent-blue-600"
        />
        <div className="flex justify-between text-xs text-gray-500 dark:text-slate-400 mt-1">
          <span>Lenient (0%)</span>
          <span>Strict (100%)</span>
        </div>
        <p className="text-sm text-gray-600 dark:text-slate-400 mt-2">
          Alert when campaign performance drops below {(performanceThreshold * 100).toFixed(0)}% of target
        </p>
      </div>

      {/* Threshold Examples */}
      <div className="bg-gray-50 dark:bg-slate-800/50 rounded-lg p-4 mt-6">
        <h4 className="font-medium text-gray-900 dark:text-slate-100 mb-3">Alert Examples</h4>
        <div className="space-y-2 text-sm">
          <div className="flex items-start gap-2">
            <Shield className="w-4 h-4 text-red-600 mt-0.5" />
            <div>
              <span className="text-gray-900 dark:text-slate-100 font-medium">Fraud:</span>
              <span className="text-gray-600 dark:text-slate-400 ml-1">
                Notify when suspicious activity has {(fraudThreshold * 100).toFixed(0)}% confidence
              </span>
            </div>
          </div>
          <div className="flex items-start gap-2">
            <TrendingDown className="w-4 h-4 text-yellow-600 mt-0.5" />
            <div>
              <span className="text-gray-900 dark:text-slate-100 font-medium">Budget:</span>
              <span className="text-gray-600 dark:text-slate-400 ml-1">
                Notify when ${(budgetThreshold * 100).toFixed(0)} of $100 budget is spent
              </span>
            </div>
          </div>
          <div className="flex items-start gap-2">
            <AlertTriangle className="w-4 h-4 text-blue-600 mt-0.5" />
            <div>
              <span className="text-gray-900 dark:text-slate-100 font-medium">Performance:</span>
              <span className="text-gray-600 dark:text-slate-400 ml-1">
                Notify when conversion rate is {(performanceThreshold * 100).toFixed(0)}% below target
              </span>
            </div>
          </div>
        </div>
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
