// frontend/src/app/dashboard/compare/page.tsx
'use client';
import { useState } from 'react';
import { Brain, TrendingUp, AlertCircle, Zap, Lock, Unlock, Trophy } from 'lucide-react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const getAuthHeaders = () => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

interface ScenarioParams {
  name: string;
  floorPrice: number;
  engagementScore: number;
  conversionProbability: number;
  deviceType: string;
  country: string;
  userSegment: string;
}

interface PredictionResult {
  predicted_bid: number;
  confidence: number;
  strategy: string;
  fraud_risk: boolean;
  reasoning: string;
}

export default function ComparePage() {
  const [campaignId] = useState('550e8400-e29b-41d4-a716-446655440000');

  // Three scenarios
  const [scenarios, setScenarios] = useState<ScenarioParams[]>([
    {
      name: 'Conservative',
      floorPrice: 1.0,
      engagementScore: 0.6,
      conversionProbability: 0.2,
      deviceType: 'desktop',
      country: 'US',
      userSegment: 'tech_enthusiasts',
    },
    {
      name: 'Balanced',
      floorPrice: 1.5,
      engagementScore: 0.75,
      conversionProbability: 0.3,
      deviceType: 'mobile',
      country: 'US',
      userSegment: 'tech_enthusiasts',
    },
    {
      name: 'Aggressive',
      floorPrice: 2.5,
      engagementScore: 0.9,
      conversionProbability: 0.4,
      deviceType: 'mobile',
      country: 'US',
      userSegment: 'tech_enthusiasts',
    },
  ]);

  const [predictions, setPredictions] = useState<(PredictionResult | null)[]>([null, null, null]);
  const [loading, setLoading] = useState<boolean[]>([false, false, false]);
  const [lockedParams, setLockedParams] = useState({
    country: false,
    userSegment: false,
  });

  const updateScenario = (index: number, field: keyof ScenarioParams, value: any) => {
    const newScenarios = [...scenarios];
    newScenarios[index] = { ...newScenarios[index], [field]: value };

    // If locked, update all scenarios
    if (lockedParams[field as keyof typeof lockedParams]) {
      newScenarios.forEach((scenario, i) => {
        if (i !== index) {
          newScenarios[i] = { ...newScenarios[i], [field]: value };
        }
      });
    }

    setScenarios(newScenarios);
  };

  const predictScenario = async (index: number) => {
    const newLoading = [...loading];
    newLoading[index] = true;
    setLoading(newLoading);

    try {
      const scenario = scenarios[index];
      const response = await fetch(API_BASE_URL + '/trpc/bidding.predict', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          campaign_id: campaignId,
          floor_price: scenario.floorPrice,
          user_segment: scenario.userSegment,
          device_type: scenario.deviceType,
          country: scenario.country,
          keywords: ['technology', 'gadgets'],
          engagement_score: scenario.engagementScore,
          conversion_probability: scenario.conversionProbability,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to get prediction');
      }

      const data = await response.json();

      if (data.result && data.result.data) {
        const newPredictions = [...predictions];
        newPredictions[index] = data.result.data;
        setPredictions(newPredictions);
      }
    } catch (err) {
      console.error('Prediction error:', err);
    } finally {
      const newLoading = [...loading];
      newLoading[index] = false;
      setLoading(newLoading);
    }
  };

  const predictAll = async () => {
    await Promise.all([0, 1, 2].map(i => predictScenario(i)));
  };

  const getBestScenario = () => {
    const validPredictions = predictions
      .map((p, i) => ({ prediction: p, index: i }))
      .filter(({ prediction }) => prediction !== null);

    if (validPredictions.length === 0) return null;

    // Find scenario with highest predicted bid
    return validPredictions.reduce((best, current) => {
      if (!best.prediction || !current.prediction) return best;
      return current.prediction.predicted_bid > best.prediction.predicted_bid ? current : best;
    });
  };

  // Add this function near the top with other functions
  const resetToDefaults = () => {
    setScenarios([
      {
        name: 'Conservative',
        floorPrice: 1.0,
        engagementScore: 0.6,
        conversionProbability: 0.2,
        deviceType: 'desktop',
        country: 'US',
        userSegment: 'tech_enthusiasts',
      },
      {
        name: 'Balanced',
        floorPrice: 1.5,
        engagementScore: 0.75,
        conversionProbability: 0.3,
        deviceType: 'mobile',
        country: 'US',
        userSegment: 'tech_enthusiasts',
      },
      {
        name: 'Aggressive',
        floorPrice: 2.5,
        engagementScore: 0.9,
        conversionProbability: 0.4,
        deviceType: 'mobile',
        country: 'US',
        userSegment: 'tech_enthusiasts',
      },
    ]);
    setPredictions([null, null, null]);
    setLockedParams({ country: false, userSegment: false });
  };

  const bestScenario = getBestScenario();

  const scenarioColors = [
    'from-blue-500 to-blue-600',
    'from-purple-500 to-purple-600',
    'from-green-500 to-green-600',
  ];

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">Bid Strategy Comparison</h1>
          <p className="text-slate-600 mt-1">Compare multiple bid scenarios side-by-side</p>
        </div>
        <div className="flex items-center space-x-3">
          <button
            onClick={resetToDefaults}
            className="px-4 py-2 bg-white border-2 border-slate-300 text-slate-700 rounded-lg font-semibold hover:bg-slate-50 transition-all"
          >
            Reset
          </button>
          <button
            onClick={predictAll}
            disabled={loading.some(l => l)}
            className="px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg font-semibold hover:from-blue-700 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center space-x-2"
          >
            <Brain className="w-5 h-5" />
            <span>Compare All Scenarios</span>
          </button>
        </div>
      </div>

      {/* Locked Parameters */}
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-4">
        <h3 className="text-sm font-semibold text-slate-900 mb-3">Synchronized Parameters</h3>
        <div className="flex items-center space-x-4">
          <button
            onClick={() => setLockedParams({ ...lockedParams, country: !lockedParams.country })}
            className={'flex items-center space-x-2 px-3 py-2 rounded-lg transition-colors ' + (lockedParams.country ? 'bg-blue-50 text-blue-700' : 'bg-slate-50 text-slate-600')}
          >
            {lockedParams.country ? <Lock className="w-4 h-4" /> : <Unlock className="w-4 h-4" />}
            <span className="text-sm font-medium">Country</span>
          </button>
          <button
            onClick={() => setLockedParams({ ...lockedParams, userSegment: !lockedParams.userSegment })}
            className={'flex items-center space-x-2 px-3 py-2 rounded-lg transition-colors ' + (lockedParams.userSegment ? 'bg-blue-50 text-blue-700' : 'bg-slate-50 text-slate-600')}
          >
            {lockedParams.userSegment ? <Lock className="w-4 h-4" /> : <Unlock className="w-4 h-4" />}
            <span className="text-sm font-medium">User Segment</span>
          </button>
        </div>
      </div>

      {/* Scenarios Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {scenarios.map((scenario, index) => (
          <div key={index} className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
            {/* Scenario Header */}
            <div className="mb-6">
              <input
                type="text"
                value={scenario.name}
                onChange={(e) => updateScenario(index, 'name', e.target.value)}
                className={'text-lg font-bold mb-2 px-3 py-1 rounded-lg bg-gradient-to-r ' + scenarioColors[index] + ' text-white w-full'}
              />
              {bestScenario && bestScenario.index === index && (
                <div className="flex items-center space-x-1 text-yellow-600 text-sm font-semibold mt-2">
                  <Trophy className="w-4 h-4" />
                  <span>Best Performance</span>
                </div>
              )}
            </div>

            {/* Parameters */}
            <div className="space-y-4 mb-6">
              {/* Floor Price */}
              <div>
                <label className="flex items-center justify-between text-xs font-medium text-slate-700 mb-1">
                  <span>Floor Price</span>
                  <span className="text-blue-600 font-semibold">${scenario.floorPrice.toFixed(2)}</span>
                </label>
                <input
                  type="range"
                  min="0.5"
                  max="5"
                  step="0.1"
                  value={scenario.floorPrice}
                  onChange={(e) => updateScenario(index, 'floorPrice', parseFloat(e.target.value))}
                  className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-blue-600"
                />
              </div>

              {/* Engagement Score */}
              <div>
                <label className="flex items-center justify-between text-xs font-medium text-slate-700 mb-1">
                  <span>Engagement</span>
                  <span className="text-purple-600 font-semibold">{(scenario.engagementScore * 100).toFixed(0)}%</span>
                </label>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  value={scenario.engagementScore}
                  onChange={(e) => updateScenario(index, 'engagementScore', parseFloat(e.target.value))}
                  className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
                />
              </div>

              {/* Conversion Probability */}
              <div>
                <label className="flex items-center justify-between text-xs font-medium text-slate-700 mb-1">
                  <span>Conversion</span>
                  <span className="text-green-600 font-semibold">{(scenario.conversionProbability * 100).toFixed(0)}%</span>
                </label>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  value={scenario.conversionProbability}
                  onChange={(e) => updateScenario(index, 'conversionProbability', parseFloat(e.target.value))}
                  className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-green-600"
                />
              </div>

              {/* Device Type */}
              <div>
                <label className="text-xs font-medium text-slate-700 mb-1 block">Device</label>
                <div className="grid grid-cols-3 gap-1">
                  {['mobile', 'desktop', 'tablet'].map((type) => (
                    <button
                      key={type}
                      onClick={() => updateScenario(index, 'deviceType', type)}
                      className={'px-2 py-1 rounded text-xs font-medium capitalize transition-all ' + (scenario.deviceType === type ? 'bg-blue-600 text-white' : 'bg-slate-100 text-slate-600')}
                    >
                      {type.slice(0, 3)}
                    </button>
                  ))}
                </div>
              </div>

              {/* Country */}
              <div>
                <label className="text-xs font-medium text-slate-700 mb-1 flex items-center">
                  Country
                  {lockedParams.country && <Lock className="w-3 h-3 ml-1 text-blue-600" />}
                </label>
                <select
                  value={scenario.country}
                  onChange={(e) => updateScenario(index, 'country', e.target.value)}
                  className="w-full px-2 py-1 text-xs bg-slate-50 border border-slate-300 rounded-lg text-slate-900 font-medium"
                >
                  <option value="US">US</option>
                  <option value="UK">UK</option>
                  <option value="CA">CA</option>
                  <option value="DE">DE</option>
                </select>
              </div>

              {/* User Segment */}
              <div>
                <label className="text-xs font-medium text-slate-700 mb-1 flex items-center">
                  Segment
                  {lockedParams.userSegment && <Lock className="w-3 h-3 ml-1 text-blue-600" />}
                </label>
                <select
                  value={scenario.userSegment}
                  onChange={(e) => updateScenario(index, 'userSegment', e.target.value)}
                  className="w-full px-2 py-1 text-xs bg-slate-50 border border-slate-300 rounded-lg text-slate-900 font-medium"
                >
                  <option value="tech_enthusiasts">Tech</option>
                  <option value="shoppers">Shop</option>
                  <option value="gamers">Game</option>
                  <option value="business">B2B</option>
                </select>
              </div>
            </div>

            {/* Predict Button */}
            <button
              onClick={() => predictScenario(index)}
              disabled={loading[index]}
              className={'w-full py-2 rounded-lg font-semibold text-white transition-all flex items-center justify-center space-x-2 bg-gradient-to-r ' + scenarioColors[index]}
            >
              {loading[index] ? (
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              ) : (
                <>
                  <Zap className="w-4 h-4" />
                  <span className="text-sm">Predict</span>
                </>
              )}
            </button>

            {/* Prediction Result */}
            {predictions[index] && (
              <div className="mt-4 p-4 bg-gradient-to-br from-slate-50 to-slate-100 rounded-lg">
                <div className="text-center mb-2">
                  <div className="text-xs text-slate-600 mb-1">Predicted Bid</div>
                  <div className={'text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r ' + scenarioColors[index]}>
                    ${predictions[index]!.predicted_bid.toFixed(4)}
                  </div>
                </div>
                <div className="text-xs text-slate-600 text-center">
                  {(predictions[index]!.confidence * 100).toFixed(0)}% confidence
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Comparison Summary */}
      {predictions.some(p => p !== null) && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
          <h3 className="text-lg font-semibold text-slate-900 mb-4">Comparison Summary</h3>
          <div className="grid grid-cols-3 gap-4">
            {scenarios.map((scenario, index) => {
              const prediction = predictions[index];
              if (!prediction) return null;

              const isBest = bestScenario && bestScenario.index === index;

              return (
                <div key={index} className={'p-4 rounded-lg border-2 ' + (isBest ? 'border-yellow-400 bg-yellow-50' : 'border-slate-200')}>
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-semibold text-slate-900">{scenario.name}</span>
                    {isBest && <Trophy className="w-5 h-5 text-yellow-600" />}
                  </div>
                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span className="text-slate-600">Bid:</span>
                      <span className="font-bold">${prediction.predicted_bid.toFixed(4)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-600">vs Floor:</span>
                      <span className={'font-bold ' + (prediction.predicted_bid > scenario.floorPrice ? 'text-green-600' : 'text-red-600')}>
                        {((prediction.predicted_bid - scenario.floorPrice) / scenario.floorPrice * 100).toFixed(1)}%
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-600">Confidence:</span>
                      <span className="font-bold">{(prediction.confidence * 100).toFixed(0)}%</span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
