const API_URL = process.env.API_URL;

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  async rewrites() {
    return [
      {
        source: "/api/back/:path*",
        destination: `${API_URL}/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
