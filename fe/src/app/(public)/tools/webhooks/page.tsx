import type { Metadata } from 'next'

import { WebhooksToolPage } from '@/features/tools'

export const metadata: Metadata = {
  title: 'Webhooks - Tools',
  description:
    'Webhooks tool workspace shell. Inspect deliveries and debug integrations from System Tools.',
  openGraph: {
    title: 'Webhooks - Tools',
    description:
      'Webhooks tool workspace shell. Inspect deliveries and debug integrations from System Tools.',
    type: 'website',
  },
}

export default function WebhooksPage() {
  return <WebhooksToolPage />
}
