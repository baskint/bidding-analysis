// frontend/src/components/dashboard/BidChart.tsx
'use client';;

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const bidData = [
  { time: "00:00", bids: 450, wins: 120, revenue: 2400 },
  { time: "04:00", bids: 320, wins: 85, revenue: 1800 },
  { time: "08:00", bids: 680, wins: 190, revenue: 3200 },
  { time: "12:00", bids: 890, wins: 280, revenue: 4100 },
  { time: "16:00", bids: 1200, wins: 350, revenue: 5200 },
  { time: "20:00", bids: 750, wins: 210, revenue: 3600 },
];

export function BidChart() {
  return (
    <div className='bg-white rounded-xl shadow-sm border border-slate-200 p-6'>
      <div className='flex items-center justify-between mb-6'>
        <div>
          <h3 className='text-lg font-semibold text-slate-900'>Bidding Activity</h3>
          <p className='text-sm text-slate-600'>24-hour bid volume and performance</p>
        </div>
        <div className='flex space-x-2'>
          <div className='flex items-center space-x-2'>
            <div className='w-3 h-3 bg-blue-500 rounded-full'></div>
            <span className='text-sm text-slate-600'>Bids</span>
          </div>
          <div className='flex items-center space-x-2'>
            <div className='w-3 h-3 bg-green-500 rounded-full'></div>
            <span className='text-sm text-slate-600'>Wins</span>
          </div>
        </div>
      </div>

      <div className='h-80'>
        <ResponsiveContainer width='100%' height='100%'>
          <LineChart data={bidData}>
            <CartesianGrid strokeDasharray='3 3' stroke='#f1f5f9' />
            <XAxis dataKey='time' stroke='#64748b' fontSize={12} />
            <YAxis stroke='#64748b' fontSize={12} />
            <Tooltip
              contentStyle={{
                backgroundColor: "white",
                border: "1px solid #e2e8f0",
                borderRadius: "8px",
                boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
              }}
            />
            <Line
              type='monotone'
              dataKey='bids'
              stroke='#3b82f6'
              strokeWidth={3}
              dot={{ fill: "#3b82f6", strokeWidth: 2, r: 4 }}
              activeDot={{ r: 6, stroke: "#3b82f6", strokeWidth: 2 }}
            />
            <Line
              type='monotone'
              dataKey='wins'
              stroke='#10b981'
              strokeWidth={3}
              dot={{ fill: "#10b981", strokeWidth: 2, r: 4 }}
              activeDot={{ r: 6, stroke: "#10b981", strokeWidth: 2 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
