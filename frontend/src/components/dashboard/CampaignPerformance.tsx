// frontend/src/components/dashboard/CampaignPerformance.tsx
"use client";

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { TrendingUp } from "lucide-react";

const campaignData = [
  { name: "E-commerce", spend: 12400, revenue: 18600, roas: 1.5 },
  { name: "Finance", spend: 8900, revenue: 15800, roas: 1.77 },
  { name: "Gaming", spend: 15200, revenue: 21300, roas: 1.4 },
  { name: "Travel", spend: 6700, revenue: 12100, roas: 1.8 },
  { name: "Tech", spend: 9800, revenue: 16200, roas: 1.65 },
];

export function CampaignPerformance() {
  return (
    <div className='bg-white rounded-xl shadow-sm border border-slate-200 p-6'>
      <div className='flex items-center justify-between mb-6'>
        <div>
          <h3 className='text-lg font-semibold text-slate-900'>Campaign Performance</h3>
          <p className='text-sm text-slate-600'>Revenue vs Spend by vertical</p>
        </div>
        <div className='flex items-center space-x-2 text-sm'>
          <TrendingUp className='w-4 h-4 text-green-500' />
          <span className='text-green-600 font-medium'>+12.5% vs last week</span>
        </div>
      </div>

      <div className='h-80'>
        <ResponsiveContainer width='100%' height='100%'>
          <BarChart data={campaignData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray='3 3' stroke='#f1f5f9' />
            <XAxis dataKey='name' stroke='#64748b' fontSize={12} />
            <YAxis stroke='#64748b' fontSize={12} />
            <Tooltip
              contentStyle={{
                backgroundColor: "white",
                border: "1px solid #e2e8f0",
                borderRadius: "8px",
                boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
              }}
              formatter={(value: number, name: string) => [
                `$${value.toLocaleString()}`,
                name === "spend" ? "Spend" : "Revenue",
              ]}
            />
            <Bar dataKey='spend' fill='#f59e0b' radius={[4, 4, 0, 0]} name='spend' />
            <Bar dataKey='revenue' fill='#10b981' radius={[4, 4, 0, 0]} name='revenue' />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* ROAS Summary */}
      <div className='mt-4 pt-4 border-t border-slate-200'>
        <div className='grid grid-cols-5 gap-4'>
          {campaignData.map((campaign) => (
            <div key={campaign.name} className='text-center'>
              <div className='text-xs text-slate-500 mb-1'>{campaign.name}</div>
              <div
                className={`text-sm font-semibold ${
                  campaign.roas > 1.6
                    ? "text-green-600"
                    : campaign.roas > 1.4
                    ? "text-yellow-600"
                    : "text-red-600"
                }`}
              >
                {campaign.roas.toFixed(2)}x ROAS
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
