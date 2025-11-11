import { TrendingUp } from "lucide-react";
import type { PerformanceMetrics } from "@/lib/api/analytics";

interface PerformanceOverviewCardProps {
  metrics: PerformanceMetrics | null;
  loading: boolean;
}

export function PerformanceOverviewCard({
  metrics,
  loading,
}: PerformanceOverviewCardProps) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/3"></div>
          <div className="grid grid-cols-4 gap-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-20 bg-slate-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!metrics) return null;

  const formatCurrency = (value: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);
  const formatPercent = (value: number) => `${(value * 100).toFixed(2)}%`;
  const formatNumber = (value: number) =>
    new Intl.NumberFormat("en-US").format(value);

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center">
          <TrendingUp className="w-6 h-6 text-blue-600 mr-3" />
          <h2 className="text-xl font-bold text-slate-900">
            Performance Overview
          </h2>
        </div>
        <span className="text-sm text-slate-600">
          Total Bids: {formatNumber(metrics.total_bids)}
        </span>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div className="space-y-1">
          <p className="text-sm text-slate-600">Win Rate</p>
          <p className="text-2xl font-bold text-slate-900">
            {formatPercent(metrics.win_rate)}
          </p>
          <p className="text-xs text-slate-500">
            {formatNumber(metrics.won_bids)} won
          </p>
        </div>

        <div className="space-y-1">
          <p className="text-sm text-slate-600">Conversion Rate</p>
          <p className="text-2xl font-bold text-slate-900">
            {formatPercent(metrics.conversion_rate)}
          </p>
          <p className="text-xs text-slate-500">
            {formatNumber(metrics.conversions)} conversions
          </p>
        </div>

        <div className="space-y-1">
          <p className="text-sm text-slate-600">Total Spend</p>
          <p className="text-2xl font-bold text-slate-900">
            {formatCurrency(metrics.total_spend)}
          </p>
          <p className="text-xs text-slate-500">
            {formatCurrency(metrics.average_bid)} avg bid
          </p>
        </div>

        <div className="space-y-1">
          <p className="text-sm text-slate-600">ROAS</p>
          <p className="text-2xl font-bold text-green-600">
            {metrics.roas.toFixed(2)}x
          </p>
          <p className="text-xs text-slate-500">
            {formatCurrency(metrics.revenue)} revenue
          </p>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6 mt-6 pt-6 border-t border-slate-200">
        <div className="space-y-1">
          <p className="text-sm text-slate-600">Cost Per Acquisition</p>
          <p className="text-xl font-bold text-slate-900">
            {formatCurrency(metrics.cpa)}
          </p>
        </div>

        <div className="space-y-1">
          <p className="text-sm text-slate-600">Revenue</p>
          <p className="text-xl font-bold text-slate-900">
            {formatCurrency(metrics.revenue)}
          </p>
        </div>
      </div>
    </div>
  );
}
