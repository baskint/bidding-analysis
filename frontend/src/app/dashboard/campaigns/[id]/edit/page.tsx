// frontend/src/app/dashboard/campaigns/[id]/edit/page.tsx
// NO "use client" directive

import { CampaignEditor } from "@/components/campaigns/CampaignEditor";

interface EditCampaignPageProps {
  params: Promise<{
    id: string;
  }>;
}
export default async function EditCampaignPage({
  params,
}: EditCampaignPageProps) {
  // 3. ❗ FIX: Await the promise to get the object with the 'id' property
  const resolvedParams = await params;

  // Now you can safely access the 'id' property on the resolved object
  return <CampaignEditor campaignId={resolvedParams.id} />;
}
