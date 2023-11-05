const apiurl = "https://api-openchokin.walnuts.dev/v1";

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  async rewrites() {
    return [
      {
        source: "/api/back/:path*",
        destination: `${apiurl}/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
