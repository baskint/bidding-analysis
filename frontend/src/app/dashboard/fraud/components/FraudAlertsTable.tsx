// frontend/src/app/dashboard/fraud/components/FraudAlertsTable.tsx
'use client';

import { useState } from "react";
import { 
  AlertTriangle, 
  Eye, 
  CheckCircle, 
  XCircle,
  Clock,
  Filter
} from "lucide-react";
import {
  FraudAlert,
  getSeverityBadgeColor,
  getStatusBadgeColor,
  formatAlertType,
  updateFraudAlert,
} from "@/lib/api/fraud";

interface FraudAlertsTableProps {
  alerts: FraudAlert[];
  loading: boolean;
  onAlertUpdate?: () => void;
}

export function FraudAlertsTable({ alerts, loading, onAlertUpdate }: FraudAlertsTableProps) {
  const [selectedAlert, setSelectedAlert] = useState<FraudAlert | null>(null);
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [severityFilter, setSeverityFilter] = useState<string>("all");
  const [updating, setUpdating] = useState(false);

  const handleStatusUpdate = async (alertId: string, newStatus: string) => {
    setUpdating(true);
    try {
      await updateFraudAlert(alertId, newStatus);
      if (onAlertUpdate) {
        onAlertUpdate();
      }
      setSelectedAlert(null);
    } catch (error) {
      console.error("Failed to update alert:", error);
      alert("Failed to update alert status");
    } finally {
      setUpdating(false);
    }
  };

  // Filter alerts
  const filteredAlerts = alerts.filter((alert) => {
    if (statusFilter !== "all" && alert.status !== statusFilter) return false;
    if (severityFilter === "high" && alert.severity < 7) return false;
    if (severityFilter === "medium" && (alert.severity < 4 || alert.severity >= 7)) return false;
    if (severityFilter === "low" && alert.severity >= 4) return false;
    return true;
  });

  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/4"></div>
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-20 bg-slate-200 rounded"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200">
      {/* Header with Filters */}
      <div className="p-6 border-b border-slate-200">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <AlertTriangle className="w-5 h-5 text-red-500" />
            <h3 className="text-lg font-semibold text-slate-900">Fraud Alerts</h3>
          </div>
          <div className="flex items-center space-x-2 text-sm text-slate-600">
            <Filter className="w-4 h-4" />
            <span>Filters:</span>
          </div>
        </div>

        {/* Filter Controls */}
        <div className="flex flex-wrap gap-3">
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-3 py-2 text-sm border border-slate-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="investigating">Investigating</option>
            <option value="resolved">Resolved</option>
            <option value="false_positive">False Positive</option>
          </select>

          <select
            value={severityFilter}
            onChange={(e) => setSeverityFilter(e.target.value)}
            className="px-3 py-2 text-sm border border-slate-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="all">All Severity</option>
            <option value="high">High (7-10)</option>
            <option value="medium">Medium (4-6)</option>
            <option value="low">Low (1-3)</option>
          </select>
        </div>
      </div>

      {/* Alerts List */}
      <div className="divide-y divide-slate-200">
        {filteredAlerts.length === 0 ? (
          <div className="p-8 text-center text-slate-600">
            <Shield className="w-12 h-12 mx-auto mb-3 text-slate-400" />
            <p className="text-lg font-medium">No fraud alerts found</p>
            <p className="text-sm">
              {statusFilter !== "all" || severityFilter !== "all"
                ? "Try adjusting your filters"
                : "Your campaigns are currently secure"}
            </p>
          </div>
        ) : (
          filteredAlerts.map((alert) => (
            <div
              key={alert.id}
              className={`p-5 hover:bg-slate-50 transition-colors ${
                selectedAlert?.id === alert.id ? "bg-blue-50" : ""
              }`}
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1">
                  <div className="flex items-center space-x-2 mb-2">
                    <AlertTriangle
                      className={`w-4 h-4 ${
                        alert.severity >= 7
                          ? "text-red-500"
                          : alert.severity >= 4
                          ? "text-yellow-500"
                          : "text-blue-500"
                      }`}
                    />
                    <h4 className="font-medium text-slate-900">
                      {formatAlertType(alert.alert_type)}
                    </h4>
                    <span
                      className={`px-2 py-0.5 rounded-full text-xs font-medium ${getSeverityBadgeColor(
                        alert.severity
                      )}`}
                    >
                      Severity: {alert.severity}
                    </span>
                    <span
                      className={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusBadgeColor(
                        alert.status
                      )}`}
                    >
                      {alert.status.replace("_", " ").toUpperCase()}
                    </span>
                  </div>
                  <p className="text-sm text-slate-600 mb-2">{alert.description}</p>
                  <div className="flex items-center space-x-4 text-xs text-slate-500">
                    <span className="flex items-center">
                      <Clock className="w-3 h-3 mr-1" />
                      {new Date(alert.detected_at).toLocaleString()}
                    </span>
                    {alert.affected_user_ids && alert.affected_user_ids.length > 0 && (
                      <span>
                        {alert.affected_user_ids.length} affected{" "}
                        {alert.affected_user_ids.length === 1 ? "user" : "users"}
                      </span>
                    )}
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex items-center space-x-2 ml-4">
                  <button
                    onClick={() =>
                      setSelectedAlert(selectedAlert?.id === alert.id ? null : alert)
                    }
                    className="p-2 hover:bg-white rounded-lg transition-colors"
                    title="View details"
                  >
                    <Eye className="w-4 h-4 text-slate-400" />
                  </button>

                  {alert.status === "active" && (
                    <>
                      <button
                        onClick={() => handleStatusUpdate(alert.id, "investigating")}
                        disabled={updating}
                        className="px-3 py-1 text-xs bg-yellow-600 text-white rounded hover:bg-yellow-700 transition-colors disabled:opacity-50"
                        title="Mark as investigating"
                      >
                        Investigate
                      </button>
                      <button
                        onClick={() => handleStatusUpdate(alert.id, "resolved")}
                        disabled={updating}
                        className="p-2 hover:bg-green-50 rounded-lg transition-colors disabled:opacity-50"
                        title="Mark as resolved"
                      >
                        <CheckCircle className="w-4 h-4 text-green-600" />
                      </button>
                      <button
                        onClick={() => handleStatusUpdate(alert.id, "false_positive")}
                        disabled={updating}
                        className="p-2 hover:bg-gray-50 rounded-lg transition-colors disabled:opacity-50"
                        title="Mark as false positive"
                      >
                        <XCircle className="w-4 h-4 text-gray-600" />
                      </button>
                    </>
                  )}
                </div>
              </div>

              {/* Expanded Details */}
              {selectedAlert?.id === alert.id && (
                <div className="mt-4 p-4 bg-white rounded-lg border border-slate-200">
                  <h5 className="font-medium text-slate-900 mb-2">Alert Details</h5>
                  <dl className="grid grid-cols-2 gap-3 text-sm">
                    <div>
                      <dt className="font-medium text-slate-700">Alert ID:</dt>
                      <dd className="text-slate-600 font-mono text-xs">{alert.id}</dd>
                    </div>
                    <div>
                      <dt className="font-medium text-slate-700">Campaign ID:</dt>
                      <dd className="text-slate-600 font-mono text-xs">
                        {alert.campaign_id}
                      </dd>
                    </div>
                    <div>
                      <dt className="font-medium text-slate-700">Detected:</dt>
                      <dd className="text-slate-600">
                        {new Date(alert.detected_at).toLocaleString()}
                      </dd>
                    </div>
                    {alert.resolved_at && (
                      <div>
                        <dt className="font-medium text-slate-700">Resolved:</dt>
                        <dd className="text-slate-600">
                          {new Date(alert.resolved_at).toLocaleString()}
                        </dd>
                      </div>
                    )}
                  </dl>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}

function Shield(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
    </svg>
  );
}