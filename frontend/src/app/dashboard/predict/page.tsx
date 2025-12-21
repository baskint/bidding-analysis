// frontend/src/app/dashboard/predict/page.tsx
'use client';
import { useState } from 'react';
import { Brain, TrendingUp, AlertCircle, Zap } from 'lucide-react';
import { apiPost, formatCurrency, formatPercent } from '@/lib/utils';
import type { PredictionRequest, PredictionResult } from '@/lib/types';

export default function PredictPage() {
  const [campaignId] = useState('550e8400-e29b-41d4-a716-446655440000');
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
      const requestData: PredictionRequest = {
        campaign_id: campaignId,
        floor_price: floorPrice,
        user_segment: userSegment,
        device_type: deviceType,
        country: country,
        keywords: ['technology', 'gadgets'],
        engagement_score: engagementScore,
        conversion_probability: conversionProbability,
      };

      const result = await apiPost<PredictionResult>('/trpc/bidding.predict', requestData);
      setPrediction(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to get prediction');
      console.error('Prediction error:', err);
    } finally {
      setLoading(false);
    }
  };

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 0.8) return 'text-green-600';
    if (confidence >= 0.6) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getConfidenceBg = (confidence: number) => {
    if (confidence >= 0.8) return 'bg-green-50 border-green-200';
    if (confidence >= 0.6) return 'bg-yellow-50 border-yellow-200';
    return 'bg-red-50 border-red-200';
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900">ML Bid Predictor</h1>
        <p className="text-slate-600 mt-1">Get AI-powered bid predictions in real-time</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Input Parameters */}
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
          <h2 className="text-lg font-semibold text-slate-900 mb-6">Bid Parameters</h2>

          <div className="space-y-6">
            {/* Floor Price */}
            <div>
              <label className="flex items-center justify-between text-sm font-medium text-slate-700 mb-2">
                <span>Floor Price</span>
                <span className="text-blue-600 font-semibold">{formatCurrency(floorPrice)}</span>
              </label>
              <input
                type="range"
                min="0.5"
                max="10"
                step="0.1"
                value={floorPrice}
                onChange={(e) => setFloorPrice(parseFloat(e.target.value))}
                className="w-full h-3 bg-gradient-to-r from-blue-200 to-blue-500 rounded-lg appearance-none cursor-pointer"
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
                <span className="text-purple-600 font-semibold">{formatPercent(engagementScore)}</span>
              </label>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                value={engagementScore}
                onChange={(e) => setEngagementScore(parseFloat(e.target.value))}
                className="w-full h-3 bg-gradient-to-r from-purple-200 to-purple-500 rounded-lg appearance-none cursor-pointer"
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
                <span className="text-green-600 font-semibold">{formatPercent(conversionProbability)}</span>
              </label>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                value={conversionProbability}
                onChange={(e) => setConversionProbability(parseFloat(e.target.value))}
                className="w-full h-3 bg-gradient-to-r from-green-200 to-green-500 rounded-lg appearance-none cursor-pointer"
              />
              <div className="flex justify-between text-xs text-slate-500 mt-1">
                <span>0%</span>
                <span>100%</span>
              </div>
            </div>

            {/* Device Type */}
            <div>
              <label className="text-sm font-medium text-slate-700 mb-2 block">Device Type</label>
              <div className="grid grid-cols-3 gap-2">
                {['mobile', 'desktop', 'tablet'].map((type) => (
                  <button
                    key={type}
                    onClick={() => setDeviceType(type)}
                    className={'px-4 py-2 rounded-lg text-sm font-medium capitalize transition-all ' +
                      (deviceType === type
                        ? 'bg-blue-600 text-white shadow-md'
                        : 'bg-slate-100 text-slate-700 hover:bg-slate-200')}
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
                className="w-full px-4 py-3 bg-slate-50 border border-slate-300 rounded-lg text-slate-900 font-medium focus:ring-2 focus:ring-blue-500 hover:bg-slate-100 transition-colors cursor-pointer"
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
                className="w-full px-4 py-3 bg-slate-50 border border-slate-300 rounded-lg text-slate-900 font-medium focus:ring-2 focus:ring-blue-500 hover:bg-slate-100 transition-colors cursor-pointer"
              >
                <option value="tech_enthusiasts">Tech Enthusiasts</option>
                <option value="shoppers">Shoppers</option>
                <option value="gamers">Gamers</option>
                <option value="business">Business</option>
                <option value="students">Students</option>
              </select>
            </div>

            {/* Predict Button */}
            <button
              onClick={handlePredict}
              disabled={loading}
              className="w-full py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg font-semibold hover:from-blue-700 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center space-x-2"
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

        {/* Prediction Results */}
        <div className="space-y-6">
          {/* Instructions (shown when no prediction) */}
          {!prediction && !error && (
            <div className="bg-gradient-to-br from-blue-50 to-purple-50 rounded-xl shadow-sm border border-blue-200 p-6">
              <div className="flex items-start space-x-3">
                <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
                  <Zap className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <h3 className="font-semibold text-slate-900 mb-2">How it works</h3>
                  <ul className="space-y-2 text-sm text-slate-600">
                    <li className="flex items-start">
                      <span className="text-blue-600 mr-2">•</span>
                      <span>Adjust the parameters using the sliders and dropdowns</span>
                    </li>
                    <li className="flex items-start">
                      <span className="text-purple-600 mr-2">•</span>
                      <span>Click "Get ML Prediction" to get AI-powered recommendations</span>
                    </li>
                    <li className="flex items-start">
                      <span className="text-green-600 mr-2">•</span>
                      <span>Our XGBoost model analyzes 300 trees with 97.91% accuracy</span>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          )}

          {/* Error State */}
          {error && (
            <div className="bg-red-50 rounded-xl shadow-sm border border-red-200 p-6">
              <div className="flex items-start space-x-3">
                <AlertCircle className="w-6 h-6 text-red-600 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold text-red-900 mb-1">Prediction Failed</h3>
                  <p className="text-sm text-red-700">{error}</p>
                </div>
              </div>
            </div>
          )}

          {/* Prediction Result */}
          {prediction && (
            <div className={'rounded-xl shadow-lg border-2 p-6 ' + getConfidenceBg(prediction.confidence)}>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-slate-900">Prediction Result</h3>
                <div className={'px-3 py-1 rounded-full text-xs font-bold ' +
                  (prediction.strategy === 'ml_optimized'
                    ? 'bg-green-100 text-green-700'
                    : 'bg-yellow-100 text-yellow-700')}>
                  {prediction.strategy === 'ml_optimized' ? 'ML Optimized' : 'Rule Based'}
                </div>
              </div>

              <div className="space-y-4">
                {/* Predicted Bid */}
                <div className="text-center py-4 bg-white rounded-lg">
                  <div className="text-sm text-slate-600 mb-1">Recommended Bid</div>
                  <div className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 text-transparent bg-clip-text">
                    {formatCurrency(prediction.predicted_bid)}
                  </div>
                </div>

                {/* Confidence */}
                <div className="flex items-center justify-between py-3 border-t border-slate-200">
                  <span className="text-sm text-slate-600">Confidence</span>
                  <div className="flex items-center space-x-2">
                    <div className="w-32 h-2 bg-slate-200 rounded-full overflow-hidden">
                      <div
                        className={'h-full transition-all ' +
                          (prediction.confidence >= 0.8 ? 'bg-green-500' :
                            prediction.confidence >= 0.6 ? 'bg-yellow-500' : 'bg-red-500')}
                        style={{ width: `${prediction.confidence * 100}%` }}
                      />
                    </div>
                    <span className={'font-bold ' + getConfidenceColor(prediction.confidence)}>
                      {formatPercent(prediction.confidence, 0)}
                    </span>
                  </div>
                </div>

                {/* Comparison to Floor */}
                <div className="flex items-center justify-between py-3 border-t border-slate-200">
                  <span className="text-sm text-slate-600">vs Floor Price</span>
                  <div className="flex items-center space-x-2">
                    <TrendingUp className={'w-4 h-4 ' +
                      (prediction.predicted_bid > floorPrice ? 'text-green-600' : 'text-red-600')} />
                    <span className={'font-bold ' +
                      (prediction.predicted_bid > floorPrice ? 'text-green-600' : 'text-red-600')}>
                      {formatPercent((prediction.predicted_bid - floorPrice) / floorPrice, 1)}
                    </span>
                  </div>
                </div>

                {/* Fraud Risk */}
                <div className="flex items-center justify-between py-3 border-t border-slate-200">
                  <span className="text-sm text-slate-600">Fraud Risk</span>
                  <span className={'px-3 py-1 rounded-full text-xs font-bold ' +
                    (prediction.fraud_risk ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700')}>
                    {prediction.fraud_risk ? 'High' : 'Low'}
                  </span>
                </div>

                {/* Reasoning */}
                {prediction.reasoning && (
                  <div className="pt-3 border-t border-slate-200">
                    <div className="text-sm text-slate-600 mb-1">AI Reasoning</div>
                    <p className="text-sm text-slate-800 bg-white rounded-lg p-3">
                      {prediction.reasoning}
                    </p>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
