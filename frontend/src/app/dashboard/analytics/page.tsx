// frontend/src/app/dashboard/analytics/page.tsx
'use client';

import { useState, useEffect, useCallback } from "react";
import {
  BarChart3,
  Globe,
  Clock,
  Calendar,
  Download,
} from "lucide-react";
import {
  getPerformanceOverview,
  getKeywordAnalysis,
  getDeviceBreakdown,
  getGeoBreakdown,
  getHourlyPerformance,
  getCompetitiveAnalysis,
  type PerformanceMetrics,
  type KeywordAnalysis,
  type DeviceBreakdown,
  type GeoBreakdown,
  type HourlyPerformance,
  type CompetitiveAnalysis,
} from "@/lib/api/analytics";
import { PerformanceOverviewCard } from "./components/PerformanceOverviewCard";
import { KeywordAnalysisTable } from "./components/KeywordAnalysisTable";
import { DeviceBreakdownChart } from "./components/DeviceBreakdownChart";
import { getDateRange, formatPercent, formatCurrency } from "@/lib/utils";

// Date range selector component
function DateRangeSelector({
  value,
  onChange,
}: {
  value: string;
  onChange: (value: string) => void;
}) {
  return (
    <div className="flex items-center space-x-3">
      <Calendar className="w-5 h-5 text-slate-600" />
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="px-4 py-2 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
      >
        <option value="7">Last 7 Days</option>
        <option value="14">Last 14 Days</option>
        <option value="30">Last 30 Days</option>
        <option value="90">Last 90 Days</option>
      </select>
    </div>
  );
}

// Geographic Breakdown Component
function GeoBreakdownTable({
  geos,
  loading,
}: {
  geos: GeoBreakdown[];
  loading: boolean;
}) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/3"></div>
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-12 bg-slate-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (geos.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 text-center text-slate-600">
        No geographic data available
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <Globe className="w-6 h-6 text-emerald-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Geographic Performance
        </h2>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="text-left text-sm text-slate-600 border-b border-slate-200">
              <th className="pb-3 font-semibold">Location</th>
              <th className="pb-3 font-semibold text-right">Bids</th>
              <th className="pb-3 font-semibold text-right">Won</th>
              <th className="pb-3 font-semibold text-right">Win Rate</th>
              <th className="pb-3 font-semibold text-right">Conv. Rate</th>
              <th className="pb-3 font-semibold text-right">Spend</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {geos && geos.length > 0 && geos.map((geo, idx) => (
              <tr key={idx} className="hover:bg-slate-50">
                <td className="py-3">
                  <div>
                    <div className="font-medium text-slate-900">
                      {geo.country}
                    </div>
                    {geo.region && (
                      <div className="text-sm text-slate-500">
                        {geo.region}
                      </div>
                    )}
                  </div>
                </td>
                <td className="py-3 text-right text-slate-600">
                  {(geo?.totalBids || 0).toLocaleString()}
                </td>
                <td className="py-3 text-right text-slate-600">
                  {(geo?.wonBids || 0).toLocaleString()}
                </td>
                <td className="py-3 text-right">
                  <span className="inline-flex items-center px-2 py-1 rounded-full bg-blue-100 text-blue-700 text-xs font-medium">
                    {formatPercent(geo.winRate)}
                  </span>
                </td>
                <td className="py-3 text-right">
                  <span className="inline-flex items-center px-2 py-1 rounded-full bg-green-100 text-green-700 text-xs font-medium">
                    {formatPercent(geo.conversionRate)}
                  </span>
                </td>
                <td className="py-3 text-right text-slate-900 font-medium">
                  {formatCurrency(geo.spend)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

// Hourly Performance Chart
function HourlyPerformanceChart({
  hourly,
  loading,
}: {
  hourly: HourlyPerformance[];
  loading: boolean;
}) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/3"></div>
          <div className="h-48 bg-slate-200 rounded"></div>
        </div>
      </div>
    );
  }

  if (hourly.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 text-center text-slate-600">
        No hourly data available
      </div>
    );
  }

  const maxBids = (hourly && hourly.length > 0) ? Math.max(...hourly.map((h) => h.total_bids), 1) : 0;

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <Clock className="w-6 h-6 text-orange-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Performance by Hour of Day
        </h2>
      </div>

      <div className="grid grid-cols-12 gap-2">
        {hourly && hourly.length > 0 && hourly.map((h) => (
          <div key={h.hour} className="flex flex-col items-center space-y-1">
            <div
              className="w-full bg-gradient-to-t from-orange-500 to-yellow-400 rounded-t"
              style={{
                height: `${Math.max((h.total_bids / maxBids) * 120, 4)}px`,
              }}
              title={`${h.total_bids} bids at ${h.hour}:00`}
            ></div>
            <span className="text-xs text-slate-600">{h.hour}</span>
          </div>
        ))}
      </div>

      <div className="mt-6 grid grid-cols-3 gap-4 text-center text-sm">
        <div>
          <p className="text-slate-600">Best Hour</p>
          <p className="text-lg font-bold text-slate-900">
            {hourly && hourly.length > 0 && hourly.reduce((max, h) =>
              h.total_bids > max.total_bids ? h : max,
            ).hour}
            :00
          </p>
        </div>
        <div>
          <p className="text-slate-600">Peak Win Rate</p>
          <p className="text-lg font-bold text-blue-600">
            {hourly && hourly.length > 0 && formatPercent(
              hourly.reduce((max, h) =>
                h.win_rate > max.win_rate ? h : max,
              ).win_rate
            )}
          </p>
        </div>
        <div>
          <p className="text-slate-600">Peak Conv. Rate</p>
          <p className="text-lg font-bold text-green-600">
            {hourly && hourly.length > 0 && formatPercent(
              hourly.reduce((max, h) =>
                h.conversion_rate > max.conversion_rate ? h : max,
              ).conversion_rate
            )}
          </p>
        </div>
      </div>
    </div>
  );
}

