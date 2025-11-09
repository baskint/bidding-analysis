// frontend/src/app/dashboard/campaigns/new/page.tsx
"use client";

import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { CampaignForm } from "@/components/campaigns/CampaignForm";

export default function NewCampaignPage() {
  return (
    <div className="p-8">
      <div className="mb-8">
        <Link
          href="/dashboard/campaigns"
          className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back to Campaigns
        </Link>

        <h1 className="text-3xl font-bold text-gray-900">
          Create New Campaign
        </h1>
        <p className="text-gray-600 mt-1">Set up a new advertising campaign</p>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <CampaignForm mode="create" />
      </div>
    </div>
  );
}
