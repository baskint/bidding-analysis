// frontend/src/components/dashboard/MetricsCards.tsx
'use client';
import { useState, useEffect } from 'react';
import { TrendingUp, TrendingDown, DollarSign, Activity, Shield, Brain } from 'lucide-react';

// API helper
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const getAuthHeaders = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

interface DashboardMetrics {
  active_bids: number;
  total_spend: number;
  win_rate: number;
  fraud_alerts: number;
  model_accuracy: number;
  avg_bid: number;
  conversions: number;
  total_campaigns: number;
  last_updated: string;
}

export function MetricsCards() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchMetrics() {
      try {
        const response = await fetch(`${API_BASE_URL}/trpc/analytics.getDashboardMetrics`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({}),
        });

        if (!response.ok) {
          throw new Error('Failed to fetch metrics');
        }

        const data = await response.json();
        console.log('API Response:', data); // Debug log

        // Handle tRPC response format
        if (data.result && data.result.data) {
          console.log('Setting metrics:', data.result.data); // Debug log
          setMetrics(data.result.data);
        } else {
          throw new Error('Invalid response format');
        }
      } catch (err) {
        console.error('Error fetching metrics:', err);
        setError(err instanceof Error ? err.message : 'Failed to load metrics');
      } finally {
        setLoading(false);
      }
    }

    fetchMetrics();
  }, []);

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 animate-pulse">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 rounded-lg bg-slate-200"></div>
              <div className="w-16 h-4 bg-slate-200 rounded"></div>
            </div>
            <div className="w-20 h-8 bg-slate-200 rounded mb-2"></div>
            <div className="w-24 h-4 bg-slate-200 rounded"></div>
          </div>
        ))}
      </div>
    );
  }

  if (error || !metrics) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-800">Failed to load metrics: {error || 'No data'}</p>
      </div>
    );
  }

  console.log('Rendering with metrics:', metrics); // Debug log

  // Safe defaults
  const activeBids = metrics.active_bids ?? 0;
  const totalSpend = metrics.total_spend ?? 0;
  const winRate = metrics.win_rate ?? 0;
  const fraudAlerts = metrics.fraud_alerts ?? 0;
  const modelAccuracy = metrics.model_accuracy ?? 0;

  const metricCards = [
    {
      name: 'Active Bids',
      value: activeBids.toLocaleString(),
      change: '+12.5%',
      trend: 'up' as const,
      icon: Activity,
      color: 'blue' as const
    },
    {
      name: 'Total Spend',
      value: `$${totalSpend.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`,
      change: '+8.3%',
      trend: 'up' as const,
      icon: DollarSign,
      color: 'green' as const
    },
    {
      name: 'Win Rate',
      value: `${(winRate * 100).toFixed(1)}%`,
      change: '-2.1%',
      trend: 'down' as const,
      icon: TrendingUp,
      color: 'purple' as const
    },
    {
      name: 'Fraud Alerts',
      value: fraudAlerts.toLocaleString(),
      change: '+15.8%',
      trend: 'up' as const,
      icon: Shield,
      color: 'red' as const
    },
    {
      name: 'ML Accuracy',
      value: `${(modelAccuracy * 100).toFixed(1)}%`,
      change: '+1.2%',
      trend: 'up' as const,
      icon: Brain,
      color: 'indigo' as const
    }
  ];

  const colorClasses: Record<'blue' | 'green' | 'purple' | 'red' | 'indigo', string> = {
    blue: 'from-blue-500 to-blue-600',
    green: 'from-green-500 to-green-600',
    purple: 'from-purple-500 to-purple-600',
    red: 'from-red-500 to-red-600',
    indigo: 'from-indigo-500 to-indigo-600'
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
      {metricCards.map((metric) => {
        const Icon = metric.icon;
        const TrendIcon = metric.trend === 'up' ? TrendingUp : TrendingDown;

        return (
          <div
            key={metric.name}
            className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 hover:shadow-md transition-shadow"
          >
            <div className="flex items-center justify-between mb-4">
              <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${colorClasses[metric.color]} flex items-center justify-center`}>
                <Icon className="w-6 h-6 text-white" />
              </div>
              <div className={`flex items-center space-x-1 text-sm font-medium ${metric.trend === 'up' ? 'text-green-600' : 'text-red-600'
                }`}>
                <TrendIcon className="w-4 h-4" />
                <span>{metric.change}</span>
              </div>
            </div>

            <div>
              <p className="text-2xl font-bold text-slate-900 mb-1">{metric.value}</p>
              <p className="text-sm text-slate-600">{metric.name}</p>
            </div>
          </div>
        );
      })}
    </div>
  );
}
