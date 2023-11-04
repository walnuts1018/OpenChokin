const API_URL = "https://api-openchokin.walnuts.dev/v1";

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
