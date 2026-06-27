import type { Metadata } from 'next'

import { LandingPage } from '@/features/landing'

export const metadata: Metadata = {
  title: 'htnghia — Web3 Scrollytelling',
  description: 'GSAP scroll-driven Web3 landing prototype with pinned sections and scrubbed animations.',
  openGraph: {
    title: 'htnghia — Web3 Scrollytelling',
    description: 'GSAP scroll-driven Web3 landing prototype with pinned sections and scrubbed animations.',
    type: 'website',
  },
}

export default function Page() {
  return <LandingPage />
}
