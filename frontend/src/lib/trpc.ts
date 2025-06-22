// frontend/src/lib/trpc.ts
import { createTRPCReact } from '@trpc/react-query';
import { httpBatchLink } from '@trpc/client';

// Define your tRPC router type (this should match your Go backend)
export type AppRouter = {
  bidding: {
    processBid: {
      input: {
        campaignId: string;
        userId: string;
        floorPrice: number;
        deviceType: string;
        os: string;
        browser: string;
        country: string;
        region: string;
        city: string;
        keywords: string[];
        segmentId: string;
        segmentCategory: string;
        engagementScore: number;
        conversionProbability: number;
      };
      output: {
        bid_price: number;
        confidence: number;
        strategy: string;
        fraud_risk: boolean;
        prediction_id: string;
      };
    };
  };
  campaign: {
    getStats: {
      input: {
        campaignId: string;
        startTime?: string;
        endTime?: string;
      };
      output: any;
    };
    getBidHistory: {
      input: {
        campaignId: string;
        startTime: string;
        endTime: string;
        limit?: number;
        offset?: number;
      };
      output: any;
    };
  };
  analytics: {
    getFraudAlerts: {
      input: {
        startTime: string;
        endTime: string;
        severityThreshold?: number;
      };
      output: any;
    };
    getModelAccuracy: {
      input: {
        startTime: string;
        endTime: string;
        modelVersion?: string;
      };
      output: any;
    };
  };
};

export const trpc = createTRPCReact<AppRouter>();

export const trpcClient = trpc.createClient({
  links: [
    httpBatchLink({
      url: process.env.NODE_ENV === 'production' 
        ? 'https://bidding-analysis-539382269313.us-central1.run.app/trpc'
        : 'http://localhost:8080/trpc',
      headers: {
        'Content-Type': 'application/json',
      },
    }),
  ],
});
