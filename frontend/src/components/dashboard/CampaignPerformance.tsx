// frontend/src/components/dashboard/CampaignPerformance.tsx
'use client';
import { useState, useEffect } from 'react';
import { TrendingUp, Target, DollarSign, MousePointerClick } from 'lucide-react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const getAuthHeaders = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

interface CampaignStats {
  total_bids: number;
  won_bids: number;
  win_rate: number;
  total_spend: number;
  conversions: number;
  avg_cpa: number;
}

export function CampaignPerformance() {
  const [stats, setStats] = useState<CampaignStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchStats() {
      try {
        const response = await fetch(`${API_BASE_URL}/trpc/campaign.getStats`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({}),
        });

        if (!response.ok) {
          throw new Error('Failed to fetch campaign stats');
        }

        const data = await response.json();

        if (data.result && data.result.data) {
          setStats(data.result.data);
        } else {
          throw new Error('Invalid response format');
        }
      } catch (err) {
        console.error('Error fetching campaign stats:', err);
        setError(err instanceof Error ? err.message : 'Failed to load stats');
      } finally {
        setLoading(false);
      }
    }

    fetchStats();
  }, []);

  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-semibold text-slate-900">Campaign Performance</h3>
        </div>
        <div className="space-y-4 animate-pulse">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
              <div className="w-32 h-4 bg-slate-200 rounded"></div>
              <div className="w-20 h-6 bg-slate-200 rounded"></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error || !stats) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <h3 className="text-lg font-semibold text-slate-900 mb-4">Campaign Performance</h3>
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800 text-sm">Failed to load campaign stats: {error || 'No data'}</p>
        </div>
      </div>
    );
  }

  const metrics = [
    {
      label: 'Total Bids',
      value: stats.total_bids.toLocaleString(),
      icon: MousePointerClick,
      color: 'blue'
    },
    {
      label: 'Won Bids',
      value: stats.won_bids.toLocaleString(),
      icon: Target,
      color: 'green'
    },
    {
      label: 'Win Rate',
      value: `${(stats.win_rate * 100).toFixed(1)}%`,
      icon: TrendingUp,
      color: 'purple'
    },
    {
      label: 'Total Spend',
      value: `$${stats.total_spend.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`,
      icon: DollarSign,
      color: 'emerald'
    },
    {
      label: 'Conversions',
      value: stats.conversions.toLocaleString(),
      icon: Target,
      color: 'indigo'
    },
    {
      label: 'Avg CPA',
      value: `$${stats.avg_cpa.toFixed(2)}`,
      icon: DollarSign,
      color: 'orange'
    }
  ];

  const colorClasses: Record<string, { bg: string; text: string }> = {
    blue: { bg: 'bg-blue-50', text: 'text-blue-600' },
    green: { bg: 'bg-green-50', text: 'text-green-600' },
    purple: { bg: 'bg-purple-50', text: 'text-purple-600' },
    emerald: { bg: 'bg-emerald-50', text: 'text-emerald-600' },
    indigo: { bg: 'bg-indigo-50', text: 'text-indigo-600' },
    orange: { bg: 'bg-orange-50', text: 'text-orange-600' }
  };

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-slate-900">Campaign Performance</h3>
        <button className="text-sm text-blue-600 hover:text-blue-700 font-medium">
          View Details →
        </button>
      </div>

      <div className="space-y-3">
        {metrics.map((metric) => {
          const Icon = metric.icon;
          const colors = colorClasses[metric.color];

          return (
            <div
              key={metric.label}
              className="flex items-center justify-between p-4 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors"
            >
              <div className="flex items-center space-x-3">
                <div className={`w-10 h-10 rounded-lg ${colors.bg} flex items-center justify-center`}>
                  <Icon className={`w-5 h-5 ${colors.text}`} />
                </div>
                <span className="text-sm font-medium text-slate-700">{metric.label}</span>
              </div>
              <span className="text-lg font-bold text-slate-900">{metric.value}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
