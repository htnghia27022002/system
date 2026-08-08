import type { Metadata, Viewport } from 'next'
import type { ReactNode } from 'react'

import { AppProviders } from '@/components/providers/app-providers'
import { PwaInstallPrompt } from '@/components/common/pwa-install-prompt'
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
  manifest: '/manifest.webmanifest',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'black-translucent',
    title: 'System',
  },
  icons: {
    icon: [{ url: '/favicon.svg', type: 'image/svg+xml' }],
    apple: [{ url: '/icons/apple-touch-icon.png', sizes: '180x180' }],
  },
  other: {
    'mobile-web-app-capable': 'yes',
    'format-detection': 'telephone=no',
  },
}

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  viewportFit: 'cover',
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#3D72E0' },
    { media: '(prefers-color-scheme: dark)', color: '#3D72E0' },
  ],
}

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <AppProviders>
          <AuthHydrator>{children}</AuthHydrator>
          <PwaInstallPrompt />
          <Toaster richColors position="top-right" />
        </AppProviders>
      </body>
    </html>
  )
}
