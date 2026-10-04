// Lazy loaders for the mock API layer. Production code paths reach these only when
// NEXT_PUBLIC_USE_MOCK_API=true, so the mock modules ship as separate chunks that
// real deployments never download.

export { MockAccessControlError, MockAuthError } from './mock-errors'

export const loadAccessControlMock = () =>
  import('./access-control.mock').then((m) => m.mockAccessControlApi)
export const loadAuthMock = () => import('./auth.mock').then((m) => m.mockAuthApi)
export const loadDashboardMock = () =>
  import('./dashboard.mock').then((m) => m.mockDashboardApi)
export const loadHealthMock = () => import('./health.mock').then((m) => m.mockHealthApi)
