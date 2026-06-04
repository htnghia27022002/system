import type { Metadata } from 'next'

import { HomePage } from '@/features/home'

export const metadata: Metadata = {
  title: 'Home',
  description: 'Welcome to System App — your dashboard and tools hub.',
  openGraph: {
    title: 'Home | System App',
    description: 'Welcome to System App — your dashboard and tools hub.',
    type: 'website',
  },
}

export default function Page() {
  return <HomePage />
}
