// frontend/src/components/campaigns/CampaignEditor.tsx
"use client"; // REQUIRED: This enables all the hooks

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Loader2 } from "lucide-react"; // Added Loader2 icon for spinner
import { CampaignForm } from "@/components/campaigns/CampaignForm";
import { getCampaign, Campaign } from "@/lib/api/campaigns";

// We only need the ID, which will be passed from the Server Component
interface CampaignEditorProps {
  campaignId: string;
}

export function CampaignEditor({ campaignId }: CampaignEditorProps) {
  const router = useRouter();
  const [campaign, setCampaign] = useState<Campaign | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadCampaign = async () => {
      try {
        setLoading(true);
        const data = await getCampaign(campaignId);
        const campaignData: Campaign = data;
        setCampaign(campaignData);
        setError(null);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "An unknown error occurred.";
        setError(`Failed to load campaign: ${errorMessage}`);
        console.error("Campaign fetch error:", err);
      } finally {
        setLoading(false);
      }
    };

    loadCampaign();
  }, [campaignId]); // ✅ Now only depends on campaignId

  // ---
  // ⏳ Loading State
  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          {/* Using lucide-react Loader2 icon for a cleaner spinner */}
          <Loader2 className="animate-spin w-8 h-8 text-blue-600 mx-auto" />
          <p className="mt-4 text-gray-600">Loading campaign details...</p>
        </div>
      </div>
    );
  }

  // ---
  // 🛑 Error/Not Found State
  if (error || !campaign) {
    return (
      <div className="flex items-center justify-center min-h-[400px] p-8">
        <div className="text-center bg-white p-8 rounded-lg shadow-xl border border-red-200">
          <p className="text-red-600 font-semibold text-lg">
            {error || "Campaign not found"}
          </p>
          <p className="text-gray-600 mt-2">
            Could not load campaign ID: **{campaignId}**
          </p>
          <button
            onClick={() => router.push("/dashboard/campaigns")}
            className="mt-6 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition duration-150"
          >
            Back to Campaigns
          </button>
        </div>
      </div>
    );
  }

  // ---
  // ✅ Success Render: Pass loaded data to the CampaignForm
  return (
    <div className="p-8">
      <div className="mb-8">
        <Link
          href={`/dashboard/campaigns/${campaign.id}`}
          className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900 mb-4 transition duration-150"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back to Campaign Details
        </Link>

        <h1 className="text-3xl font-bold text-gray-900">Edit Campaign</h1>
        <p className="text-gray-600 mt-1">Editing: **{campaign.name}**</p>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        {/* Pass the fully loaded and typed campaign object to the form */}
        <CampaignForm mode="edit" initialCampaign={campaign} />
      </div>
    </div>
  );
}
