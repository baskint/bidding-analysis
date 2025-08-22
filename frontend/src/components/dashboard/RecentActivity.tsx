// frontend/src/components/dashboard/RecentActivity.tsx
'use client';;

import { Activity, DollarSign, Shield, TrendingUp, Clock } from "lucide-react";
import React from "react";

const activities: Array<{
  id: number;
  type: string;
  title: string;
  description: string;
  amount: string | null;
  time: string;
  icon: React.ComponentType<{ className?: string }>;
  color: "green" | "red" | "blue" | "yellow";
}> = [
    {
      id: 1,
      type: "bid_won",
      title: "High-value bid won",
      description: 'Campaign "Black Friday Sale" won bid for $24.50',
      amount: "$24.50",
      time: "2 minutes ago",
      icon: TrendingUp,
      color: "green",
    },
    {
      id: 2,
      type: "fraud_blocked",
      title: "Fraud attempt blocked",
      description: "Suspicious traffic from IP 192.168.1.1 blocked",
      amount: null,
      time: "5 minutes ago",
      icon: Shield,
      color: "red",
    },
    {
      id: 3,
      type: "campaign_started",
      title: "New campaign launched",
      description: 'Campaign "Holiday Electronics" is now live',
      amount: null,
      time: "12 minutes ago",
      icon: Activity,
      color: "blue",
    },
    {
      id: 4,
      type: "high_spend",
      title: "Budget milestone reached",
      description: 'Campaign "Gaming Accessories" reached 80% of daily budget',
      amount: "$4,200",
      time: "18 minutes ago",
      icon: DollarSign,
      color: "yellow",
    },
    {
      id: 5,
      type: "bid_won",
      title: "Premium placement secured",
      description: 'Won premium ad slot for "Luxury Watches"',
      amount: "$47.80",
      time: "25 minutes ago",
      icon: TrendingUp,
      color: "green",
    },
    {
      id: 6,
      type: "fraud_blocked",
      title: "Click fraud prevented",
      description: "Blocked 150 fraudulent clicks from bot network",
      amount: null,
      time: "32 minutes ago",
      icon: Shield,
      color: "red",
    },
  ];

const colorClasses: Record<"green" | "red" | "blue" | "yellow", string> = {
  green: "bg-green-100 text-green-600",
  red: "bg-red-100 text-red-600",
  blue: "bg-blue-100 text-blue-600",
  yellow: "bg-yellow-100 text-yellow-600",
};

export function RecentActivity() {
  return (
    <div className='bg-white rounded-xl shadow-sm border border-slate-200 p-6'>
      <div className='flex items-center justify-between mb-6'>
        <div>
          <h3 className='text-lg font-semibold text-slate-900'>Recent Activity</h3>
          <p className='text-sm text-slate-600'>Latest bidding events and system updates</p>
        </div>
        <button className='text-sm text-blue-600 hover:text-blue-700 font-medium'>View All</button>
      </div>

      <div className='space-y-4'>
        {activities.map((activity) => {
          const Icon = activity.icon;
          return (
            <div
              key={activity.id}
              className='flex items-start space-x-4 p-4 rounded-lg hover:bg-slate-50 transition-colors'
            >
              <div
                className={`flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center ${colorClasses[activity.color]
                  }`}
              >
                <Icon className='w-5 h-5' />
              </div>

              <div className='flex-1 min-w-0'>
                <div className='flex items-center justify-between'>
                  <p className='text-sm font-medium text-slate-900'>{activity.title}</p>
                  {activity.amount && (
                    <span className='text-sm font-semibold text-slate-900'>{activity.amount}</span>
                  )}
                </div>
                <p className='text-sm text-slate-600 mt-1'>{activity.description}</p>
                <div className='flex items-center mt-2 text-xs text-slate-500'>
                  <Clock className='w-3 h-3 mr-1' />
                  {activity.time}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className='mt-6 pt-4 border-t border-slate-200'>
        <div className='flex items-center justify-between text-sm'>
          <span className='text-slate-600'>Last updated: {new Date().toLocaleTimeString()}</span>
          <button className='text-blue-600 hover:text-blue-700 font-medium'>Refresh</button>
        </div>
      </div>
    </div>
  );
}
