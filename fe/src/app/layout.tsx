import type { Metadata } from 'next'
import type { ReactNode } from 'react'

import { AppProviders } from '@/components/providers/app-providers'
import { Toaster } from '@/components/ui/sonner'
import { AuthHydrator } from '@/components/common/auth-hydrator'

import '@/styles/index.css'

export const metadata: Metadata = {
  title: {
    default: 'System App',
    template: '%s | System App',
  },
  description: 'System application with admin dashboard and access control.',
  robots: { index: true, follow: true },
  icons: {
    icon: [{ url: '/favicon.svg', type: 'image/svg+xml' }],
    apple: [{ url: '/favicon.svg', type: 'image/svg+xml' }],
  },
}

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <AppProviders>
          <AuthHydrator>{children}</AuthHydrator>
          <Toaster richColors position="top-right" />
        </AppProviders>
      </body>
    </html>
  )
}
