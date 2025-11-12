import { Target } from "lucide-react";
import type { KeywordAnalysis } from "@/lib/api/analytics";

interface KeywordAnalysisTableProps {
  keywords: KeywordAnalysis[];
  loading: boolean;
}

export function KeywordAnalysisTable({
  keywords,
  loading,
}: KeywordAnalysisTableProps) {
  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-slate-200 rounded w-1/4"></div>
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-12 bg-slate-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (keywords.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 text-center text-slate-600">
        No keyword data available
      </div>
    );
  }

  const formatCurrency = (value: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);
  const formatPercent = (value: number) => `${(value * 100).toFixed(1)}%`;

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <div className="flex items-center mb-6">
        <Target className="w-6 h-6 text-purple-600 mr-3" />
        <h2 className="text-xl font-bold text-slate-900">
          Top Keywords by Performance
        </h2>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="text-left text-sm text-slate-600 border-b border-slate-200">
              <th className="pb-3 font-semibold">Keyword</th>
              <th className="pb-3 font-semibold text-right">Bids</th>
              <th className="pb-3 font-semibold text-right">Win Rate</th>
              <th className="pb-3 font-semibold text-right">Conv. Rate</th>
              <th className="pb-3 font-semibold text-right">Spend</th>
              <th className="pb-3 font-semibold text-right">ROAS</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {keywords.length > 0 && keywords?.map((kw, idx) => (
              <tr key={idx} className="hover:bg-slate-50">
                <td className="py-3 font-medium text-slate-900">
                  {kw.keyword}
                </td>
                <td className="py-3 text-right text-slate-600">
                  {(kw.totalBids ?? 0).toLocaleString()}
                </td>
                <td className="py-3 text-right">
                  <span className="inline-flex items-center px-2 py-1 rounded-full bg-blue-100 text-blue-700 text-xs font-medium">
                    {formatPercent(kw.winRate ?? 0)}
                  </span>
                </td>
                <td className="py-3 text-right">
                  <span className="inline-flex items-center px-2 py-1 rounded-full bg-green-100 text-green-700 text-xs font-medium">
                    {formatPercent(kw.conversionRate ?? 0)}
                  </span>
                </td>
                <td className="py-3 text-right text-slate-900 font-medium">
                  {formatCurrency(kw.spend ?? 0)}
                </td>
                <td
                  className={`py-3 text-right font-semibold ${
                    (kw.roas ?? 0) >= 2
                      ? "text-green-600"
                      : (kw.roas ?? 0) >= 1
                        ? "text-yellow-600"
                        : "text-red-600"
                  }`}
                >
                  {(kw.roas ?? 0).toFixed(2)}x
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
