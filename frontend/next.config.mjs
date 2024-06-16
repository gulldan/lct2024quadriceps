/** @type {import('next').NextConfig} */



const nextConfig = {
  output: "standalone",
  async rewrites() {
    const host = `http://bff:8888`

    return [
      {
        source: "/api/v1/:path*",
        destination: `${host}/api/v1/:path*`, // Proxy to Backend
      },
    ];
  },
};

export default nextConfig;
