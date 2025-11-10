// frontend/src/app/dashboard/campaigns/edit/page.tsx
"use client";

import { useSearchParams, useRouter } from "next/navigation";
import { CampaignEditor } from "@/components/campaigns/CampaignEditor";
import { useEffect } from "react";

export default function EditCampaignPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const campaignId = searchParams.get('id');

  useEffect(() => {
    if (!campaignId) {
      router.push('/dashboard/campaigns');
    }
  }, [campaignId, router]);

  if (!campaignId) {
    return <div>Loading...</div>;
  }

  return <CampaignEditor campaignId={campaignId} />;
}