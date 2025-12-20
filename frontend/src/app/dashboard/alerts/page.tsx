// frontend/src/app/dashboard/alerts/page.tsx
'use client';

import { useState, useEffect, useCallback } from "react";
import { Bell, RefreshCw, Calendar, Settings } from "lucide-react";
import {
  getAlerts,
  getAlertOverview,
  Alert,
  AlertOverview as AlertOverviewType,
  AlertType,
  AlertSeverity,
  AlertStatus,
} from "@/lib/api/alerts";
import { AlertOverview } from "./components/AlertOverview";
import { AlertsTable } from "./components/AlertsTable";
import { AlertFilters } from "./components/AlertFilters";

export default function AlertsPage() {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [overview, setOverview] = useState<AlertOverviewType | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [dateRange, setDateRange] = useState(30);

  // Filters
  const [selectedType, setSelectedType] = useState<AlertType | "all">("all");
  const [selectedSeverity, setSelectedSeverity] = useState<AlertSeverity | "all">("all");
  const [selectedStatus, setSelectedStatus] = useState<AlertStatus | "all">("all");

  const fetchData = useCallback(
    async (showRefreshing = false) => {
      if (showRefreshing) setRefreshing(true);
      try {
        // Build filter params
        const params: {
          type?: AlertType;
          severity?: AlertSeverity;
          status?: AlertStatus;
          limit?: number;
        } = {
          limit: 100,
        };

        if (selectedType !== "all") params.type = selectedType;
        if (selectedSeverity !== "all") params.severity = selectedSeverity;
        if (selectedStatus !== "all") params.status = selectedStatus;

        const [alertsData, overviewData] = await Promise.all([
          getAlerts(params),
          getAlertOverview(dateRange),
        ]);

        setAlerts(alertsData);
        setOverview(overviewData);
      } catch (error) {
        console.error("Failed to fetch alerts data:", error);
      } finally {
        setLoading(false);
        setRefreshing(false);
      }
    },
    [dateRange, selectedType, selectedSeverity, selectedStatus]
  );

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleRefresh = () => {
    fetchData(true);
  };

  const handleClearFilters = () => {
    setSelectedType("all");
    setSelectedSeverity("all");
    setSelectedStatus("all");
  };

  return (
    <div className="min-h-screen bg-slate-50 p-6 dark:bg-slate-900">
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center space-x-3">
            <div className="p-3 bg-blue-100 rounded-xl">
              <Bell className="w-8 h-8 text-blue-600" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-100">
                Alerts & Notifications
              </h1>
              <p className="text-slate-600">
                Monitor and manage all platform alerts in one place
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-3">
            {/* Date Range Selector */}
            <div className="flex items-center space-x-2 bg-white px-4 py-2 rounded-lg border border-slate-200">
              <Calendar className="w-4 h-4 text-slate-600" />
              <select
                value={dateRange}
                onChange={(e) => setDateRange(Number(e.target.value))}
                className="bg-transparent border-none focus:ring-0 text-sm font-medium text-slate-700"
              >
                <option value={7}>Last 7 Days</option>
                <option value={14}>Last 14 Days</option>
                <option value={30}>Last 30 Days</option>
                <option value={60}>Last 60 Days</option>
                <option value={90}>Last 90 Days</option>
              </select>
            </div>

            {/* Refresh Button */}
            <button
              onClick={handleRefresh}
              disabled={refreshing}
              className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <RefreshCw
                className={`w-4 h-4 ${refreshing ? "animate-spin" : ""}`}
              />
              <span>Refresh</span>
            </button>

            {/* Settings Button */}
            <button
              className="flex items-center space-x-2 px-4 py-2 bg-white border border-slate-200 text-slate-700 rounded-lg hover:bg-slate-50 transition-colors"
              title="Configure alert rules (Coming soon)"
            >
              <Settings className="w-4 h-4" />
              <span>Settings</span>
            </button>
          </div>
        </div>
      </div>

      {/* Overview Cards */}
      <AlertOverview data={overview} loading={loading} />

      {/* Filters */}
      <AlertFilters
        selectedType={selectedType}
        selectedSeverity={selectedSeverity}
        selectedStatus={selectedStatus}
        onTypeChange={setSelectedType}
        onSeverityChange={setSelectedSeverity}
        onStatusChange={setSelectedStatus}
        onClearFilters={handleClearFilters}
      />

      {/* Alerts Table */}
      <AlertsTable
        alerts={alerts}
        loading={loading}
        onAlertUpdate={() => fetchData(true)}
      />

      {/* Info Footer */}
      <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
        <div className="flex items-start space-x-3">
          <Bell className="w-5 h-5 text-blue-600 mt-0.5" />
          <div className="flex-1">
            <h4 className="font-medium text-blue-900 mb-1">
              About Alerts & Notifications
            </h4>
            <p className="text-sm text-blue-800">
              This unified alerts center monitors all aspects of your bidding platform including
              fraud detection, budget thresholds, campaign performance, ML model accuracy, and
              system health. Configure custom alert rules and thresholds to stay informed about
              what matters most to your campaigns.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
