// frontend/src/app/dashboard/page.tsx
'use client';

import { MetricsCards } from '@/components/dashboard/MetricsCards';
import { BidChart } from '@/components/dashboard/BidChart';
import { CampaignPerformance } from '@/components/dashboard/CampaignPerformance';
import { RecentActivity } from '@/components/dashboard/RecentActivity';
import { FraudAlerts } from '@/components/dashboard/FraudAlerts';

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-100">Dashboard Overview</h1>
          <p className="text-slate-600 mt-1">Real-time bidding analytics and performance metrics</p>
        </div>
        <div className="flex space-x-3">
          <button className="px-4 py-2 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors">
            Export Data
          </button>
          <button className="px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:from-blue-700 hover:to-purple-700 transition-all">
            New Campaign
          </button>
        </div>
      </div>

      {/* Metrics Cards */}
      <MetricsCards />

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <BidChart />
        <CampaignPerformance />
      </div>

      {/* Bottom Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <RecentActivity />
        </div>
        <div>
          <FraudAlerts />
        </div>
      </div>
    </div>
  );
}
