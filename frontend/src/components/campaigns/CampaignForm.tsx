"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  createCampaign,
  updateCampaign,
  Campaign,
  CreateCampaignInput,
  UpdateCampaignInput,
} from "@/lib/api/campaigns";

interface CampaignFormProps {
  initialCampaign?: Campaign;
  mode: "create" | "edit";
}

export function CampaignForm({ initialCampaign, mode }: CampaignFormProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    name: initialCampaign?.name || "",
    status: initialCampaign?.status || "active",
    budget: initialCampaign?.budget?.toString() || "",
    daily_budget: initialCampaign?.daily_budget?.toString() || "",
    target_cpa: initialCampaign?.target_cpa?.toString() || "",
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name || formData.name.length < 3) {
      newErrors.name = "Campaign name must be at least 3 characters";
    }
    if (formData.name.length > 255) {
      newErrors.name = "Campaign name must be less than 255 characters";
    }

    if (formData.budget && parseFloat(formData.budget) < 0) {
      newErrors.budget = "Budget must be positive";
    }

    if (formData.daily_budget && parseFloat(formData.daily_budget) < 0) {
      newErrors.daily_budget = "Daily budget must be positive";
    }

    if (formData.target_cpa && parseFloat(formData.target_cpa) < 0) {
      newErrors.target_cpa = "Target CPA must be positive";
    }

    if (formData.budget && formData.daily_budget) {
      const budget = parseFloat(formData.budget);
      const dailyBudget = parseFloat(formData.daily_budget);
      if (dailyBudget > budget) {
        newErrors.daily_budget = "Daily budget cannot exceed total budget";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      if (mode === "create") {
        const data: CreateCampaignInput = {
          name: formData.name,
          budget: formData.budget ? parseFloat(formData.budget) : undefined,
          daily_budget: formData.daily_budget
            ? parseFloat(formData.daily_budget)
            : undefined,
          target_cpa: formData.target_cpa
            ? parseFloat(formData.target_cpa)
            : undefined,
        };
        const result = await createCampaign(data);
        router.push(`/dashboard/campaigns/edit?id=${result.id}`);
      } else if (initialCampaign) {
        const data: UpdateCampaignInput = {
          id: initialCampaign.id,
          name: formData.name,
          status: formData.status as "active" | "paused" | "archived",
          budget: formData.budget ? parseFloat(formData.budget) : undefined,
          daily_budget: formData.daily_budget
            ? parseFloat(formData.daily_budget)
            : undefined,
          target_cpa: formData.target_cpa
            ? parseFloat(formData.target_cpa)
            : undefined,
        };
        await updateCampaign(data);
        router.push(`/dashboard/campaigns/edit?id= ${initialCampaign?.id}`);
      }
    } catch (err: unknown) {
      const errorMessage = err instanceof Error
        ? err.message
        : "Failed to save campaign";
      setError(errorMessage);
      console.error("Failed to save campaign:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    // Clear error for this field when user starts typing
    if (errors[name]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-2xl">
      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}

      <div className="space-y-6">
        {/* Campaign Name */}
        <div>
          <label
            htmlFor="name"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Campaign Name <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            id="name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.name ? "border-red-500" : "border-gray-300"
              }`}
            placeholder="e.g., Summer Sale 2024"
            required
          />
          {errors.name && (
            <p className="mt-1 text-sm text-red-600">{errors.name}</p>
          )}
        </div>

        {/* Status (Edit only) */}
        {mode === "edit" && (
          <div>
            <label
              htmlFor="status"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Status
            </label>
            <select
              id="status"
              name="status"
              value={formData.status}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="active">Active</option>
              <option value="paused">Paused</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        )}

        {/* Budget */}
        <div>
          <label
            htmlFor="budget"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Total Budget (USD)
          </label>
          <div className="relative">
            <span className="absolute left-3 top-2 text-gray-500">$</span>
            <input
              type="number"
              id="budget"
              name="budget"
              value={formData.budget}
              onChange={handleChange}
              className={`w-full pl-8 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.budget ? "border-red-500" : "border-gray-300"
                }`}
              placeholder="10000.00"
              step="0.01"
              min="0"
            />
          </div>
          {errors.budget && (
            <p className="mt-1 text-sm text-red-600">{errors.budget}</p>
          )}
          <p className="mt-1 text-xs text-gray-500">
            Optional: Leave empty for unlimited budget
          </p>
        </div>

        {/* Daily Budget */}
        <div>
          <label
            htmlFor="daily_budget"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Daily Budget (USD)
          </label>
          <div className="relative">
            <span className="absolute left-3 top-2 text-gray-500">$</span>
            <input
              type="number"
              id="daily_budget"
              name="daily_budget"
              value={formData.daily_budget}
              onChange={handleChange}
              className={`w-full pl-8 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.daily_budget ? "border-red-500" : "border-gray-300"
                }`}
              placeholder="500.00"
              step="0.01"
              min="0"
            />
          </div>
          {errors.daily_budget && (
            <p className="mt-1 text-sm text-red-600">{errors.daily_budget}</p>
          )}
          <p className="mt-1 text-xs text-gray-500">
            Optional: Maximum daily spend limit
          </p>
        </div>

        {/* Target CPA */}
        <div>
          <label
            htmlFor="target_cpa"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Target CPA (Cost Per Acquisition)
          </label>
          <div className="relative">
            <span className="absolute left-3 top-2 text-gray-500">$</span>
            <input
              type="number"
              id="target_cpa"
              name="target_cpa"
              value={formData.target_cpa}
              onChange={handleChange}
              className={`w-full pl-8 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.target_cpa ? "border-red-500" : "border-gray-300"
                }`}
              placeholder="25.00"
              step="0.01"
              min="0"
            />
          </div>
          {errors.target_cpa && (
            <p className="mt-1 text-sm text-red-600">{errors.target_cpa}</p>
          )}
          <p className="mt-1 text-xs text-gray-500">
            Optional: Target cost per conversion
          </p>
        </div>

        {/* Actions */}
        <div className="flex gap-4 pt-6 border-t border-gray-200">
          <button
            type="submit"
            disabled={loading}
            className="flex-1 px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {loading
              ? "Saving..."
              : mode === "create"
                ? "Create Campaign"
                : "Save Changes"}
          </button>
          <button
            type="button"
            onClick={() => router.back()}
            className="px-6 py-3 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 focus:ring-4 focus:ring-gray-200 transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </form>
  );
}
