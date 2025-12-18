// frontend/src/components/dashboard/FraudAlerts.tsx
'use client';
import { useState, useEffect } from 'react';
import { AlertTriangle, Shield, Clock, AlertCircle } from 'lucide-react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const getAuthHeaders = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

interface FraudAlert {
  id: string;
  type: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  status: string;
  campaign_id: string;
  detected_at: string;
}

export function FraudAlerts() {
  const [alerts, setAlerts] = useState<FraudAlert[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchAlerts() {
      try {
        const response = await fetch(`${API_BASE_URL}/trpc/analytics.getFraudAlerts`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({}),
        });

        if (!response.ok) {
          throw new Error('Failed to fetch fraud alerts');
        }

        const data = await response.json();

        if (data.result && data.result.data) {
          setAlerts(data.result.data || []);
        } else {
          throw new Error('Invalid response format');
        }
      } catch (err) {
        console.error('Error fetching fraud alerts:', err);
        setError(err instanceof Error ? err.message : 'Failed to load alerts');
      } finally {
        setLoading(false);
      }
    }

    fetchAlerts();
  }, []);

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'bg-red-100 text-red-800 border-red-200';
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      default:
        return 'bg-blue-100 text-blue-800 border-blue-200';
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case 'critical':
      case 'high':
        return AlertTriangle;
      case 'medium':
        return AlertCircle;
      default:
        return Shield;
    }
  };

  const formatAlertType = (type: string) => {
    return type
      .split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');
  };

  const formatTimeAgo = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);

    if (diffMins < 60) {
      return `${diffMins}m ago`;
    } else if (diffMins < 1440) {
      return `${Math.floor(diffMins / 60)}h ago`;
    } else {
      return `${Math.floor(diffMins / 1440)}d ago`;
    }
  };

  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-slate-900">Fraud Alerts</h3>
        </div>
        <div className="space-y-3 animate-pulse">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="p-3 border border-slate-200 rounded-lg">
              <div className="w-full h-4 bg-slate-200 rounded mb-2"></div>
              <div className="w-24 h-3 bg-slate-200 rounded"></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <h3 className="text-lg font-semibold text-slate-900 mb-4">Fraud Alerts</h3>
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800 text-sm">Failed to load alerts: {error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-slate-900">Fraud Alerts</h3>
        <span className="px-2 py-1 bg-red-100 text-red-800 text-xs font-semibold rounded-full">
          {alerts.length} Active
        </span>
      </div>

      {alerts.length === 0 ? (
        <div className="text-center py-8">
          <Shield className="w-12 h-12 text-green-500 mx-auto mb-3" />
          <p className="text-sm text-slate-600">No fraud alerts detected</p>
          <p className="text-xs text-slate-500 mt-1">All systems operating normally</p>
        </div>
      ) : (
        <div className="space-y-3">
          {alerts.map((alert) => {
            const SeverityIcon = getSeverityIcon(alert.severity);

            return (
              <div
                key={alert.id}
                className={`p-3 border rounded-lg ${getSeverityColor(alert.severity)} transition-all hover:shadow-md`}
              >
                <div className="flex items-start justify-between mb-2">
                  <div className="flex items-center space-x-2">
                    <SeverityIcon className="w-4 h-4" />
                    <span className="text-sm font-semibold">
                      {formatAlertType(alert.type)}
                    </span>
                  </div>
                  <span className="text-xs font-medium uppercase px-2 py-1 bg-white/50 rounded">
                    {alert.severity}
                  </span>
                </div>

                <div className="flex items-center justify-between text-xs">
                  <div className="flex items-center space-x-1 text-slate-600">
                    <Clock className="w-3 h-3" />
                    <span>{formatTimeAgo(alert.detected_at)}</span>
                  </div>
                  <span className="text-slate-600">
                    Status: <span className="font-medium">{alert.status}</span>
                  </span>
                </div>
              </div>
            );
          })}
        </div>
      )}

      <button className="w-full mt-4 px-4 py-2 text-sm text-blue-600 hover:bg-blue-50 rounded-lg transition-colors font-medium">
        View All Alerts →
      </button>
    </div>
  );
}
