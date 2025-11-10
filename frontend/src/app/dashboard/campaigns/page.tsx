// frontend/src/app/dashboard/campaigns/page.tsx
"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import {
  TrendingUp,
  Pause,
  Play,
  Trash2,
  Edit,
  Plus,
  DollarSign,
  Target,
  Activity,
  Calendar,
} from "lucide-react";
import {
  listCampaigns,
  pauseCampaign,
  activateCampaign,
  deleteCampaign,
  CampaignSummary,
} from "@/lib/api/campaigns";

export default function CampaignsPage() {
  const [campaigns, setCampaigns] = useState<CampaignSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<
    "all" | "active" | "paused" | "archived"
  >("all");

  useEffect(() => {
    loadCampaigns();
  }, []);

  const loadCampaigns = async () => {
    try {
      setLoading(true);
      const data = await listCampaigns();
      setCampaigns(data);
      setError(null);
    } catch (err) {
      setError("Failed to load campaigns");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handlePauseCampaign = async (id: string) => {
    try {
      await pauseCampaign(id);
      await loadCampaigns();
    } catch (err) {
      console.error("Failed to pause campaign:", err);
      alert("Failed to pause campaign");
    }
  };

  const handleActivateCampaign = async (id: string) => {
    try {
      await activateCampaign(id);
      await loadCampaigns();
    } catch (err) {
      console.error("Failed to activate campaign:", err);
      alert("Failed to activate campaign");
    }
  };

  const handleDeleteCampaign = async (id: string) => {
    if (!confirm("Are you sure you want to archive this campaign?")) {
      return;
    }

    try {
      await deleteCampaign(id);
      await loadCampaigns();
    } catch (err) {
      console.error("Failed to delete campaign:", err);
      alert("Failed to delete campaign");
    }
  };

  const filteredCampaigns = campaigns.filter((campaign) => {
    const matchesSearch = campaign.name
      .toLowerCase()
      .includes(searchQuery.toLowerCase());
    const matchesStatus =
      statusFilter === "all" || campaign.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

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

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading campaigns...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <p className="text-red-600">{error}</p>
          <button
            onClick={loadCampaigns}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Campaigns</h1>
            <p className="text-gray-600 mt-1">
              Manage your advertising campaigns
            </p>
          </div>
          <Link
            href="/dashboard/campaigns/new"
            className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            New Campaign
          </Link>
        </div>
      </div>

      {/* Filters */}
      <div className="mb-6 flex gap-4">
        <div className="flex-1">
          <input
            type="text"
            placeholder="Search campaigns..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
        {(() => {
          type StatusFilter = "all" | "active" | "paused" | "archived";
          interface StatusOption {
            value: StatusFilter;
            label: string;
          }

          const statusOptions: StatusOption[] = [
            { value: "all", label: "All Statuses" },
            { value: "active", label: "Active" },
            { value: "paused", label: "Paused" },
            { value: "archived", label: "Archived" },
          ];

          return (
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as StatusFilter)}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              {statusOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          );
        })()}
      </div>

      {/* Campaign Grid */}
      {filteredCampaigns.length === 0 ? (
        <div className="text-center py-12">
          <TrendingUp className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            No campaigns found
          </h3>
          <p className="text-gray-600 mb-6">
            {searchQuery || statusFilter !== "all"
              ? "Try adjusting your filters"
              : "Get started by creating your first campaign"}
          </p>
          {!searchQuery && statusFilter === "all" && (
            <Link
              href="/dashboard/campaigns/new"
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              <Plus className="w-5 h-5 mr-2" />
              Create Campaign
            </Link>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCampaigns.map((campaign) => (
            <div
              key={campaign.id}
              className="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow"
            >
              {/* Card Header */}
              <div className="p-6 border-b border-gray-100">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <Link href={`/dashboard/campaigns/${campaign.id}`}>
                      <h3 className="text-lg font-semibold text-gray-900 hover:text-blue-600 cursor-pointer">
                        {campaign.name}
                      </h3>
                    </Link>
                    <span
                      className={`inline-block mt-2 px-2 py-1 text-xs font-medium rounded-full ${getStatusColor(campaign.status)}`}
                    >
                      {campaign.status}
                    </span>
                  </div>
                </div>

                {/* Budget Progress */}
                {campaign.budget && (
                  <div className="mt-4">
                    <div className="flex justify-between text-sm mb-2">
                      <span className="text-gray-600">Budget</span>
                      <span className="font-medium text-gray-900">
                        {formatCurrency(campaign.total_spend)} /{" "}
                        {formatCurrency(campaign.budget)}
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className="bg-blue-600 h-2 rounded-full transition-all"
                        style={{
                          width: `${Math.min(100, (campaign.total_spend / campaign.budget) * 100)}%`,
                        }}
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* Metrics */}
              <div className="p-6 space-y-3">
                <div className="grid grid-cols-2 gap-3">
                  <div className="flex items-center">
                    <Activity className="w-4 h-4 text-blue-600 mr-2" />
                    <div>
                      <div className="text-xs text-gray-600">Win Rate</div>
                      <div className="text-sm font-semibold text-gray-900">
                        {formatPercent(campaign.win_rate)}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center">
                    <Target className="w-4 h-4 text-green-600 mr-2" />
                    <div>
                      <div className="text-xs text-gray-600">Conv. Rate</div>
                      <div className="text-sm font-semibold text-gray-900">
                        {formatPercent(campaign.conversion_rate)}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center">
                    <DollarSign className="w-4 h-4 text-purple-600 mr-2" />
                    <div>
                      <div className="text-xs text-gray-600">Avg Bid</div>
                      <div className="text-sm font-semibold text-gray-900">
                        {formatCurrency(campaign.average_bid)}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center">
                    <TrendingUp className="w-4 h-4 text-orange-600 mr-2" />
                    <div>
                      <div className="text-xs text-gray-600">Cost/Conv</div>
                      <div className="text-sm font-semibold text-gray-900">
                        {campaign.conversions > 0
                          ? formatCurrency(campaign.cost_per_conversion)
                          : "N/A"}
                      </div>
                    </div>
                  </div>
                </div>

                {campaign.last_activity_at && (
                  <div className="flex items-center text-xs text-gray-500 pt-2 border-t border-gray-100">
                    <Calendar className="w-3 h-3 mr-1" />
                    Last activity:{" "}
                    {new Date(campaign.last_activity_at).toLocaleDateString()}
                  </div>
                )}
              </div>

              {/* Actions */}
              <div className="p-4 bg-gray-50 border-t border-gray-100 flex gap-2">
                <Link
                  href={`/dashboard/campaigns/edit?id=${campaign.id}`}
                  className="flex-1 flex items-center justify-center px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-white transition-colors"
                >
                  <Edit className="w-4 h-4 mr-1" />
                  Edit
                </Link>
                {campaign.status === "active" ? (
                  <button
                    onClick={() => handlePauseCampaign(campaign.id)}
                    className="flex-1 flex items-center justify-center px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-white transition-colors"
                  >
                    <Pause className="w-4 h-4 mr-1" />
                    Pause
                  </button>
                ) : campaign.status === "paused" ? (
                  <button
                    onClick={() => handleActivateCampaign(campaign.id)}
                    className="flex-1 flex items-center justify-center px-3 py-2 text-sm border border-green-300 text-green-700 rounded-lg hover:bg-green-50 transition-colors"
                  >
                    <Play className="w-4 h-4 mr-1" />
                    Activate
                  </button>
                ) : null}
                <button
                  onClick={() => handleDeleteCampaign(campaign.id)}
                  className="px-3 py-2 text-sm border border-red-300 text-red-700 rounded-lg hover:bg-red-50 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
