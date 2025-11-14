// frontend/src/app/dashboard/fraud/page.tsx
'use client';

import { useState, useEffect } from "react";
import { Shield, RefreshCw, Calendar } from "lucide-react";
import {
  getFraudOverview,
  getRealFraudAlerts,
  getFraudTrends,
  getDeviceFraudAnalysis,
  getGeoFraudAnalysis,
  type FraudOverview as FraudOverviewType,
  type FraudAlert,
  type FraudTrend,
  type DeviceFraud,
  type GeoFraud,
} from "@/lib/api/fraud";
import { FraudOverview } from "./components/FraudOverview";
import { FraudAlertsTable } from "./components/FraudAlertsTable";
import { FraudTrendsChart } from "./components/FraudTrendsChart";
import { DeviceGeoAnalysis } from "./components/DeviceGeoAnalysis";

export default function FraudPage() {
  const [overview, setOverview] = useState<FraudOverviewType | null>(null);
  const [alerts, setAlerts] = useState<FraudAlert[]>([]);
  const [trends, setTrends] = useState<FraudTrend[]>([]);
  const [deviceData, setDeviceData] = useState<DeviceFraud[]>([]);
  const [geoData, setGeoData] = useState<GeoFraud[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [dateRange, setDateRange] = useState(30);

  const fetchData = async (showRefreshing = false) => {
    if (showRefreshing) setRefreshing(true);
    try {
      const [overviewData, alertsData, trendsData, deviceDataResp, geoDataResp] =
        await Promise.all([
          getFraudOverview(dateRange),
          getRealFraudAlerts({ limit: 50 }),
          getFraudTrends(dateRange),
          getDeviceFraudAnalysis(dateRange),
          getGeoFraudAnalysis(dateRange),
        ]);

      setOverview(overviewData);
      setAlerts(alertsData);
      setTrends(trendsData);
      setDeviceData(deviceDataResp);
      setGeoData(geoDataResp);
    } catch (error) {
      console.error("Failed to fetch fraud data:", error);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [dateRange]);

  const handleRefresh = () => {
    fetchData(true);
  };

  return (
    <div className="min-h-screen bg-slate-50 p-6">
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center space-x-3">
            <div className="p-3 bg-red-100 rounded-xl">
              <Shield className="w-8 h-8 text-red-600" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-slate-900">Fraud Detection</h1>
              <p className="text-slate-600">
                Monitor and prevent fraudulent bidding activities
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
              <RefreshCw className={`w-4 h-4 ${refreshing ? "animate-spin" : ""}`} />
              <span>Refresh</span>
            </button>
          </div>
        </div>
      </div>

      {/* Overview Section */}
      <FraudOverview data={overview} loading={loading} />

      {/* Charts and Analysis */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <FraudTrendsChart trends={trends} loading={loading} />
        <DeviceGeoAnalysis
          deviceData={deviceData}
          geoData={geoData}
          loading={loading}
        />
      </div>

      {/* Alerts Table */}
      <FraudAlertsTable
        alerts={alerts}
        loading={loading}
        onAlertUpdate={() => fetchData(true)}
      />

      {/* Footer Info */}
      <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
        <div className="flex items-start space-x-3">
          <Shield className="w-5 h-5 text-blue-600 mt-0.5" />
          <div className="flex-1">
            <h4 className="font-medium text-blue-900 mb-1">About Fraud Detection</h4>
            <p className="text-sm text-blue-800">
              Our fraud detection system monitors bidding patterns in real-time to identify
              and prevent fraudulent activities. Suspicious patterns are automatically
              flagged for review. You can mark alerts as false positives or take action to
              block malicious actors.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
