import { Navigate, Outlet } from 'react-router-dom'

import { useAuthStore } from '@/store/auth-store'

export function AdminRoute() {
  const isAdmin = useAuthStore((state) => state.isAdmin)

  if (!isAdmin) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}
