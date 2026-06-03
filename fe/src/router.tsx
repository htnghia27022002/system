import { createBrowserRouter, Navigate } from 'react-router-dom'

import {
  AdminRoute,
  GuestRoute,
  ProtectedRoute,
} from '@/features/auth'
import { PermissionRoute } from '@/features/access-control'
import { AdminLayout } from '@/layouts/admin-layout'
import { MainLayout } from '@/layouts/main-layout'
import { AdminOverviewPage } from '@/pages/admin-overview-page'
import { AdminRolesPage } from '@/pages/admin-roles-page'
import { AdminUsersPage } from '@/pages/admin-users-page'
import { HomePage } from '@/pages/home-page'
import { LoginPage } from '@/pages/login-page'
import { RegisterPage } from '@/pages/register-page'

export const router = createBrowserRouter([
  {
    element: <GuestRoute />,
    children: [
      {
        path: '/login',
        element: <LoginPage />,
      },
      {
        path: '/register',
        element: <RegisterPage />,
      },
    ],
  },
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
    ],
  },
  {
    path: '/admin',
    element: <ProtectedRoute />,
    children: [
      {
        element: <AdminRoute />,
        children: [
          {
            element: <AdminLayout />,
            children: [
              {
                element: <PermissionRoute permission="dashboard:read" />,
                handle: { breadcrumbKey: 'nav.dashboard' },
                children: [
                  {
                    index: true,
                    element: <AdminOverviewPage />,
                  },
                ],
              },
              {
                path: 'users',
                element: <PermissionRoute permission="users:read" />,
                handle: { breadcrumbKey: 'nav.users' },
                children: [
                  {
                    index: true,
                    element: <AdminUsersPage />,
                  },
                ],
              },
              {
                path: 'roles',
                element: <PermissionRoute permission="roles:read" />,
                handle: { breadcrumbKey: 'nav.roles' },
                children: [
                  {
                    index: true,
                    element: <AdminRolesPage />,
                  },
                ],
              },
              {
                path: 'users/roles',
                element: <Navigate to="/admin/roles" replace />,
              },
              {
                path: 'permissions',
                element: <Navigate to="/admin/roles" replace />,
              },
            ],
          },
        ],
      },
    ],
  },
])
