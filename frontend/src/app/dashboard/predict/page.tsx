// frontend/src/app/dashboard/predict/page.tsx
'use client';
import { useState } from 'react';
import { Brain, TrendingUp, AlertCircle, Zap } from 'lucide-react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const getAuthHeaders = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

interface PredictionResult {
  predicted_bid: number;
  confidence: number;
  strategy: string;
  fraud_risk: boolean;
  reasoning: string;
}

export default function PredictPage() {
  const [campaignId] = useState('550e8400-e29b-41d4-a716-446655440000'); // Default campaign
  const [floorPrice, setFloorPrice] = useState(1.5);
  const [engagementScore, setEngagementScore] = useState(0.75);
  const [conversionProbability, setConversionProbability] = useState(0.3);
  const [deviceType, setDeviceType] = useState('mobile');
  const [country, setCountry] = useState('US');
  const [userSegment, setUserSegment] = useState('tech_enthusiasts');

  const [prediction, setPrediction] = useState<PredictionResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handlePredict = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_BASE_URL}/trpc/bidding.predict`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          campaign_id: campaignId,
          floor_price: floorPrice,
          user_segment: userSegment,
          device_type: deviceType,
          country: country,
          keywords: ['technology', 'gadgets'],
          engagement_score: engagementScore,
          conversion_probability: conversionProbability,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to get prediction');
      }

      const data = await response.json();

      if (data.result && data.result.data) {
        setPrediction(data.result.data);
      } else {
        throw new Error('Invalid response format');
      }
    } catch (err) {
      console.error('Prediction error:', err);
      setError(err instanceof Error ? err.message : 'Failed to get prediction');
    } finally {
      setLoading(false);
    }
  };

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 0.8) return 'text-green-600 bg-green-50 border-green-200';
    if (confidence >= 0.6) return 'text-yellow-600 bg-yellow-50 border-yellow-200';
    return 'text-red-600 bg-red-50 border-red-200';
  };

  const getStrategyBadge = (strategy: string) => {
    if (strategy === 'ml_optimized') {
      return (
        <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gradient-to-r from-blue-500 to-purple-600 text-white">
          <Brain className="w-4 h-4 mr-1" />
          ML Optimized
        </span>
      );
    }
    return (
      <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gray-100 text-gray-800">
        Rule-Based
      </span>
    );
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-100">ML Bid Predictor</h1>
          <p className="text-slate-600 mt-1">Use machine learning to predict optimal bid prices</p>
        </div>
        <div className="flex items-center space-x-2 px-4 py-2 bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg border border-blue-200">
          <Zap className="w-5 h-5 text-blue-600" />
          <span className="text-sm font-medium text-blue-900">XGBoost Model (97.9% Accuracy)</span>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Input Form */}
        <div className="lg:col-span-2 bg-white rounded-xl shadow-sm border border-slate-200 p-6">
          <h2 className="text-lg font-semibold text-slate-900 mb-6">Bid Parameters</h2>

          <div className="space-y-6">
            {/* Floor Price */}
            <div>
              <label className="flex items-center justify-between text-sm font-medium text-slate-700 mb-2">
                <span>Floor Price</span>
                <span className="text-blue-600 font-semibold">${floorPrice.toFixed(2)}</span>
              </label>
              <input
                type="range"
                min="0.5"
                max="10"
                step="0.1"
                value={floorPrice}
                onChange={(e) => setFloorPrice(parseFloat(e.target.value))}
                className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-blue-600"
              />
              <div className="flex justify-between text-xs text-slate-500 mt-1">
                <span>$0.50</span>
                <span>$10.00</span>
              </div>
            </div>

            {/* Engagement Score */}
            <div>
              <label className="flex items-center justify-between text-sm font-medium text-slate-700 mb-2">
                <span>Engagement Score</span>
                <span className="text-blue-600 font-semibold">{(engagementScore * 100).toFixed(0)}%</span>
              </label>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                value={engagementScore}
                onChange={(e) => setEngagementScore(parseFloat(e.target.value))}
                className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
              />
              <div className="flex justify-between text-xs text-slate-500 mt-1">
                <span>0%</span>
                <span>100%</span>
              </div>
            </div>

            {/* Conversion Probability */}
            <div>
              <label className="flex items-center justify-between text-sm font-medium text-slate-700 mb-2">
                <span>Conversion Probability</span>
                <span className="text-blue-600 font-semibold">{(conversionProbability * 100).toFixed(0)}%</span>
              </label>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                value={conversionProbability}
                onChange={(e) => setConversionProbability(parseFloat(e.target.value))}
                className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-green-600"
              />
              <div className="flex justify-between text-xs text-slate-500 mt-1">
                <span>0%</span>
                <span>100%</span>
              </div>
            </div>

            {/* Device Type */}
            <div>
              <label className="text-sm font-medium text-slate-700 mb-2 block">Device Type</label>
              <div className="grid grid-cols-3 gap-3">
                {['mobile', 'desktop', 'tablet'].map((type) => (
                  <button
                    key={type}
                    onClick={() => setDeviceType(type)}
                    className={`px-4 py-3 rounded-lg border-2 font-medium text-sm capitalize transition-all ${deviceType === type
                      ? 'border-blue-600 bg-blue-50 text-blue-900'
                      : 'border-slate-200 bg-white text-slate-700 hover:border-slate-300'
                      }`}
                  >
                    {type}
                  </button>
                ))}
              </div>
            </div>

            {/* Country */}
            <div>
              <label className="text-sm font-medium text-slate-700 mb-2 block">Country</label>
              <select
                value={country}
                onChange={(e) => setCountry(e.target.value)}
                className="w-full px-4 py-3 bg-slate-50 border border-slate-300 rounded-lg text-slate-900 font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 hover:bg-slate-100 transition-colors cursor-pointer"
              >
                <option value="US">United States</option>
                <option value="UK">United Kingdom</option>
                <option value="CA">Canada</option>
                <option value="AU">Australia</option>
                <option value="DE">Germany</option>
                <option value="FR">France</option>
                <option value="JP">Japan</option>
              </select>
            </div>

            {/* User Segment */}
            <div>
              <label className="text-sm font-medium text-slate-700 mb-2 block">User Segment</label>
              <select
                value={userSegment}
                onChange={(e) => setUserSegment(e.target.value)}
                className="w-full px-4 py-3 bg-slate-50 border border-slate-300 rounded-lg text-slate-900 font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 hover:bg-slate-100 transition-colors cursor-pointer"
              >
                <option value="tech_enthusiasts">Tech Enthusiasts</option>
                <option value="shoppers">Shoppers</option>
                <option value="gamers">Gamers</option>
                <option value="business">Business Professionals</option>
                <option value="students">Students</option>
              </select>
            </div>

            {/* Predict Button */}
            <button
              onClick={handlePredict}
              disabled={loading}
              className="w-full py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg font-semibold hover:from-blue-700 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center space-x-2"
            >
              {loading ? (
                <>
                  <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  <span>Predicting...</span>
                </>
              ) : (
                <>
                  <Brain className="w-5 h-5" />
                  <span>Get ML Prediction</span>
                </>
              )}
            </button>
          </div>
        </div>

        {/* Results Panel */}
        <div className="space-y-6">
          {/* Prediction Result */}
          {prediction && (
            <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-slate-900">Prediction</h3>
                {getStrategyBadge(prediction.strategy)}
              </div>

              <div className="space-y-4">
                {/* Predicted Bid */}
                <div className="text-center p-6 bg-gradient-to-br from-blue-50 to-purple-50 rounded-lg">
                  <div className="text-sm text-slate-600 mb-1">Optimal Bid Price</div>
                  <div className="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-purple-600">
                    ${prediction.predicted_bid.toFixed(4)}
                  </div>
                </div>

                {/* Confidence */}
                <div className={`p-4 rounded-lg border-2 ${getConfidenceColor(prediction.confidence)}`}>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium">Confidence</span>
                    <span className="text-lg font-bold">{(prediction.confidence * 100).toFixed(0)}%</span>
                  </div>
                  <div className="w-full bg-white rounded-full h-2">
                    <div
                      className="h-2 rounded-full bg-gradient-to-r from-blue-500 to-purple-600"
                      style={{ width: `${prediction.confidence * 100}%` }}
                    />
                  </div>
                </div>

                {/* Comparison */}
                <div className="p-4 bg-slate-50 rounded-lg">
                  <div className="text-xs text-slate-600 mb-2">vs Floor Price</div>
                  <div className="flex items-center space-x-2">
                    <TrendingUp className={`w-5 h-5 ${prediction.predicted_bid > floorPrice ? 'text-green-600' : 'text-red-600'
                      }`} />
                    <span className={`text-lg font-bold ${prediction.predicted_bid > floorPrice ? 'text-green-600' : 'text-red-600'
                      }`}>
                      {prediction.predicted_bid > floorPrice ? '+' : ''}
                      {((prediction.predicted_bid - floorPrice) / floorPrice * 100).toFixed(1)}%
                    </span>
                  </div>
                </div>

                {/* Reasoning */}
                <div className="p-4 bg-blue-50 rounded-lg">
                  <div className="flex items-start space-x-2">
                    <AlertCircle className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" />
                    <p className="text-sm text-blue-900">{prediction.reasoning}</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Error Display */}
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <div className="flex items-start space-x-2">
                <AlertCircle className="w-5 h-5 text-red-600 mt-0.5" />
                <div>
                  <div className="text-sm font-medium text-red-900">Prediction Failed</div>
                  <div className="text-sm text-red-700 mt-1">{error}</div>
                </div>
              </div>
            </div>
          )}

          {/* Instructions */}
          {!prediction && !error && (
            <div className="bg-slate-50 rounded-xl border border-slate-200 p-6">
              <h3 className="text-sm font-semibold text-slate-900 mb-3">How it works</h3>
              <ul className="space-y-2 text-sm text-slate-600">
                <li className="flex items-start">
                  <span className="text-blue-600 mr-2">1.</span>
                  Adjust the bid parameters using the sliders and dropdowns
                </li>
                <li className="flex items-start">
                  <span className="text-blue-600 mr-2">2.</span>
                  Click  &ldquo;Get ML Prediction&rdquo; to query the XGBoost model
                </li>
                <li className="flex items-start">
                  <span className="text-blue-600 mr-2">3.</span>
                  View the optimal bid price with confidence score
                </li>
                <li className="flex items-start">
                  <span className="text-blue-600 mr-2">4.</span>
                  Compare different scenarios to optimize your bidding strategy
                </li>
              </ul>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
