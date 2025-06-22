// frontend/src/components/dashboard/FraudAlerts.tsx
"use client";

import { AlertTriangle, Shield, Eye, X } from "lucide-react";

const alerts: Array<{
  id: number;
  severity: "high" | "medium" | "low";
  title: string;
  description: string;
  source: string;
  blocked: number;
  time: string;
  status: "active" | "investigating" | "resolved";
}> = [
  {
    id: 1,
    severity: "high",
    title: "Click Farm Detected",
    description: "Unusual traffic pattern from AS64512",
    source: "192.168.1.0/24",
    blocked: 247,
    time: "3 min ago",
    status: "active",
  },
  {
    id: 2,
    severity: "medium",
    title: "Suspicious User Agent",
    description: "Bot-like behavior detected",
    source: "User-Agent: HeadlessChrome",
    blocked: 12,
    time: "8 min ago",
    status: "investigating",
  },
  {
    id: 3,
    severity: "low",
    title: "Geolocation Mismatch",
    description: "IP location differs from device GPS",
    source: "Campaign: Electronics",
    blocked: 5,
    time: "15 min ago",
    status: "resolved",
  },
  {
    id: 4,
    severity: "high",
    title: "Rapid Click Pattern",
    description: "Multiple clicks within 100ms",
    source: "10.0.0.45",
    blocked: 89,
    time: "22 min ago",
    status: "active",
  },
];

const severityColors: Record<"high" | "medium" | "low", string> = {
  high: "border-red-200 bg-red-50",
  medium: "border-yellow-200 bg-yellow-50",
  low: "border-blue-200 bg-blue-50",
};

const severityBadges: Record<"high" | "medium" | "low", string> = {
  high: "bg-red-100 text-red-800",
  medium: "bg-yellow-100 text-yellow-800",
  low: "bg-blue-100 text-blue-800",
};

const statusColors: Record<"active" | "investigating" | "resolved", string> = {
  active: "text-red-600",
  investigating: "text-yellow-600",
  resolved: "text-green-600",
};

export function FraudAlerts() {
  return (
    <div className='bg-white rounded-xl shadow-sm border border-slate-200 p-6'>
      <div className='flex items-center justify-between mb-6'>
        <div className='flex items-center space-x-2'>
          <Shield className='w-5 h-5 text-red-500' />
          <div>
            <h3 className='text-lg font-semibold text-slate-900'>Fraud Alerts</h3>
            <p className='text-sm text-slate-600'>Security incidents and threats</p>
          </div>
        </div>
        <div className='flex items-center space-x-2'>
          <div className='w-2 h-2 bg-red-500 rounded-full animate-pulse'></div>
          <span className='text-xs text-red-600 font-medium'>2 Active</span>
        </div>
      </div>

      <div className='space-y-3'>
        {alerts.map((alert) => (
          <div key={alert.id} className={`border rounded-lg p-4 ${severityColors[alert.severity]}`}>
            <div className='flex items-start justify-between mb-2'>
              <div className='flex items-start space-x-3'>
                <AlertTriangle
                  className={`w-4 h-4 mt-0.5 ${
                    alert.severity === "high"
                      ? "text-red-500"
                      : alert.severity === "medium"
                      ? "text-yellow-500"
                      : "text-blue-500"
                  }`}
                />
                <div className='flex-1'>
                  <div className='flex items-center space-x-2 mb-1'>
                    <h4 className='text-sm font-medium text-slate-900'>{alert.title}</h4>
                    <span
                      className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                        severityBadges[alert.severity]
                      }`}
                    >
                      {alert.severity.toUpperCase()}
                    </span>
                  </div>
                  <p className='text-xs text-slate-600 mb-2'>{alert.description}</p>
                  <div className='flex items-center justify-between'>
                    <div className='text-xs text-slate-500'>
                      <span className='font-medium'>Source:</span> {alert.source}
                    </div>
                    <div className='text-xs text-slate-500'>{alert.time}</div>
                  </div>
                </div>
              </div>

              <div className='flex items-center space-x-1'>
                <button className='p-1 hover:bg-white rounded transition-colors'>
                  <Eye className='w-3 h-3 text-slate-400' />
                </button>
                <button className='p-1 hover:bg-white rounded transition-colors'>
                  <X className='w-3 h-3 text-slate-400' />
                </button>
              </div>
            </div>

            <div className='flex items-center justify-between'>
              <div className='flex items-center space-x-4'>
                <div className='text-xs'>
                  <span className='text-slate-500'>Blocked:</span>
                  <span className='font-medium text-slate-900 ml-1'>{alert.blocked}</span>
                </div>
                <div className={`text-xs font-medium ${statusColors[alert.status]}`}>
                  {alert.status.charAt(0).toUpperCase() + alert.status.slice(1)}
                </div>
              </div>

              {alert.status === "active" && (
                <button className='text-xs bg-red-600 text-white px-2 py-1 rounded hover:bg-red-700 transition-colors'>
                  Block All
                </button>
              )}
            </div>
          </div>
        ))}
      </div>

      <div className='mt-4 pt-4 border-t border-slate-200'>
        <button className='w-full text-sm text-slate-600 hover:text-slate-900 font-medium'>
          View All Security Events →
        </button>
      </div>
    </div>
  );
}
