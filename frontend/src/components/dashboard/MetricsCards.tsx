// frontend/src/components/dashboard/MetricsCards.tsx
"use client";

import { TrendingUp, TrendingDown, DollarSign, Activity, Shield, Brain } from "lucide-react";

const metrics = [
  {
    name: "Total Bids Today",
    value: "24,857",
    change: "+12.5%",
    trend: "up",
    icon: Activity,
    color: "blue",
  },
  {
    name: "Revenue",
    value: "$47,892",
    change: "+8.3%",
    trend: "up",
    icon: DollarSign,
    color: "green",
  },
  {
    name: "Win Rate",
    value: "68.4%",
    change: "-2.1%",
    trend: "down",
    icon: TrendingUp,
    color: "purple",
  },
  {
    name: "Fraud Blocked",
    value: "1,247",
    change: "+15.8%",
    trend: "up",
    icon: Shield,
    color: "red",
  },
  {
    name: "ML Accuracy",
    value: "94.7%",
    change: "+1.2%",
    trend: "up",
    icon: Brain,
    color: "indigo",
  },
];

const colorClasses = {
  blue: "from-blue-500 to-blue-600",
  green: "from-green-500 to-green-600",
  purple: "from-purple-500 to-purple-600",
  red: "from-red-500 to-red-600",
  indigo: "from-indigo-500 to-indigo-600",
};

export function MetricsCards() {
  return (
    <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4'>
      {metrics.map((metric) => {
        const Icon = metric.icon;
        const TrendIcon = metric.trend === "up" ? TrendingUp : TrendingDown;

        return (
          <div
            key={metric.name}
            className='bg-white rounded-xl shadow-sm border border-slate-200 p-6 hover:shadow-md transition-shadow'
          >
            <div className='flex items-center justify-between mb-4'>
              <div
                className={`w-12 h-12 rounded-lg bg-gradient-to-r ${
                  colorClasses[metric.color]
                } flex items-center justify-center`}
              >
                <Icon className='w-6 h-6 text-white' />
              </div>
              <div
                className={`flex items-center space-x-1 text-sm font-medium ${
                  metric.trend === "up" ? "text-green-600" : "text-red-600"
                }`}
              >
                <TrendIcon className='w-4 h-4' />
                <span>{metric.change}</span>
              </div>
            </div>

            <div>
              <p className='text-2xl font-bold text-slate-900 mb-1'>{metric.value}</p>
              <p className='text-sm text-slate-600'>{metric.name}</p>
            </div>
          </div>
        );
      })}
    </div>
  );
}
