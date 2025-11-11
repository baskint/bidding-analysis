// frontend/src/components/dashboard/Sidebar.tsx
'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  BarChart3,
  Zap,
  Shield,
  Settings,
  TrendingUp,
  Brain,
  AlertTriangle,
  Activity,
} from 'lucide-react';

const navigation = [
  { name: 'Overview', href: '/dashboard', icon: BarChart3 },
  { name: 'Live Bidding', href: '/dashboard/bidding', icon: Zap },
  { name: 'Campaigns', href: '/dashboard/campaigns', icon: TrendingUp },
  { name: 'Analytics', href: '/dashboard/analytics', icon: Activity },
  { name: 'ML Models', href: '/dashboard/models', icon: Brain },
  { name: 'Fraud Detection', href: '/dashboard/fraud', icon: Shield },
  { name: 'Alerts', href: '/dashboard/alerts', icon: AlertTriangle },
  { name: 'Settings', href: '/dashboard/settings', icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <div className="fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-xl border-r border-slate-200">
      <div className="flex flex-col h-full">
        {/* Logo */}
        <div className="flex items-center px-6 py-4 border-b border-slate-200">
          <div className="flex items-center space-x-3">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <Zap className="w-5 h-5 text-white" />
            </div>
            <span className="text-xl font-bold text-slate-900">BidAnalyzer</span>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-4 py-6 space-y-2 overflow-y-auto">
          {navigation.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.name}
                href={item.href}
                className={`
                  flex items-center px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200
                  ${isActive
                    ? 'bg-gradient-to-r from-blue-50 to-purple-50 text-blue-700 border-l-4 border-blue-600'
                    : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
                  }
                `}
              >
                <item.icon className={`w-5 h-5 mr-3 ${isActive ? 'text-blue-600' : 'text-slate-400'}`} />
                {item.name}
              </Link>
            );
          })}
        </nav>

        {/* Status */}
        <div className="px-4 py-4 border-t border-slate-200">
          <div className="flex items-center space-x-3 p-3 bg-green-50 rounded-lg">
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span className="text-sm font-medium text-green-700">System Online</span>
          </div>
        </div>
      </div>
    </div>
  );
}
