// frontend/src/app/dashboard/fraud/components/FraudOverview.tsx
'use client';

import { Shield, AlertTriangle, DollarSign, TrendingDown } from "lucide-react";
import { FraudOverview as FraudOverviewType } from "@/lib/api/fraud";
import { getThreatLevelColor, getThreatLevelBadgeColor } from "@/lib/api/fraud";

interface FraudOverviewProps {
  data: FraudOverviewType | null;
  loading: boolean;
}

export function FraudOverview({ data, loading }: FraudOverviewProps) {
  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 animate-pulse">
            <div className="h-4 bg-slate-200 rounded w-1/2 mb-2"></div>
            <div className="h-8 bg-slate-200 rounded w-3/4"></div>
          </div>
        ))}
      </div>
    );
  }

  if (!data) {
    return null;
  }

  const stats = [
    {
      label: "Total Alerts",
      value: data.total_alerts,
      icon: AlertTriangle,
      color: "text-orange-600",
      bgColor: "bg-orange-50",
    },
    {
      label: "Active Alerts",
      value: data.active_alerts,
      icon: Shield,
      color: data.active_alerts > 5 ? "text-red-600" : "text-green-600",
      bgColor: data.active_alerts > 5 ? "bg-red-50" : "bg-green-50",
    },
    {
      label: "Blocked Bids",
      value: data.blocked_bids.toLocaleString(),
      icon: TrendingDown,
      color: "text-blue-600",
      bgColor: "bg-blue-50",
    },
    {
      label: "Amount Saved",
      value: `$${data.amount_saved.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      })}`,
      icon: DollarSign,
      color: "text-emerald-600",
      bgColor: "bg-emerald-50",
    },
  ];

  return (
    <div className="space-y-6 mb-6">
      {/* Threat Level Banner */}
      <div
        className={`rounded-xl p-4 border ${
          data.threat_level === "critical"
            ? "bg-red-50 border-red-200"
            : data.threat_level === "high"
            ? "bg-orange-50 border-orange-200"
            : data.threat_level === "medium"
            ? "bg-yellow-50 border-yellow-200"
            : "bg-green-50 border-green-200"
        }`}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Shield
              className={`w-6 h-6 ${getThreatLevelColor(data.threat_level)}`}
            />
            <div>
              <div className="font-semibold text-slate-900">
                Current Threat Level
              </div>
              <span
                className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${getThreatLevelBadgeColor(
                  data.threat_level
                )}`}
              >
                {data.threat_level.toUpperCase()}
              </span>
            </div>
          </div>
          {data.active_alerts > 0 && (
            <div className="text-right">
              <div className="text-2xl font-bold text-slate-900">
                {data.active_alerts}
              </div>
              <div className="text-sm text-slate-600">
                Active {data.active_alerts === 1 ? "Alert" : "Alerts"}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, index) => (
          <div
            key={index}
            className="bg-white rounded-xl shadow-sm border border-slate-200 p-6"
          >
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium text-slate-600">
                {stat.label}
              </span>
              <div className={`p-2 rounded-lg ${stat.bgColor}`}>
                <stat.icon className={`w-5 h-5 ${stat.color}`} />
              </div>
            </div>
            <div className={`text-2xl font-bold ${stat.color}`}>
              {stat.value}
            </div>
          </div>
        ))}
      </div>

      {/* Fraud Types Breakdown */}
      {Object.keys(data.alerts_by_type).length > 0 && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
          <h3 className="text-lg font-semibold text-slate-900 mb-4">
            Alerts by Type
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
            {Object.entries(data.alerts_by_type).map(([type, count]) => (
              <div key={type} className="text-center">
                <div className="text-2xl font-bold text-slate-900">{count}</div>
                <div className="text-sm text-slate-600 capitalize">
                  {type.replace(/_/g, " ")}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Top Affected Campaigns */}
      {data.top_affected_campaigns.length > 0 && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
          <h3 className="text-lg font-semibold text-slate-900 mb-4">
            Most Affected Campaigns
          </h3>
          <div className="space-y-3">
            {data.top_affected_campaigns.map((campaign) => (
              <div
                key={campaign.campaign_id}
                className="flex items-center justify-between p-3 rounded-lg bg-slate-50 hover:bg-slate-100 transition-colors"
              >
                <div className="flex-1">
                  <div className="font-medium text-slate-900">
                    {campaign.campaign_name}
                  </div>
                  <div className="text-sm text-slate-600">
                    {campaign.fraud_attempts} fraud{" "}
                    {campaign.fraud_attempts === 1 ? "attempt" : "attempts"}
                  </div>
                </div>
                <div className="flex items-center space-x-3">
                  <div className="text-right">
                    <div className="text-sm font-medium text-slate-900">
                      Risk Score: {campaign.risk_score.toFixed(1)}
                    </div>
                    <span
                      className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${getThreatLevelBadgeColor(
                        campaign.threat_level
                      )}`}
                    >
                      {campaign.threat_level.toUpperCase()}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
