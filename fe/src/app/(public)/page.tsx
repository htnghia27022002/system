import type { Metadata } from 'next'

import { LandingPage } from '@/features/landing'

export const metadata: Metadata = {
  title: 'System - Build tools people open',
  description:
    'Product home with a shared tools catalog and a dedicated tools hub.',
  openGraph: {
    title: 'System - Build tools people open',
    description:
      'Product home with a shared tools catalog and a dedicated tools hub.',
    type: 'website',
  },
}

export default function Page() {
  return <LandingPage />
}
