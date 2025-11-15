// frontend/src/app/dashboard/alerts/components/AlertOverview.tsx
import { AlertTriangle, CheckCircle, Clock, XCircle } from "lucide-react";
import { AlertOverview as AlertOverviewType } from "@/lib/api/alerts";

interface AlertOverviewProps {
  data: AlertOverviewType | null;
  loading: boolean;
}

export function AlertOverview({ data, loading }: AlertOverviewProps) {
  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="bg-white rounded-xl p-6 shadow-sm animate-pulse">
            <div className="h-4 bg-slate-200 rounded w-1/2 mb-4"></div>
            <div className="h-8 bg-slate-200 rounded w-3/4 mb-2"></div>
            <div className="h-3 bg-slate-200 rounded w-full"></div>
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
      color: "bg-blue-100 text-blue-600",
      textColor: "text-blue-900",
    },
    {
      label: "Unread",
      value: data.unread_alerts,
      icon: Clock,
      color: "bg-yellow-100 text-yellow-600",
      textColor: "text-yellow-900",
    },
    {
      label: "Critical",
      value: data.critical_alerts,
      icon: XCircle,
      color: "bg-red-100 text-red-600",
      textColor: "text-red-900",
    },
    {
      label: "Resolved",
      value: data.total_alerts - data.unread_alerts,
      icon: CheckCircle,
      color: "bg-green-100 text-green-600",
      textColor: "text-green-900",
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      {stats.map((stat) => {
        const Icon = stat.icon;
        return (
          <div
            key={stat.label}
            className="bg-white rounded-xl p-6 shadow-sm hover:shadow-md transition-shadow"
          >
            <div className="flex items-center justify-between mb-4">
              <span className="text-sm font-medium text-slate-600">
                {stat.label}
              </span>
              <div className={`p-2 rounded-lg ${stat.color}`}>
                <Icon className="w-5 h-5" />
              </div>
            </div>
            <div className={`text-3xl font-bold ${stat.textColor} mb-1`}>
              {stat.value.toLocaleString()}
            </div>
            <div className="text-xs text-slate-500">
              Last 30 days
            </div>
          </div>
        );
      })}
    </div>
  );
}
