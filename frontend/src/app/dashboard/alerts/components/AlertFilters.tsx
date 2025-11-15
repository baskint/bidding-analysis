// frontend/src/app/dashboard/alerts/components/AlertFilters.tsx
"use client";

import { Filter, X } from "lucide-react";
import { AlertType, AlertSeverity, AlertStatus } from "@/lib/api/alerts";

interface AlertFiltersProps {
  selectedType: AlertType | "all";
  selectedSeverity: AlertSeverity | "all";
  selectedStatus: AlertStatus | "all";
  onTypeChange: (type: AlertType | "all") => void;
  onSeverityChange: (severity: AlertSeverity | "all") => void;
  onStatusChange: (status: AlertStatus | "all") => void;
  onClearFilters: () => void;
}

export function AlertFilters({
  selectedType,
  selectedSeverity,
  selectedStatus,
  onTypeChange,
  onSeverityChange,
  onStatusChange,
  onClearFilters,
}: AlertFiltersProps) {
  const hasActiveFilters =
    selectedType !== "all" ||
    selectedSeverity !== "all" ||
    selectedStatus !== "all";

  return (
    <div className="bg-white rounded-xl shadow-sm p-6 mb-6">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-2">
          <Filter className="w-5 h-5 text-slate-600" />
          <h3 className="text-lg font-semibold text-slate-900">Filters</h3>
        </div>
        {hasActiveFilters && (
          <button
            onClick={onClearFilters}
            className="flex items-center space-x-1 px-3 py-1.5 text-sm text-slate-600 hover:text-slate-900 transition-colors"
          >
            <X className="w-4 h-4" />
            <span>Clear All</span>
          </button>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Type Filter */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Alert Type
          </label>
          <select
            value={selectedType}
            onChange={(e) => onTypeChange(e.target.value as AlertType | "all")}
            className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
          >
            <option value="all">All Types</option>
            <option value="fraud">Fraud</option>
            <option value="budget">Budget</option>
            <option value="performance">Performance</option>
            <option value="model">ML Model</option>
            <option value="system">System</option>
            <option value="campaign">Campaign</option>
          </select>
        </div>

        {/* Severity Filter */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Severity
          </label>
          <select
            value={selectedSeverity}
            onChange={(e) =>
              onSeverityChange(e.target.value as AlertSeverity | "all")
            }
            className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
          >
            <option value="all">All Severities</option>
            <option value="critical">Critical</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </div>

        {/* Status Filter */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Status
          </label>
          <select
            value={selectedStatus}
            onChange={(e) =>
              onStatusChange(e.target.value as AlertStatus | "all")
            }
            className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
          >
            <option value="all">All Statuses</option>
            <option value="unread">Unread</option>
            <option value="read">Read</option>
            <option value="acknowledged">Acknowledged</option>
            <option value="resolved">Resolved</option>
            <option value="dismissed">Dismissed</option>
          </select>
        </div>
      </div>

      {/* Active Filters Summary */}
      {hasActiveFilters && (
        <div className="mt-4 flex flex-wrap gap-2">
          {selectedType !== "all" && (
            <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              Type: {selectedType}
              <button
                onClick={() => onTypeChange("all")}
                className="ml-1 hover:text-blue-900"
              >
                <X className="w-3 h-3" />
              </button>
            </span>
          )}
          {selectedSeverity !== "all" && (
            <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-orange-100 text-orange-800">
              Severity: {selectedSeverity}
              <button
                onClick={() => onSeverityChange("all")}
                className="ml-1 hover:text-orange-900"
              >
                <X className="w-3 h-3" />
              </button>
            </span>
          )}
          {selectedStatus !== "all" && (
            <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
              Status: {selectedStatus}
              <button
                onClick={() => onStatusChange("all")}
                className="ml-1 hover:text-green-900"
              >
                <X className="w-3 h-3" />
              </button>
            </span>
          )}
        </div>
      )}
    </div>
  );
}
