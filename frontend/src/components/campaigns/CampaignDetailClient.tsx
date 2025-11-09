// frontend/src/components/campaigns/CampaignDetailClient.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import {
  ArrowLeft,
  Edit,
  Pause,
  Play,
  TrendingUp,
  DollarSign,
  Target,
  Activity,
  Calendar,
  Globe,
  Smartphone,
  Hash,
} from "lucide-react";
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import {
  getCampaign,
  pauseCampaign,
  activateCampaign,
  CampaignDetail,
} from "@/lib/api/campaigns";

// ---
// 🎯 PROP DEFINITION
// The Server Component now passes the fully fetched data as a prop.
interface CampaignDetailClientProps {
  initialCampaign: CampaignDetail;
}

const COLORS = ["#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"];

// ---
// 🚀 COMPONENT START
export default function CampaignDetailClient({
  initialCampaign,
}: CampaignDetailClientProps) {
  const router = useRouter();

  // 1. STATE INITIALIZATION: Initialize the state directly with the pre-fetched data.
  const [campaign, setCampaign] = useState<CampaignDetail>(initialCampaign);

  // 2. STATE FOR ACTIONS: We only need 'loading' and 'error' states for user actions (pause/activate/reload).
  const [isActionLoading, setIsActionLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // ---
  // 🔄 REFRESH DATA LOGIC
  // This function is now only responsible for reloading the campaign *after* an action.
  const loadCampaign = async () => {
    setIsActionLoading(true);
    try {
      // Use the ID from the current campaign state
      const data = await getCampaign(campaign.id);
      setCampaign(data);
      setError(null);
    } catch (err) {
      setError("Failed to reload campaign data after action.");
      console.error(err);
    } finally {
      setIsActionLoading(false);
    }
  };

  // ---
  // ⏸️ ACTION HANDLERS
  const handlePause = async () => {
    if (isActionLoading) return;
    try {
      await pauseCampaign(campaign.id);
      await loadCampaign();
    } catch (err) {
      console.error("Failed to pause campaign:", err);
      alert("Failed to pause campaign");
    }
  };

  const handleActivate = async () => {
    if (isActionLoading) return;
    try {
      await activateCampaign(campaign.id);
      await loadCampaign();
    } catch (err) {
      console.error("Failed to activate campaign:", err);
      alert("Failed to activate campaign");
    }
  };

  // ---
  // 🖼️ FORMATTERS (moved lower for better readability)
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const formatPercent = (value: number) => {
    return `${(value * 100).toFixed(1)}%`;
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "active":
        return "bg-green-100 text-green-800";
      case "paused":
        return "bg-yellow-100 text-yellow-800";
      case "archived":
        return "bg-gray-100 text-gray-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  // ---
  // ⏳ ACTION LOADING OVERLAY
  // Show a simpler loading state if the user performs an action (Pause/Activate)
  if (isActionLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Updating campaign status...</p>
        </div>
      </div>
    );
  }

  // NOTE: Initial load/404 handling should be done in the parent Server Component.
  // We only show error here if a subsequent action failed.
  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <p className="text-red-600">Error: {error}</p>
          <button
            onClick={() => router.push("/dashboard/campaigns")}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Back to Campaigns
          </button>
        </div>
      </div>
    );
  }

  // ---
  // 📊 MAIN UI RENDER
  return (
    <div className="p-8">
      {/* Header */}
      <div className="mb-8">
        <Link
          href="/dashboard/campaigns"
          className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back to Campaigns
        </Link>

        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">
              {campaign.name}
            </h1>
            <span
              className={`inline-block mt-2 px-3 py-1 text-sm font-medium rounded-full ${getStatusColor(campaign.status)}`}
            >
              {campaign.status}
            </span>
          </div>
          <div className="flex gap-2">
            <Link
              href={`/dashboard/campaigns/${campaign.id}/edit`}
              className="flex items-center px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50"
            >
              <Edit className="w-4 h-4 mr-2" />
              Edit
            </Link>
            {campaign.status === "active" ? (
              <button
                onClick={handlePause}
                className="flex items-center px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 disabled:opacity-50"
                disabled={isActionLoading}
              >
                <Pause className="w-4 h-4 mr-2" />
                Pause
              </button>
            ) : campaign.status === "paused" ? (
              <button
                onClick={handleActivate}
                className="flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
                disabled={isActionLoading}
              >
                <Play className="w-4 h-4 mr-2" />
                Activate
              </button>
            ) : null}
          </div>
        </div>
      </div>

      {/* Key Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Total Spend</p>
              <p className="text-2xl font-bold text-gray-900">
                {formatCurrency(campaign.total_spend)}
              </p>
              {campaign.budget && (
                <p className="text-xs text-gray-500 mt-1">
                  of {formatCurrency(campaign.budget)}
                </p>
              )}
            </div>
            <DollarSign className="w-8 h-8 text-blue-600" />
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Win Rate</p>
              <p className="text-2xl font-bold text-gray-900">
                {formatPercent(campaign.win_rate)}
              </p>
              <p className="text-xs text-gray-500 mt-1">
                {campaign.won_bids} / {campaign.total_bids} bids
              </p>
            </div>
            <Activity className="w-8 h-8 text-green-600" />
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Conversions</p>
              <p className="text-2xl font-bold text-gray-900">
                {campaign.conversions}
              </p>
              <p className="text-xs text-gray-500 mt-1">
                {formatPercent(campaign.conversion_rate)} rate
              </p>
            </div>
            <Target className="w-8 h-8 text-purple-600" />
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Cost / Conversion</p>
              <p className="text-2xl font-bold text-gray-900">
                {campaign.conversions > 0
                  ? formatCurrency(campaign.cost_per_conversion)
                  : "N/A"}
              </p>
              {campaign.target_cpa && (
                <p className="text-xs text-gray-500 mt-1">
                  target: {formatCurrency(campaign.target_cpa)}
                </p>
              )}
            </div>
            <TrendingUp className="w-8 h-8 text-orange-600" />
          </div>
        </div>
      </div>

      {/* Daily Performance Chart */}
      {campaign.daily_metrics && campaign.daily_metrics.length > 0 && (
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200 mb-8">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Daily Performance
          </h2>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={campaign.daily_metrics}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="date"
                tickFormatter={(date) =>
                  new Date(date).toLocaleDateString("en-US", {
                    month: "short",
                    day: "numeric",
                  })
                }
              />
              <YAxis yAxisId="left" />
              <YAxis yAxisId="right" orientation="right" />
              <Tooltip
                labelFormatter={(date) => new Date(date).toLocaleDateString()}
                formatter={(value: number, name: string) => [
                  // Ensure currency/percentage formatting for tooltips if needed
                  name === "Win Rate"
                    ? formatPercent(value)
                    : name === "Spend"
                      ? formatCurrency(value)
                      : value.toFixed(0),
                  name,
                ]}
              />
              <Legend />
              <Line
                yAxisId="left"
                type="monotone"
                dataKey="total_bids"
                stroke="#3b82f6"
                name="Total Bids"
              />
              <Line
                yAxisId="left"
                dataKey="conversions"
                stroke="#10b981"
                name="Conversions"
              />
              <Line
                yAxisId="right"
                dataKey="win_rate"
                stroke="#f59e0b"
                name="Win Rate"
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      )}

      {/* Device & Geo Breakdown */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Device Breakdown */}
        {campaign.device_breakdown && campaign.device_breakdown.length > 0 && (
          <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
              <Smartphone className="w-5 h-5 mr-2" />
              Device Performance
            </h2>
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={campaign.device_breakdown}
                  dataKey="bids"
                  nameKey="device_type"
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  label={(entry: { device_type: string; bids: number }) => 
                    `${entry.device_type}: ${entry.bids}`
                  }
                >
                  {campaign.device_breakdown.map((entry, index) => (
                    <Cell
                      key={`cell-${index}`}
                      fill={COLORS[index % COLORS.length]}
                    />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        )}

        {/* Geographic Breakdown */}
        {campaign.geo_breakdown && campaign.geo_breakdown.length > 0 && (
          <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
              <Globe className="w-5 h-5 mr-2" />
              Geographic Performance
            </h2>
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={campaign.geo_breakdown}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="country" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="bids" fill="#3b82f6" name="Bids" />
                <Bar dataKey="conversions" fill="#10b981" name="Conversions" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        )}
      </div>

      {/* Top Keywords */}
      {campaign.top_keywords && campaign.top_keywords.length > 0 && (
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200 mb-8">
          <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
            <Hash className="w-5 h-5 mr-2" />
            Top Keywords
          </h2>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead>
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Keyword
                  </th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    Bids
                  </th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    Won
                  </th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    Conversions
                  </th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    Spend
                  </th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    Win Rate
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {campaign.top_keywords.map((keyword, index) => (
                  <tr key={index} className="hover:bg-gray-50">
                    <td className="px-4 py-3 text-sm font-medium text-gray-900">
                      {keyword.keyword}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 text-right">
                      {keyword.bids}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 text-right">
                      {keyword.won_bids}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 text-right">
                      {keyword.conversions}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 text-right">
                      {formatCurrency(keyword.spend)}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 text-right">
                      {formatPercent(keyword.win_rate)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Recent Bids */}
      {campaign.recent_bids && campaign.recent_bids.length > 0 && (
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
            <Calendar className="w-5 h-5 mr-2" />
            Recent Bids
          </h2>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead>
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Time
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Bid Price
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Result
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Device
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Location
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {campaign.recent_bids.map((bid) => (
                  <tr key={bid.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {new Date(bid.timestamp).toLocaleString()}
                    </td>
                    <td className="px-4 py-3 text-sm font-medium text-gray-900">
                      {formatCurrency(bid.bid_price)}
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <span
                        className={`px-2 py-1 rounded-full text-xs font-medium ${
                          bid.won
                            ? "bg-green-100 text-green-800"
                            : "bg-red-100 text-red-800"
                        }`}
                      >
                        {bid.won ? "Won" : "Lost"}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {bid.device_type || "N/A"}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {bid.country
                        ? `${bid.city || ""}, ${bid.country}`
                        : "N/A"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}
