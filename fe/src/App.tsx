import { useEffect } from 'react'
import { RouterProvider } from 'react-router-dom'

import { AppProviders } from '@/components/providers/app-providers'
import { Toaster } from '@/components/ui/sonner'
import { router } from '@/router'
import { useAuthStore } from '@/store/auth-store'

function AuthHydrator({ children }: { children: React.ReactNode }) {
  const hydrateFromStorage = useAuthStore((state) => state.hydrateFromStorage)

  useEffect(() => {
    hydrateFromStorage()
  }, [hydrateFromStorage])

  return children
}

export default function App() {
  return (
    <AppProviders>
      <AuthHydrator>
        <RouterProvider router={router} />
      </AuthHydrator>
      <Toaster richColors position="top-right" />
    </AppProviders>
  )
}
