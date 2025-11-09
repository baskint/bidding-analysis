// frontend/src/app/dashboard/campaigns/[id]/page.tsx
// REMOVE: "use client";

import { notFound } from "next/navigation";
import { getCampaign, CampaignDetail } from "@/lib/api/campaigns"; // Import server-side data fetcher
import CampaignDetailClient from "@/components/campaigns/CampaignDetailClient"; // Import new Client Component

// 👈 FIX 1: The correct Server Component prop type for a dynamic route
interface CampaignDetailPageProps {
  params: Promise<{
    id: string;
  }>;
}

// 👈 FIX 2: The Page component is now an async function
export default async function CampaignDetailPage({
  params,
}: CampaignDetailPageProps) {
  let initialCampaignData: CampaignDetail;
  const resolvedParams = await params;
  // 💡 Server-side Data Fetching
  try {
    initialCampaignData = await getCampaign(resolvedParams.id);
  } catch (error) {
    // Handle error (e.g., campaign not found)
    console.error("Failed to fetch campaign on server:", error);
    notFound(); // Use Next.js utility for 404
  }

  // Pass the initial data to the Client Component.
  // The Client Component will still contain the activation/pause logic.
  return <CampaignDetailClient initialCampaign={initialCampaignData} />;
}
