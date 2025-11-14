// frontend/src/app/dashboard/fraud/components/DeviceGeoAnalysis.tsx
'use client';

import { useState } from "react";
import { Monitor, Globe } from "lucide-react";
import { DeviceFraud, GeoFraud } from "@/lib/api/fraud";

interface DeviceGeoAnalysisProps {
  deviceData: DeviceFraud[];
  geoData: GeoFraud[];
  loading: boolean;
}

export function DeviceGeoAnalysis({ deviceData, geoData, loading }: DeviceGeoAnalysisProps) {
  const [activeTab, setActiveTab] = useState<"device" | "geo">("device");

  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/3"></div>
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-16 bg-slate-200 rounded"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200">
      {/* Tabs */}
      <div className="border-b border-slate-200">
        <div className="flex">
          <button
            onClick={() => setActiveTab("device")}
            className={`flex items-center space-x-2 px-6 py-4 font-medium transition-colors ${
              activeTab === "device"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-600 hover:text-slate-900"
            }`}
          >
            <Monitor className="w-4 h-4" />
            <span>Device Analysis</span>
          </button>
          <button
            onClick={() => setActiveTab("geo")}
            className={`flex items-center space-x-2 px-6 py-4 font-medium transition-colors ${
              activeTab === "geo"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-600 hover:text-slate-900"
            }`}
          >
            <Globe className="w-4 h-4" />
            <span>Geographic Analysis</span>
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="p-6">
        {activeTab === "device" ? (
          <DeviceTable data={deviceData} />
        ) : (
          <GeoTable data={geoData} />
        )}
      </div>
    </div>
  );
}

function DeviceTable({ data }: { data: DeviceFraud[] }) {
  if (data.length === 0) {
    return (
      <div className="text-center py-12 text-slate-600">
        <Monitor className="w-12 h-12 mx-auto mb-3 text-slate-400" />
        <p>No suspicious device activity detected</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="text-left text-sm font-medium text-slate-700 border-b border-slate-200">
            <th className="pb-3">Device Type</th>
            <th className="pb-3">Browser</th>
            <th className="pb-3">Operating System</th>
            <th className="pb-3 text-right">Total Bids</th>
            <th className="pb-3 text-right">Fraud Bids</th>
            <th className="pb-3 text-right">Fraud Rate</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {(data && data.length > 0) && data?.map((item, index) => (
            <tr key={index} className="hover:bg-slate-50 transition-colors">
              <td className="py-3">
                <span className="font-medium text-slate-900">{item.device_type || "Unknown"}</span>
              </td>
              <td className="py-3 text-slate-600">{item.browser || "N/A"}</td>
              <td className="py-3 text-slate-600">{item.os || "N/A"}</td>
              <td className="py-3 text-right text-slate-900">{item.total_bids.toLocaleString()}</td>
              <td className="py-3 text-right">
                <span className="font-medium text-red-600">{item.fraud_bids.toLocaleString()}</span>
              </td>
              <td className="py-3 text-right">
                <span
                  className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${
                    item.fraud_rate > 0.5
                      ? "bg-red-100 text-red-800"
                      : item.fraud_rate > 0.2
                      ? "bg-yellow-100 text-yellow-800"
                      : "bg-blue-100 text-blue-800"
                  }`}
                >
                  {(item.fraud_rate * 100).toFixed(1)}%
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function GeoTable({ data }: { data: GeoFraud[] }) {
  if (data.length === 0) {
    return (
      <div className="text-center py-12 text-slate-600">
        <Globe className="w-12 h-12 mx-auto mb-3 text-slate-400" />
        <p>No suspicious geographic activity detected</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="text-left text-sm font-medium text-slate-700 border-b border-slate-200">
            <th className="pb-3">Country</th>
            <th className="pb-3">Region</th>
            <th className="pb-3">City</th>
            <th className="pb-3 text-right">Total Bids</th>
            <th className="pb-3 text-right">Fraud Bids</th>
            <th className="pb-3 text-right">Fraud Rate</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {data.map((item, index) => (
            <tr key={index} className="hover:bg-slate-50 transition-colors">
              <td className="py-3">
                <span className="font-medium text-slate-900">{item.country || "Unknown"}</span>
              </td>
              <td className="py-3 text-slate-600">{item.region || "N/A"}</td>
              <td className="py-3 text-slate-600">{item.city || "N/A"}</td>
              <td className="py-3 text-right text-slate-900">{item.total_bids.toLocaleString()}</td>
              <td className="py-3 text-right">
                <span className="font-medium text-red-600">{item.fraud_bids.toLocaleString()}</span>
              </td>
              <td className="py-3 text-right">
                <span
                  className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${
                    item.fraud_rate > 0.5
                      ? "bg-red-100 text-red-800"
                      : item.fraud_rate > 0.2
                      ? "bg-yellow-100 text-yellow-800"
                      : "bg-blue-100 text-blue-800"
                  }`}
                >
                  {(item.fraud_rate * 100).toFixed(1)}%
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
