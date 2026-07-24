import type { Metadata } from 'next'

import { ToolsHubPage } from '@/features/tools'

export const metadata: Metadata = {
  title: 'Tools - System',
  description:
    'Browse the shared tools catalog. Open Webhooks and other utilities from one hub.',
  openGraph: {
    title: 'Tools - System',
    description:
      'Browse the shared tools catalog. Open Webhooks and other utilities from one hub.',
    type: 'website',
  },
}

export default function ToolsPage() {
  return <ToolsHubPage />
}
