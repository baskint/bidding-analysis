// frontend/src/app/dashboard/campaigns/[id]/edit/page.tsx
"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { CampaignForm } from "@/components/campaigns/CampaignForm";
import { getCampaign, Campaign } from "@/lib/api/campaigns";

interface EditCampaignPageProps {
  params: {
    id: string;
  };
}

export default function EditCampaignPage({ params }: EditCampaignPageProps) {
  const router = useRouter();
  const [campaign, setCampaign] = useState<Campaign | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadCampaign();
  }, [params.id]);

  const loadCampaign = async () => {
    try {
      setLoading(true);
      const data = await getCampaign(params.id);
      // Extract just the campaign fields we need
      const campaignData: Campaign = {
        id: data.id,
        name: data.name,
        user_id: data.user_id,
        status: data.status,
        budget: data.budget,
        daily_budget: data.daily_budget,
        target_cpa: data.target_cpa,
        created_at: data.created_at,
        updated_at: data.updated_at,
      };
      setCampaign(campaignData);
      setError(null);
    } catch (err) {
      setError("Failed to load campaign");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading campaign...</p>
        </div>
      </div>
    );
  }

  if (error || !campaign) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <p className="text-red-600">{error || "Campaign not found"}</p>
          <button
            onClick={() => router.push("/dashboard/campaigns")}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Back to Campaigns
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="mb-8">
        <Link
          href={`/dashboard/campaigns/${campaign.id}`}
          className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back to Campaign
        </Link>

        <h1 className="text-3xl font-bold text-gray-900">Edit Campaign</h1>
        <p className="text-gray-600 mt-1">{campaign.name}</p>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <CampaignForm mode="edit" campaign={campaign} />
      </div>
    </div>
  );
}
