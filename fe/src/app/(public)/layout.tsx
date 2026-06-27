import type { ReactNode } from 'react'

import { LandingLayout } from '@/layouts/landing-layout'

export default function PublicGroupLayout({ children }: { children: ReactNode }) {
  return <LandingLayout>{children}</LandingLayout>
}
