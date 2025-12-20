// frontend/src/app/dashboard/activity/page.tsx
'use client';
import { useState, useEffect, useCallback } from 'react';
import {
  Activity,
  TrendingDown,
  CheckCircle,
  XCircle,
  Smartphone,
  Monitor,
  Tablet,
  Globe,
  RefreshCw,
  Filter,
} from 'lucide-react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const getAuthHeaders = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': 'Bearer ' + token }),
  };
};

interface BidEvent {
  bid_event_id: string;
  campaign_id: string;
  bid_price: number;
  win_price?: number;
  won: boolean;
  converted: boolean;
  timestamp: string;
  segment_category: string;
  device_type: string;
  country: string;
}

export default function ActivityPage() {
  const [events, setEvents] = useState<BidEvent[]>([]);
  const [filteredEvents, setFilteredEvents] = useState<BidEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date());

  // Filters
  const [statusFilter, setStatusFilter] = useState<'all' | 'won' | 'lost' | 'converted'>('all');
  const [deviceFilter, setDeviceFilter] = useState<string>('all');
  const [countryFilter, setCountryFilter] = useState<string>('all');

  const fetchEvents = useCallback(async () => {
    try {
      const response = await fetch(API_BASE_URL + '/trpc/bidding.stream', {
        method: 'GET',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        throw new Error('Failed to fetch bid stream');
      }

      const data = await response.json();
      const bidEvents = data.result?.data || [];
      setEvents(bidEvents);
      setLastRefresh(new Date());
    } catch (err) {
      console.error('Error fetching bid stream:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Apply filters
  useEffect(() => {
    let filtered = [...events];

    if (statusFilter === 'won') {
      filtered = filtered.filter(e => e.won);
    } else if (statusFilter === 'lost') {
      filtered = filtered.filter(e => !e.won);
    } else if (statusFilter === 'converted') {
      filtered = filtered.filter(e => e.converted);
    }

    if (deviceFilter !== 'all') {
      filtered = filtered.filter(e => e.device_type === deviceFilter);
    }

    if (countryFilter !== 'all') {
      filtered = filtered.filter(e => e.country === countryFilter);
    }

    setFilteredEvents(filtered);
  }, [events, statusFilter, deviceFilter, countryFilter]);

  // Initial load
  useEffect(() => {
    fetchEvents();
  }, [fetchEvents]);

  // Auto-refresh every 5 seconds
  useEffect(() => {
    if (!autoRefresh) return;

    const interval = setInterval(() => {
      fetchEvents();
    }, 5000);

    return () => clearInterval(interval);
  }, [autoRefresh, fetchEvents]);

  const getDeviceIcon = (device: string) => {
    switch (device.toLowerCase()) {
      case 'mobile': return Smartphone;
      case 'desktop': return Monitor;
      case 'tablet': return Tablet;
      default: return Monitor;
    }
  };

  const getStatusColor = (event: BidEvent) => {
    if (event.converted) return 'border-yellow-400 bg-yellow-50';
    if (event.won) return 'border-green-400 bg-green-50';
    return 'border-slate-300 bg-slate-50';
  };

  const getStatusIcon = (event: BidEvent) => {
    if (event.converted) return { Icon: CheckCircle, color: 'text-yellow-600' };
    if (event.won) return { Icon: CheckCircle, color: 'text-green-600' };
    return { Icon: XCircle, color: 'text-slate-500' };
  };

  const uniqueDevices = Array.from(new Set(events.map(e => e.device_type)));
  const uniqueCountries = Array.from(new Set(events.map(e => e.country)));

  const stats = {
    total: filteredEvents.length,
    won: filteredEvents.filter(e => e.won).length,
    lost: filteredEvents.filter(e => !e.won).length,
    converted: filteredEvents.filter(e => e.converted).length,
    totalSpent: filteredEvents.reduce((sum, e) => sum + (e.win_price || 0), 0),
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">Live Bid Activity</h1>
          <p className="text-slate-600 mt-1">Real-time stream of bid events</p>
        </div>
        <div className="flex items-center space-x-3">
          <button
            onClick={() => setAutoRefresh(!autoRefresh)}
            className={'px-4 py-2 rounded-lg font-medium transition-all flex items-center space-x-2 ' +
              (autoRefresh ? 'bg-green-100 text-green-700 border-2 border-green-300' : 'bg-slate-100 text-slate-700 border-2 border-slate-300')}
          >
            <RefreshCw className={'w-4 h-4 ' + (autoRefresh ? 'animate-spin' : '')} />
            <span>{autoRefresh ? 'Auto-refresh On' : 'Auto-refresh Off'}</span>
          </button>
          <button
            onClick={fetchEvents}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 flex items-center space-x-2"
          >
            <RefreshCw className="w-4 h-4" />
            <span>Refresh Now</span>
          </button>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        <div className="bg-white rounded-lg shadow-sm border border-slate-200 p-4">
          <div className="text-sm text-slate-600">Total Bids</div>
          <div className="text-2xl font-bold text-slate-900">{stats.total}</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-green-200 p-4">
          <div className="text-sm text-green-600">Won</div>
          <div className="text-2xl font-bold text-green-600">{stats.won}</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-slate-200 p-4">
          <div className="text-sm text-slate-600">Lost</div>
          <div className="text-2xl font-bold text-slate-600">{stats.lost}</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-yellow-200 p-4">
          <div className="text-sm text-yellow-600">Converted</div>
          <div className="text-2xl font-bold text-yellow-600">{stats.converted}</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-blue-200 p-4">
          <div className="text-sm text-blue-600">Total Spent</div>
          <div className="text-2xl font-bold text-blue-600">${stats.totalSpent.toFixed(2)}</div>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-4">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <Filter className="w-5 h-5 text-slate-600" />
            <span className="text-sm font-semibold text-slate-700">Filters:</span>
          </div>

          {/* Status Filter */}
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value as 'all' | 'won' | 'lost' | 'converted')}
            className="px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-sm font-medium"
          >
            <option value="all">All Status</option>
            <option value="won">Won Only</option>
            <option value="lost">Lost Only</option>
            <option value="converted">Converted Only</option>
          </select>

          {/* Device Filter */}
          <select
            value={deviceFilter}
            onChange={(e) => setDeviceFilter(e.target.value)}
            className="px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-sm font-medium"
          >
            <option value="all">All Devices</option>
            {uniqueDevices.map(device => (
              <option key={device} value={device}>{device}</option>
            ))}
          </select>

          {/* Country Filter */}
          <select
            value={countryFilter}
            onChange={(e) => setCountryFilter(e.target.value)}
            className="px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-sm font-medium"
          >
            <option value="all">All Countries</option>
            {uniqueCountries.map(country => (
              <option key={country} value={country}>{country}</option>
            ))}
          </select>

          <div className="flex-1"></div>
          <div className="text-xs text-slate-500">
            Last updated: {lastRefresh.toLocaleTimeString()}
          </div>
        </div>
      </div>

      {/* Event Timeline */}
      {loading ? (
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
            <p className="text-slate-600">Loading activity...</p>
          </div>
        </div>
      ) : filteredEvents.length === 0 ? (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-12 text-center">
          <Activity className="w-16 h-16 text-slate-400 mx-auto mb-4" />
          <p className="text-slate-600">No bid events match your filters</p>
        </div>
      ) : (
        <div className="space-y-3">
          {filteredEvents.map((event, index) => {
            const DeviceIcon = getDeviceIcon(event.device_type);
            const { Icon: StatusIcon, color: statusColor } = getStatusIcon(event);

            return (
              <div
                key={event.bid_event_id}
                className={'rounded-lg border-l-4 p-4 transition-all hover:shadow-md ' + getStatusColor(event)}
                style={{ animation: index < 5 ? 'slideIn 0.3s ease-out' : 'none' }}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-4 flex-1">
                    {/* Status Icon */}
                    <div className={'w-10 h-10 rounded-full flex items-center justify-center ' +
                      (event.converted ? 'bg-yellow-100' : event.won ? 'bg-green-100' : 'bg-slate-200')}>
                      <StatusIcon className={'w-5 h-5 ' + statusColor} />
                    </div>

                    {/* Bid Info */}
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-1">
                        <span className="font-semibold text-slate-900">
                          ${event.bid_price.toFixed(4)}
                        </span>
                        {event.won && event.win_price && (
                          <>
                            <TrendingDown className="w-4 h-4 text-green-600" />
                            <span className="text-sm font-medium text-green-600">
                              Won at ${event.win_price.toFixed(4)}
                            </span>
                          </>
                        )}
                        {event.converted && (
                          <span className="px-2 py-1 bg-yellow-200 text-yellow-800 text-xs font-bold rounded-full">
                            CONVERTED
                          </span>
                        )}
                      </div>
                      <div className="flex items-center space-x-4 text-sm text-slate-600">
                        <div className="flex items-center space-x-1">
                          <DeviceIcon className="w-4 h-4" />
                          <span className="capitalize">{event.device_type}</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <Globe className="w-4 h-4" />
                          <span>{event.country}</span>
                        </div>
                        <span className="px-2 py-0.5 bg-slate-200 text-slate-700 text-xs rounded-full">
                          {event.segment_category}
                        </span>
                      </div>
                    </div>

                    {/* Timestamp */}
                    <div className="text-right">
                      <div className="text-xs text-slate-500">
                        {new Date(event.timestamp).toLocaleTimeString()}
                      </div>
                      <div className="text-xs text-slate-400">
                        {new Date(event.timestamp).toLocaleDateString()}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

      <style jsx>{`
        @keyframes slideIn {
          from {
            opacity: 0;
            transform: translateX(-20px);
          }
          to {
            opacity: 1;
            transform: translateX(0);
          }
        }
      `}</style>
    </div>
  );
}