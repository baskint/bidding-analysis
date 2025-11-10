import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'export',
  trailingSlash: false,
  images: {
    unoptimized: true,
  },
   // Optional: Add if you have dynamic routes
  skipTrailingSlashRedirect: true,
};

export default nextConfig;
