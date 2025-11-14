// frontend/src/app/dashboard/fraud/components/FraudTrendsChart.tsx
'use client';

import { TrendingDown } from "lucide-react";
import { FraudTrend } from "@/lib/api/fraud";

interface FraudTrendsChartProps {
  trends: FraudTrend[];
  loading: boolean;
}

export function FraudTrendsChart({ trends, loading }: FraudTrendsChartProps) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse">
          <div className="h-6 bg-slate-200 rounded w-1/3 mb-4"></div>
          <div className="h-64 bg-slate-200 rounded"></div>
        </div>
      </div>
    );
  }

  if (trends.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="flex items-center mb-4">
          <TrendingDown className="w-5 h-5 text-blue-600 mr-2" />
          <h3 className="text-lg font-semibold text-slate-900">Fraud Trends</h3>
        </div>
        <div className="text-center py-12 text-slate-600">
          <p>No fraud trend data available</p>
        </div>
      </div>
    );
  }

  // Aggregate by date
  const dateMap = new Map<string, { attempts: number; blocked: number; saved: number }>();
  (trends && trends.length > 0) && trends?.forEach((trend) => {
    const existing = dateMap.get(trend.date) || { attempts: 0, blocked: 0, saved: 0 };
    dateMap.set(trend.date, {
      attempts: existing.attempts + trend.fraud_attempts,
      blocked: existing.blocked + trend.blocked_bids,
      saved: existing.saved + trend.amount_saved,
    });
  });

  const sortedDates = Array.from(dateMap.keys()).sort();
  const maxAttempts = Math.max(...Array.from(dateMap.values()).map((v) => v.attempts));

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center">
          <TrendingDown className="w-5 h-5 text-blue-600 mr-2" />
          <h3 className="text-lg font-semibold text-slate-900">Fraud Trends</h3>
        </div>
        <div className="flex items-center space-x-4 text-sm">
          <div className="flex items-center">
            <div className="w-3 h-3 bg-red-500 rounded mr-2"></div>
            <span className="text-slate-600">Fraud Attempts</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 bg-green-500 rounded mr-2"></div>
            <span className="text-slate-600">Blocked</span>
          </div>
        </div>
      </div>

      {/* Simple Bar Chart */}
      <div className="space-y-3">
        {sortedDates.map((date) => {
          const data = dateMap.get(date)!;
          const percentage = (data.attempts / maxAttempts) * 100;
          const blockedPercentage = data.blocked > 0 ? (data.blocked / data.attempts) * percentage : 0;

          return (
            <div key={date} className="space-y-1">
              <div className="flex items-center justify-between text-sm">
                <span className="text-slate-700 font-medium">
                  {new Date(date).toLocaleDateString("en-US", {
                    month: "short",
                    day: "numeric",
                  })}
                </span>
                <div className="flex items-center space-x-3">
                  <span className="text-slate-600">
                    <span className="font-medium text-red-600">{data.attempts}</span> attempts
                  </span>
                  {data.blocked > 0 && (
                    <span className="text-slate-600">
                      <span className="font-medium text-green-600">{data.blocked}</span> blocked
                    </span>
                  )}
                  {data.saved > 0 && (
                    <span className="text-emerald-600 font-medium">
                      ${data.saved.toFixed(2)} saved
                    </span>
                  )}
                </div>
              </div>
              <div className="h-8 bg-slate-100 rounded-lg overflow-hidden relative">
                <div
                  className="absolute left-0 top-0 h-full bg-red-200 transition-all"
                  style={{ width: `${percentage}%` }}
                ></div>
                <div
                  className="absolute left-0 top-0 h-full bg-green-500 transition-all"
                  style={{ width: `${blockedPercentage}%` }}
                ></div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Summary */}
      <div className="mt-6 pt-6 border-t border-slate-200">
        <div className="grid grid-cols-3 gap-4 text-center">
          <div>
            <div className="text-2xl font-bold text-red-600">
              {Array.from(dateMap.values()).reduce((sum, v) => sum + v.attempts, 0)}
            </div>
            <div className="text-sm text-slate-600">Total Attempts</div>
          </div>
          <div>
            <div className="text-2xl font-bold text-green-600">
              {Array.from(dateMap.values()).reduce((sum, v) => sum + v.blocked, 0)}
            </div>
            <div className="text-sm text-slate-600">Blocked</div>
          </div>
          <div>
            <div className="text-2xl font-bold text-emerald-600">
              $
              {Array.from(dateMap.values())
                .reduce((sum, v) => sum + v.saved, 0)
                .toFixed(2)}
            </div>
            <div className="text-sm text-slate-600">Amount Saved</div>
          </div>
        </div>
      </div>
    </div>
  );
}
