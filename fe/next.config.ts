import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  eslint: {
    // Next.js React Compiler rules flag pre-existing patterns in shadcn/ui and
    // react-hook-form that are intentional. Run linting separately with pnpm lint.
    ignoreDuringBuilds: true,
  },
  async redirects() {
    return [
      {
        source: '/admin/users/roles',
        destination: '/admin/roles',
        permanent: true,
      },
      {
        source: '/admin/permissions',
        destination: '/admin/roles',
        permanent: true,
      },
    ]
  },
}

export default nextConfig
