import { Smartphone } from "lucide-react";
import type { DeviceBreakdown } from "@/lib/api/analytics";

interface DeviceBreakdownChartProps {
  devices: DeviceBreakdown[];
  loading: boolean;
}

export function DeviceBreakdownChart({
  devices,
  loading,
}:DeviceBreakdownChartProps) {
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
  const maxBids = Math.max(...devices.map((d) => d.totalBids), 1);

  console.log(devices);

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
                {device.deviceType}
              </span>
              <div className="flex items-center space-x-4 text-sm">
                <span className="text-slate-600">
                  {(device?.totalBids ?? 0).toLocaleString()} bids
                </span>
                <span className="text-blue-600 font-medium">
                  {formatPercent(device.winRate)} WR
                </span>
                <span className="text-green-600 font-medium">
                  {formatPercent(device.conversionRate)} CR
                </span>
              </div>
            </div>
            <div className="relative h-2 bg-slate-100 rounded-full overflow-hidden">
              <div
                className="absolute top-0 left-0 h-full bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full"
                style={{
                  width: `${(device.totalBids / maxBids) * 100}%`,
                }}
              ></div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
