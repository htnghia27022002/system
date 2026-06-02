import { Navigate, Outlet } from 'react-router-dom'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

export function GuestRoute() {
  const user = useAuthStore((state) => state.user)

  if (user) {
    return <Navigate to={getPostLoginPath(user.role)} replace />
  }

  return <Outlet />
}
