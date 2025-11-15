// frontend/src/app/dashboard/alerts/components/AlertsTable.tsx
'use client';

import { useState } from "react";
import {
  Shield,
  DollarSign,
  TrendingDown,
  Brain,
  AlertTriangle,
  Target,
  CheckCircle,
  Eye,
  XCircle,
  MessageSquare,
} from "lucide-react";
import {
  Alert,
  AlertType,
  AlertStatus,
  getSeverityBadgeColor,
  getStatusBadgeColor,
  updateAlertStatus,
} from "@/lib/api/alerts";

interface AlertsTableProps {
  alerts: Alert[];
  loading: boolean;
  onAlertUpdate: () => void;
}

export function AlertsTable({
  alerts,
  loading,
  onAlertUpdate,
}: AlertsTableProps) {
  const [selectedAlerts, setSelectedAlerts] = useState<Set<string>>(new Set());
  const [updatingAlerts, setUpdatingAlerts] = useState<Set<string>>(new Set());
  const [expandedAlert, setExpandedAlert] = useState<string | null>(null);
  const [noteText, setNoteText] = useState("");

  const getTypeIcon = (type: AlertType) => {
    switch (type) {
      case "fraud":
        return Shield;
      case "budget":
        return DollarSign;
      case "performance":
        return TrendingDown;
      case "model":
        return Brain;
      case "system":
        return AlertTriangle;
      case "campaign":
        return Target;
      default:
        return AlertTriangle;
    }
  };

  const handleSelectAlert = (alertId: string) => {
    const newSelected = new Set(selectedAlerts);
    if (newSelected.has(alertId)) {
      newSelected.delete(alertId);
    } else {
      newSelected.add(alertId);
    }
    setSelectedAlerts(newSelected);
  };

  const handleSelectAll = () => {
    if (selectedAlerts.size === alerts.length) {
      setSelectedAlerts(new Set());
    } else {
      setSelectedAlerts(new Set(alerts.map((a) => a.id)));
    }
  };

  const handleStatusUpdate = async (
    alertId: string,
    status: AlertStatus,
    notes?: string
  ) => {
    setUpdatingAlerts((prev) => new Set(prev).add(alertId));
    try {
      await updateAlertStatus(alertId, status, notes);
      onAlertUpdate();
      if (expandedAlert === alertId) {
        setExpandedAlert(null);
        setNoteText("");
      }
    } catch (error) {
      console.error("Failed to update alert:", error);
      alert("Failed to update alert status");
    } finally {
      setUpdatingAlerts((prev) => {
        const newSet = new Set(prev);
        newSet.delete(alertId);
        return newSet;
      });
    }
  };

  const toggleExpand = (alertId: string) => {
    setExpandedAlert(expandedAlert === alertId ? null : alertId);
    if (expandedAlert !== alertId) {
      const alert = alerts.find((a) => a.id === alertId);
      setNoteText(alert?.notes || "");
    }
  };

  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm">
        <div className="p-6 border-b border-slate-200">
          <div className="h-6 bg-slate-200 rounded w-1/4 animate-pulse"></div>
        </div>
        <div className="p-6 space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-16 bg-slate-100 rounded animate-pulse"></div>
          ))}
        </div>
      </div>
    );
  }

  if (alerts.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm p-12 text-center">
        <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
        <h3 className="text-xl font-semibold text-slate-900 mb-2">
          No Alerts
        </h3>
        <p className="text-slate-600">
          You&apos;re all caught up! No alerts to display at the moment.
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm">
      {/* Header */}
      <div className="p-6 border-b border-slate-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <input
              type="checkbox"
              checked={selectedAlerts.size === alerts.length && alerts.length > 0}
              onChange={handleSelectAll}
              className="rounded border-slate-300 text-blue-600 focus:ring-blue-500"
            />
            <h2 className="text-lg font-semibold text-slate-900">
              All Alerts ({alerts.length})
            </h2>
          </div>
          {selectedAlerts.size > 0 && (
            <div className="flex items-center space-x-2">
              <span className="text-sm text-slate-600">
                {selectedAlerts.size} selected
              </span>
            </div>
          )}
        </div>
      </div>

      {/* Alert List */}
      <div className="divide-y divide-slate-200">
        {alerts.map((alert) => {
          const Icon = getTypeIcon(alert.type);
          const isExpanded = expandedAlert === alert.id;
          const isUpdating = updatingAlerts.has(alert.id);

          return (
            <div
              key={alert.id}
              className={`${
                alert.status === "unread" ? "bg-blue-50" : "bg-white"
              } hover:bg-slate-50 transition-colors`}
            >
              <div className="p-6">
                <div className="flex items-start space-x-4">
                  {/* Checkbox */}
                  <input
                    type="checkbox"
                    checked={selectedAlerts.has(alert.id)}
                    onChange={() => handleSelectAlert(alert.id)}
                    className="mt-1 rounded border-slate-300 text-blue-600 focus:ring-blue-500"
                  />

                  {/* Icon */}
                  <div className="flex-shrink-0">
                    <div
                      className={`p-3 rounded-lg ${getSeverityBadgeColor(
                        alert.severity
                      )}`}
                    >
                      <Icon className="w-5 h-5" />
                    </div>
                  </div>

                  {/* Content */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between mb-2">
                      <div className="flex-1">
                        <div className="flex items-center space-x-2 mb-1">
                          <h3 className="text-base font-semibold text-slate-900">
                            {alert.title}
                          </h3>
                          <span
                            className={`px-2 py-1 text-xs font-medium rounded-full border ${getSeverityBadgeColor(
                              alert.severity
                            )}`}
                          >
                            {alert.severity.toUpperCase()}
                          </span>
                        </div>
                        <p className="text-sm text-slate-600 mb-2">
                          {alert.message}
                        </p>
                        {alert.campaign_name && (
                          <div className="text-xs text-slate-500">
                            Campaign: {alert.campaign_name}
                          </div>
                        )}
                      </div>
                      <span
                        className={`ml-4 px-3 py-1 text-xs font-medium rounded-full border ${getStatusBadgeColor(
                          alert.status
                        )}`}
                      >
                        {alert.status.replace("_", " ").toUpperCase()}
                      </span>
                    </div>

                    {/* Metadata */}
                    <div className="flex items-center space-x-4 text-xs text-slate-500 mb-3">
                      <span>
                        {new Date(alert.created_at).toLocaleString()}
                      </span>
                      <span>•</span>
                      <span className="capitalize">{alert.type} Alert</span>
                    </div>

                    {/* Actions */}
                    <div className="flex items-center space-x-2">
                      {alert.status === "unread" && (
                        <button
                          onClick={() => handleStatusUpdate(alert.id, "read")}
                          disabled={isUpdating}
                          className="flex items-center space-x-1 px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors disabled:opacity-50"
                        >
                          <Eye className="w-3 h-3" />
                          <span>Mark as Read</span>
                        </button>
                      )}
                      {(alert.status === "unread" || alert.status === "read") && (
                        <button
                          onClick={() =>
                            handleStatusUpdate(alert.id, "acknowledged")
                          }
                          disabled={isUpdating}
                          className="flex items-center space-x-1 px-3 py-1.5 text-xs font-medium text-green-600 bg-green-50 rounded-lg hover:bg-green-100 transition-colors disabled:opacity-50"
                        >
                          <CheckCircle className="w-3 h-3" />
                          <span>Acknowledge</span>
                        </button>
                      )}
                      {alert.status !== "resolved" && (
                        <button
                          onClick={() => toggleExpand(alert.id)}
                          className="flex items-center space-x-1 px-3 py-1.5 text-xs font-medium text-slate-600 bg-slate-100 rounded-lg hover:bg-slate-200 transition-colors"
                        >
                          <MessageSquare className="w-3 h-3" />
                          <span>Add Note & Resolve</span>
                        </button>
                      )}
                      {alert.status !== "dismissed" && (
                        <button
                          onClick={() =>
                            handleStatusUpdate(alert.id, "dismissed")
                          }
                          disabled={isUpdating}
                          className="flex items-center space-x-1 px-3 py-1.5 text-xs font-medium text-gray-600 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors disabled:opacity-50"
                        >
                          <XCircle className="w-3 h-3" />
                          <span>Dismiss</span>
                        </button>
                      )}
                    </div>
                  </div>
                </div>

                {/* Expanded Note Section */}
                {isExpanded && (
                  <div className="mt-4 ml-14 pl-4 border-l-2 border-slate-200">
                    <label className="block text-sm font-medium text-slate-700 mb-2">
                      Add Resolution Notes
                    </label>
                    <textarea
                      value={noteText}
                      onChange={(e) => setNoteText(e.target.value)}
                      placeholder="Add notes about how this alert was resolved..."
                      rows={3}
                      className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                    />
                    <div className="flex items-center space-x-2 mt-3">
                      <button
                        onClick={() =>
                          handleStatusUpdate(alert.id, "resolved", noteText)
                        }
                        disabled={isUpdating}
                        className="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50"
                      >
                        Resolve Alert
                      </button>
                      <button
                        onClick={() => {
                          setExpandedAlert(null);
                          setNoteText("");
                        }}
                        className="px-4 py-2 bg-slate-200 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-300 transition-colors"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
