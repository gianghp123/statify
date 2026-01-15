import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  experimental: {
    globalNotFound: true,
    serverActions: {
      bodySizeLimit: '50mb',
    },
  },
};

export default nextConfig;
