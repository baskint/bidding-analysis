"use client";

import { useState } from "react";
import {
  Zap,
  Play,
  Pause,
  Settings,
  TrendingUp,
  DollarSign,
} from "lucide-react";

export default function BiddingPage() {
  const [isLiveBidding, setIsLiveBidding] = useState(false);
  const [biddingStrategy, setBiddingStrategy] = useState("aggressive");

  const toggleLiveBidding = () => {
    setIsLiveBidding(!isLiveBidding);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 flex items-center">
            <Zap className="w-7 h-7 mr-3 text-blue-600" />
            Live Bidding Test
          </h1>
          <p className="text-gray-600 mt-1">
            Real-time bid optimization with AI-powered predictions
          </p>
        </div>

        {/* Live Bidding Toggle */}
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-3">
            <span className="text-sm font-medium text-gray-700">
              Live Bidding
            </span>
            <button
              onClick={toggleLiveBidding}
              className={`
                relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
                ${isLiveBidding ? "bg-green-600" : "bg-gray-200"}
              `}
            >
              <span
                className={`
                  inline-block h-4 w-4 transform rounded-full bg-white transition-transform
                  ${isLiveBidding ? "translate-x-6" : "translate-x-1"}
                `}
              />
            </button>
          </div>

          {isLiveBidding ? (
            <button
              onClick={toggleLiveBidding}
              className="flex items-center px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
            >
              <Pause className="w-4 h-4 mr-2" />
              Stop Bidding
            </button>
          ) : (
            <button
              onClick={toggleLiveBidding}
              className="flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
            >
              <Play className="w-4 h-4 mr-2" />
              Start Bidding
            </button>
          )}
        </div>
      </div>

      {/* Status Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                <Zap className="w-4 h-4 text-blue-600" />
              </div>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-500">Active Bids</p>
              <p className="text-2xl font-semibold text-gray-900">1,247</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                <TrendingUp className="w-4 h-4 text-green-600" />
              </div>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-500">Win Rate</p>
              <p className="text-2xl font-semibold text-gray-900">34.8%</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center">
                <DollarSign className="w-4 h-4 text-purple-600" />
              </div>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-500">Avg. Bid</p>
              <p className="text-2xl font-semibold text-gray-900">$2.34</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center ${
                  isLiveBidding ? "bg-green-100" : "bg-gray-100"
                }`}
              >
                <div
                  className={`w-3 h-3 rounded-full ${
                    isLiveBidding ? "bg-green-600 animate-pulse" : "bg-gray-400"
                  }`}
                />
              </div>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-500">Status</p>
              <p
                className={`text-2xl font-semibold ${
                  isLiveBidding ? "text-green-600" : "text-gray-900"
                }`}
              >
                {isLiveBidding ? "Live" : "Paused"}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Bidding Strategy */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Real-time Bid Stream
            </h3>

            {isLiveBidding ? (
              <div className="space-y-3">
                {/* Mock bid entries */}
                {[
                  {
                    id: 1,
                    campaign: "Mobile App Campaign",
                    bid: "$1.85",
                    status: "won",
                    time: "2 sec ago",
                  },
                  {
                    id: 2,
                    campaign: "E-commerce Sale",
                    bid: "$2.34",
                    status: "lost",
                    time: "5 sec ago",
                  },
                  {
                    id: 3,
                    campaign: "Brand Awareness",
                    bid: "$1.67",
                    status: "won",
                    time: "8 sec ago",
                  },
                  {
                    id: 4,
                    campaign: "Retargeting",
                    bid: "$3.12",
                    status: "won",
                    time: "12 sec ago",
                  },
                  {
                    id: 5,
                    campaign: "Local Services",
                    bid: "$0.95",
                    status: "lost",
                    time: "15 sec ago",
                  },
                ].map((bid) => (
                  <div
                    key={bid.id}
                    className="flex items-center justify-between py-3 px-4 bg-gray-50 rounded-lg"
                  >
                    <div className="flex items-center space-x-4">
                      <div
                        className={`w-3 h-3 rounded-full ${
                          bid.status === "won" ? "bg-green-500" : "bg-red-500"
                        }`}
                      />
                      <div>
                        <p className="font-medium text-gray-900">
                          {bid.campaign}
                        </p>
                        <p className="text-sm text-gray-500">{bid.time}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-gray-900">{bid.bid}</p>
                      <p
                        className={`text-sm capitalize ${
                          bid.status === "won"
                            ? "text-green-600"
                            : "text-red-600"
                        }`}
                      >
                        {bid.status}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <Pause className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">Bidding is currently paused</p>
                <p className="text-sm text-gray-400 mt-1">
                  Click &quot;Start Bidding&quot; to begin real-time bidding
                </p>
              </div>
            )}
          </div>
        </div>

        {/* Bidding Controls */}
        <div className="space-y-6">
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
              <Settings className="w-5 h-5 mr-2" />
              Bidding Strategy
            </h3>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Strategy Type
                </label>
                <select
                  value={biddingStrategy}
                  onChange={(e) => setBiddingStrategy(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-700"
                >
                  <option value="aggressive">Aggressive</option>
                  <option value="balanced">Balanced</option>
                  <option value="conservative">Conservative</option>
                  <option value="custom">Custom</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Max Bid Limit
                </label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-2.5 w-4 h-4 text-gray-700" />
                  <input
                    type="number"
                    step="0.01"
                    defaultValue="5.00"
                    className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-700"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Target CPA
                </label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-2.5 w-4 h-4 text-gray-700" />
                  <input
                    type="number"
                    step="0.01"
                    defaultValue="25.00"
                    className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-700"
                  />
                </div>
              </div>

              <button className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                Update Strategy
              </button>
            </div>
          </div>

          {/* AI Insights */}
          <div className="bg-gradient-to-r from-blue-50 to-purple-50 rounded-xl border border-blue-200 p-6">
            <h4 className="font-semibold text-gray-900 mb-3">🤖 AI Insights</h4>
            <div className="space-y-2 text-sm text-gray-700">
              <p>• Consider increasing bids for mobile traffic by 15%</p>
              <p>• High conversion probability detected for evening hours</p>
              <p>• Fraud risk is low across all campaigns</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