// Competitive Analysis Component
function CompetitiveAnalysisCard({
  competitive,
  loading,
}: {
  competitive: CompetitiveAnalysis[];
  loading: boolean;
}) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/3"></div>
          <div className="space-y-3">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-16 bg-slate-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (competitive.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 text-center text-slate-600">
        No competitive data available
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <BarChart3 className="w-6 h-6 text-rose-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Competitive Landscape
        </h2>
      </div>

      <div className="space-y-4">
        {competitive.length > 0 && competitive.slice(0, 5).map((comp, idx) => (
          <div
            key={idx}
            className="p-4 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors"
          >
            <div className="flex items-center justify-between mb-2">
              <span className="font-medium text-slate-900">
                {comp.segment_category}
              </span>
              <span className="text-sm text-slate-600">
                {(comp?.total_opportunities || 0).toLocaleString()} opportunities
              </span>
            </div>

            <div className="grid grid-cols-4 gap-4 text-sm">
              <div>
                <p className="text-slate-600">Our Win Rate</p>
                <p className="font-semibold text-slate-900">
                  {formatPercent(comp.our_win_rate)}
                </p>
              </div>
              <div>
                <p className="text-slate-600">Our Avg Bid</p>
                <p className="font-semibold text-slate-900">
                  {formatCurrency(comp.our_average_bid)}
                </p>
              </div>
              <div>
                <p className="text-slate-600">Market Avg</p>
                <p className="font-semibold text-slate-900">
                  {formatCurrency(comp.market_average_bid)}
                </p>
              </div>
              <div>
                <p className="text-slate-600">Competition</p>
                <p
                  className={`font-semibold ${comp.competition_intensity > 1.5
                      ? "text-red-600"
                      : comp.competition_intensity > 1.2
                        ? "text-yellow-600"
                        : "text-green-600"
                    }`}
                >
                  {(comp?.competition_intensity || 0).toFixed(2)}x
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// Main Analysis Page
export default function AnalysisPage() {
  const [dateRange, setDateRange] = useState("30");
  const [loading, setLoading] = useState(true);

  // State for all data
  const [performanceMetrics, setPerformanceMetrics] =
    useState<PerformanceMetrics | null>(null);
  const [keywords, setKeywords] = useState<KeywordAnalysis[]>([]);
  const [devices, setDevices] = useState<DeviceBreakdown[]>([]);
  const [geos, setGeos] = useState<GeoBreakdown[]>([]);
  const [hourly, setHourly] = useState<HourlyPerformance[]>([]);
  const [competitive, setCompetitive] = useState<CompetitiveAnalysis[]>([]);

  const loadData = useCallback(async () => {
    setLoading(true);
    const dateParams = getDateRange(parseInt(dateRange));

    try {
      // Using Promise.all for parallel fetching - ALL real API calls
      const [
        perfData,
        kwData,
        deviceData,
        geoData,
        hourlyData,
        competitiveData,
      ] = await Promise.all([
        getPerformanceOverview(dateParams),
        getKeywordAnalysis({ ...dateParams, limit: 20 }),
        getDeviceBreakdown(dateParams),
        getGeoBreakdown(dateParams),
        getHourlyPerformance(dateParams),
        getCompetitiveAnalysis(dateParams),
      ]);

      setPerformanceMetrics(perfData);
      setKeywords(kwData || []);
      setDevices(deviceData || []);
      setGeos(geoData || []);
      setHourly(hourlyData || []);
      setCompetitive(competitiveData || []);

    } catch (error) {
      console.error("Failed to load analysis data:", error);
    } finally {
      setLoading(false);
    }
  }, [dateRange]);

  // Load data on mount and when date range changes
  useEffect(() => {
    loadData();
  }, [dateRange, loadData]);

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-100">
            Campaign Analytics
          </h1>
          <p className="text-slate-600 mt-1">
            Deep dive into your bidding performance and insights
          </p>
        </div>
        <div className="flex items-center space-x-3">
          <DateRangeSelector value={dateRange} onChange={setDateRange} />
          <button className="flex items-center px-4 py-2 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors">
            <Download className="w-4 h-4 mr-2" />
            Export Report
          </button>
        </div>
      </div>

      {/* Performance Overview */}
      <PerformanceOverviewCard
        metrics={performanceMetrics}
        loading={loading}
      />

      {/* Keyword Analysis */}
      <KeywordAnalysisTable keywords={keywords} loading={loading} />

      {/* Device & Hourly Performance Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <DeviceBreakdownChart devices={devices} loading={loading} />
        <HourlyPerformanceChart hourly={hourly} loading={loading} />
      </div>

      {/* Geographic & Competitive Analysis Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <GeoBreakdownTable geos={geos} loading={loading} />
        <CompetitiveAnalysisCard
          competitive={competitive}
          loading={loading}
        />
      </div>
    </div>
  );
}
