// src/lib/trpc.ts
import { createTRPCReact } from '@trpc/react-query';
import { httpBatchLink } from '@trpc/client';

// Define your API types (you'll need to match your Go backend)
export interface User {
  id: string;
  username: string;
  createdAt: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

// Create the tRPC client
export const trpc = createTRPCReact<AppRouter>();

// Client configuration
export function getTRPCClient() {
  return trpc.createClient({
    links: [
      httpBatchLink({
        url: process.env.NODE_ENV === 'production' 
          ? 'https://bidding-analysis-539382269313.us-central1.run.app/trpc'
          : 'http://localhost:8080/trpc',
        // Add headers for authentication if needed
        headers() {
          const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
          return token ? { authorization: `Bearer ${token}` } : {};
        },
      }),
    ],
  });
}

// Type definitions for your Go backend routes
// You'll need to update these to match your actual Go tRPC routes
export type AppRouter = {
  auth: {
    login: {
      input: { username: string; password: string };
      output: AuthResponse;
    };
    register: {
      input: { username: string; password: string };
      output: AuthResponse;
    };
    me: {
      input: { token: string };
      output: { user: User };
    };
  };
};
