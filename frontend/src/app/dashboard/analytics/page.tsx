// frontend/src/app/dashboard/analytics/page.tsx
"use client";

import { useState, useEffect } from "react";
import {
  BarChart3,
  TrendingUp,
  Globe,
  Smartphone,
  Clock,
  Target,
  Calendar,
  Download,
} from "lucide-react";
import {
  getPerformanceOverview,
  // getKeywordAnalysis,
  // getDeviceBreakdown,
  // getGeoBreakdown,
  // getHourlyPerformance,
  // getCompetitiveAnalysis,
  type PerformanceMetrics,
  type KeywordAnalysis,
  type DeviceBreakdown,
  type GeoBreakdown,
  type HourlyPerformance,
  type CompetitiveAnalysis,
} from "@/lib/api/analytics";

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

// Performance Overview Card
function PerformanceOverviewCard({
  metrics,
  loading,
}: {
  metrics: PerformanceMetrics | null;
  loading: boolean;
}) {
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
  const formatPercent = (value: number) =>
    `${(value * 100).toFixed(2)}%`;
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

// Keyword Analysis Table
function KeywordAnalysisTable({
  keywords,
  loading,
}: {
  keywords: KeywordAnalysis[];
  loading: boolean;
}) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/4"></div>
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-12 bg-slate-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (keywords.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 text-center text-slate-600">
        No keyword data available
      </div>
    );
  }

  const formatCurrency = (value: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);
  const formatPercent = (value: number) =>
    `${(value * 100).toFixed(1)}%`;

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <Target className="w-6 h-6 text-purple-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Top Keywords by Performance
        </h2>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="text-left text-sm text-slate-600 border-b border-slate-200">
              <th className="pb-3 font-semibold">Keyword</th>
              <th className="pb-3 font-semibold text-right">Bids</th>
              <th className="pb-3 font-semibold text-right">Win Rate</th>
              <th className="pb-3 font-semibold text-right">Conv. Rate</th>
              <th className="pb-3 font-semibold text-right">Spend</th>
              <th className="pb-3 font-semibold text-right">ROAS</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {keywords.map((kw, idx) => (
              <tr key={idx} className="hover:bg-slate-50">
                <td className="py-3 font-medium text-slate-900">
                  {kw.keyword}
                </td>
                <td className="py-3 text-right text-slate-600">
                  {/* FIX: Use nullish coalescing (?? 0) */}
                  {(kw.total_bids ?? 0).toLocaleString()}
                </td>
                <td className="py-3 text-right">
                  <span className="inline-flex items-center px-2 py-1 rounded-full bg-blue-100 text-blue-700 text-xs font-medium">
                    {/* FIX: Use nullish coalescing (?? 0) */}
                    {formatPercent(kw.win_rate ?? 0)}
                  </span>
                </td>
                <td className="py-3 text-right">
                  <span className="inline-flex items-center px-2 py-1 rounded-full bg-green-100 text-green-700 text-xs font-medium">
                    {/* FIX: Use nullish coalescing (?? 0) */}
                    {formatPercent(kw.conversion_rate ?? 0)}
                  </span>
                </td>
                <td className="py-3 text-right text-slate-900 font-medium">
                  {/* FIX: Use nullish coalescing (?? 0) */}
                  {formatCurrency(kw.spend ?? 0)}
                </td>
                <td
                  className={`py-3 text-right font-semibold ${
                    // FIX: Use nullish coalescing (?? 0) for comparison
                    (kw.roas ?? 0) >= 2
                      ? "text-green-600"
                      : (kw.roas ?? 0) >= 1
                        ? "text-yellow-600"
                        : "text-red-600"
                    }`}
                >
                  {/* FIX: Use nullish coalescing (?? 0) to prevent .toFixed error */}
                  {(kw.roas ?? 0).toFixed(2)}x
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

// Device Breakdown Component
function DeviceBreakdownChart({
  devices,
  loading,
}: {
  devices: DeviceBreakdown[];
  loading: boolean;
}) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/3"></div>
          <div className="space-y-3">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-16 bg-slate-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (devices.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 text-center text-slate-600">
        No device data available
      </div>
    );
  }

  const formatPercent = (value: number) =>
    `${(value * 100).toFixed(1)}%`;
  const maxBids = Math.max(...devices.map((d) => d.total_bids), 1);

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <Smartphone className="w-6 h-6 text-indigo-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Device Breakdown
        </h2>
      </div>

      <div className="space-y-4">
        {devices.map((device, idx) => (
          <div key={idx} className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-slate-900">
                {device.device_type}
              </span>
              <div className="flex items-center space-x-4 text-sm">
                <span className="text-slate-600">
                  {(device?.total_bids ?? 0).toLocaleString()} bids
                </span>
                <span className="text-blue-600 font-medium">
                  {formatPercent(device.win_rate)} WR
                </span>
                <span className="text-green-600 font-medium">
                  {formatPercent(device.conversion_rate)} CR
                </span>
              </div>
            </div>
            <div className="relative h-2 bg-slate-100 rounded-full overflow-hidden">
              <div
                className="absolute top-0 left-0 h-full bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full"
                style={{
                  width: `${(device.total_bids / maxBids) * 100}%`,
                }}
              ></div>
            </div>
          </div>
        ))}
      </div>
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

  const formatPercent = (value: number) =>
    `${(value * 100).toFixed(1)}%`;
  const formatCurrency = (value: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);

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
            {geos.map((geo, idx) => (
              <tr key={idx} className="hover:bg-slate-50">
                <td className="py-3">
                  <div>
                    {/* FIX 1: Safely access the string value of the 'country' object */}
                    <div className="font-medium text-slate-900">
                      {geo.country}
                    </div>

                    {/* FIX 2: Check if the 'region' string value exists before rendering */}
                    {geo.region && (
                      <div className="text-sm text-slate-500">
                        {geo.region}
                      </div>  
                    )}
                  </div>
                </td>

                {/* Corrected metric property names (camelCase as per your data sample) */}
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

  const maxBids = Math.max(...hourly.map((h) => h.total_bids), 1);

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <Clock className="w-6 h-6 text-orange-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Performance by Hour of Day
        </h2>
      </div>

      <div className="grid grid-cols-12 gap-2">
        {hourly.map((h) => (
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
            {hourly.reduce((max, h) =>
              h.total_bids > max.total_bids ? h : max,
            ).hour}
            :00
          </p>
        </div>
        <div>
          <p className="text-slate-600">Peak Win Rate</p>
          <p className="text-lg font-bold text-blue-600">
            {(
              hourly.reduce((max, h) =>
                h.win_rate > max.win_rate ? h : max,
              ).win_rate * 100
            ).toFixed(1)}
            %
          </p>
        </div>
        <div>
          <p className="text-slate-600">Peak Conv. Rate</p>
          <p className="text-lg font-bold text-green-600">
            {(
              hourly.reduce((max, h) =>
                h.conversion_rate > max.conversion_rate ? h : max,
              ).conversion_rate * 100
            ).toFixed(1)}
            %
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

  const formatPercent = (value: number) =>
    `${(value * 100).toFixed(1)}%`;
  const formatCurrency = (value: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <BarChart3 className="w-6 h-6 text-rose-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Competitive Landscape
        </h2>
      </div>

      <div className="space-y-4">
        {competitive.slice(0, 5).map((comp, idx) => (
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
  // const [keywords, setKeywords] = useState<KeywordAnalysis[]>([]);
  // const [devices, setDevices] = useState<DeviceBreakdown[]>([]);
  // const [geos, setGeos] = useState<GeoBreakdown[]>([]);
  // const [hourly, setHourly] = useState<HourlyPerformance[]>([]);
  // const [competitive, setCompetitive] = useState<CompetitiveAnalysis[]>([]);
  const keywords: KeywordAnalysis[] = [];
  const devices: DeviceBreakdown[] = [];
  const geos: GeoBreakdown[] = [];
  const hourly: HourlyPerformance[] = [];
  const competitive: CompetitiveAnalysis[] = [];  

  // Calculate date range
  const getDateRange = () => {
    const end = new Date();
    const start = new Date();
    start.setDate(start.getDate() - parseInt(dateRange));

    return {
      start_date: start.toISOString().split("T")[0],
      end_date: end.toISOString().split("T")[0],
    };
  };

  // Load all data
  const loadData = async () => {
    setLoading(true);
    const dateParams = getDateRange();

    try {
      const [
        perfData,
        // kwData,
        // deviceData,
        // geoData,
        // hourlyData,
        // compData,
      ] = await Promise.all([
        getPerformanceOverview(dateParams),
        // getKeywordAnalysis({ ...dateParams, limit: 20 }),
        // getDeviceBreakdown(dateParams),
        // getGeoBreakdown({ ...dateParams, limit: 20 }),
        // getHourlyPerformance(dateParams),
        // getCompetitiveAnalysis(dateParams),
      ]);

      setPerformanceMetrics(perfData);
      // setKeywords(kwData || []);
      // setDevices(deviceData || []);
      // setGeos(geoData || []);
      // setHourly(hourlyData || []);
      // setCompetitive(compData || []);
    } catch (error) {
      console.error("Failed to load analysis data:", error);
    } finally {
      setLoading(false);
    }
  };

  // Load data on mount and when date range changes
  useEffect(() => {
    loadData();
  }, [dateRange]);

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">
            Campaign Analysis
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
